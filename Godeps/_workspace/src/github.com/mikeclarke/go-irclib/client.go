package irc

import (
	"github.com/mikeclarke/go-broadcast"
	"bufio"
	"crypto/tls"
	"fmt"
	"log"
	"net"
	"regexp"
	"strconv"
	"strings"
	"time"
)

const (
	VERSION = "go-irclib v0.1"
)

type IRCClient struct {
	Nickname              string   // Nickname the client will use
	Password              string   // Password used to log on to the server
	RealName              string   // Supplied to the server as "Real name" or "ircname"
	Username              string   // Supplied to the server as the "User name""
	UserInfo              string   // Sent in reply to a USERINFO CTCP query
	FingerReply           string   // Sent in reply to a FINGER CTCP query
	VersionName           string   // CTCP VERSION reply, client name
	VersionNum            string   // CTCP VERSION reply, client version
	VersionEnv            string   // CTCP VERSION reply, environment the client is running in
	SourceURL             string   // CTCP SOURCE reply, the URL of the source code of this client
	LineRate              int      // Minimum delay between lines sent to the server
	ErroneousNickFallback string   // Default nickname assigned when ERR_ERRONEUSNICKNAME

	// Connection settings
	Server                string   // Host name to attempt a connection to
	SSL                   bool     // Use TLS for secure connection
	socket                net.Conn

	registered            bool     // Whether or not the user is registered
	hostname              string   // Host name of the IRC server the client is connected to
	heartbeatInterval     float64  // Interval, in seconds, to send PING messages for keepalive

	// Communication channels
	pwrite                             chan string // Channel for writing messages to IRC server
	readerExit, writerExit, pingerExit chan bool   // Channels for notifying goroutine stops
	endping                            chan bool   // Channel for stopping ping goroutine
	errchan                            chan error  // Channel for dumping errors
	broadcast                          *broadcast.Broadcaster

	lastMessage     time.Time
	currentNickname string
	stopped         bool
	HighlightRE     *regexp.Regexp
}

type Event struct {
	Raw       string  // Raw message string
	Prefix    string
	Command   string
	Arguments []string
	Client    *IRCClient
	Highlight bool
}

func parseIRCMessage(s string) (string, string, []string) {
	// http://twistedmatrix.com/trac/browser/trunk/twisted/words/protocols/irc.py#54
	prefix := ""
	trailing := ""
	command := ""
	args := []string{}

	if s[0] == ':' {
		splits := strings.SplitN(s[1:], " ", 2)
		prefix, s = splits[0], splits[1]
	}

	if strings.Contains(s, " :") {
		splits := strings.SplitN(s, " :", 2)
		s, trailing = splits[0], splits[1]
		args = strings.Fields(s)
		args = append(args, trailing)
	} else {
		args = strings.Fields(s)
	}
	command, args = args[0], args[1:]
	return prefix, command, args
}

func (irc *IRCClient) readLoop() {
	br := bufio.NewReaderSize(irc.socket, 512)

	for {
		msg, err := br.ReadString('\n')

		if err != nil {
			irc.errchan <- err
			break
		}

		irc.lastMessage = time.Now()
		msg = strings.Trim(msg, "\r\n")

		// Parse raw message into Event struct
		prefix, cmd, args := parseIRCMessage(msg)

		event := &Event{
			Raw: msg,
			Prefix: prefix,
			Command: cmd,
			Arguments: args,
			Client: irc,
			Highlight: irc.IsHighlight(msg),
		}

		// Publish on broadcast channel
		irc.broadcast.Write(event)
	}

	irc.readerExit <- true
}

func (irc *IRCClient) writeLoop() {
	for {
		b, ok := <-irc.pwrite
		if !ok || b == "" || irc.socket == nil {
			break
		}

		log.Printf("--> %s", strings.Trim(b, "\r\n"))

		_, err := irc.socket.Write([]byte(b))
		if err != nil {
			irc.errchan <- err
			break
		}
	}
	irc.writerExit <- true
}

//Pings the server if we have not recived any messages for 5 minutes
func (irc *IRCClient) pingLoop() {
	ticker := time.NewTicker(1 * time.Minute)   //Tick every minute.
	ticker2 := time.NewTicker(15 * time.Minute) //Tick every 15 minutes.
	for {
		select {
		case <-ticker.C:
			// Ping if we haven't received anything from the server within 4 minutes
			if time.Since(irc.lastMessage) >= (4 * time.Minute) {
				irc.SendRawf("PING %d", time.Now().UnixNano())
			}
		case <-ticker2.C:
			// Ping every 15 minutes.
			irc.SendRawf("PING %d", time.Now().UnixNano())

			// Try to recapture nickname if it's not as configured.
			if irc.Nickname != irc.currentNickname {
				irc.currentNickname = irc.Nickname
				irc.SendRawf("NICK %s", irc.Nickname)
			}
		case <-irc.endping:
			// Shut down everything
			ticker.Stop()
			ticker2.Stop()
			irc.pingerExit <- true
			return
		}
	}
}

func (irc *IRCClient) Run() {
	for !irc.stopped {
		err := <-irc.errchan
		if irc.stopped {
			break
		}
		log.Printf("errchan: %s\n", err)
		irc.Disconnect()
		irc.Connect(irc.Server)
	}
}

func (irc *IRCClient) Quit() {
	irc.SendRaw("QUIT")
	irc.stopped = true
	irc.Disconnect()
}

func (irc *IRCClient) Join(channel string) {
	irc.SendRawf("JOIN %s\r\n", channel)
}

func (irc *IRCClient) Part(channel string) {
	irc.SendRawf("PART %s\r\n", channel)
}

func (irc *IRCClient) Notice(target, message string) {
	irc.SendRawf("NOTICE %s :%s\r\n", target, message)
}

func (irc *IRCClient) Noticef(target, format string, a ...interface{}) {
	irc.Notice(target, fmt.Sprintf(format, a...))
}

func (irc *IRCClient) Privmsg(target, message string) {
	irc.SendRawf("PRIVMSG %s :%s\r\n", target, message)
}

func (irc *IRCClient) Privmsgf(target, format string, a ...interface{}) {
	irc.Privmsg(target, fmt.Sprintf(format, a...))
}

func (irc *IRCClient) SendRaw(message string) {
	irc.pwrite <- message + "\r\n"
}

func (irc *IRCClient) SendRawf(format string, a ...interface{}) {
	irc.SendRaw(fmt.Sprintf(format, a...))
}

func (irc *IRCClient) Nick(n string) {
	irc.Nickname = n
	irc.SendRawf("NICK %s", n)
}

func (irc *IRCClient) GetNick() string {
	return irc.currentNickname
}

func (irc *IRCClient) IsHighlight(msg string) bool {
	return strings.Contains(msg, fmt.Sprintf("%s:", irc.Nickname))
}

// Sends all buffered messages (if possible),
// stops all goroutines and then closes the socket.
func (irc *IRCClient) Disconnect() {
	close(irc.pwrite)
	irc.endping <- true

	<-irc.readerExit
	<-irc.writerExit
	<-irc.pingerExit
	irc.socket.Close()
	irc.socket = nil
}

func (irc *IRCClient) Reconnect() error {
	return irc.Connect(irc.Server)
}

func (irc *IRCClient) Connect(server string) error {
	irc.Server = server
	irc.stopped = false

	var err error
	if irc.SSL {
		irc.socket, err = tls.Dial("tcp", irc.Server, nil)
	} else {
		irc.socket, err = net.Dial("tcp", irc.Server)
	}
	if err != nil {
		return err
	}
	log.Printf("Connected to %s (%s)\n", irc.Server, irc.socket.RemoteAddr())

	irc.pwrite = make(chan string, 1024)
	irc.errchan = make(chan error, 2)

	go irc.readLoop()
	go irc.writeLoop()
	go irc.pingLoop()

	if len(irc.Password) > 0 {
		irc.SendRawf("PASS %s\r\n", irc.Password)
	}
	irc.SendRawf("NICK %s\r\n", irc.Nickname)
	irc.SendRawf("USER %s 0.0.0.0 0.0.0.0 :%s\r\n", irc.Username, irc.Username)
	return nil
}

func (irc *IRCClient) AddHandler(f func (*Event)) {
	messages := irc.broadcast.Listen(1024)

	go func() {
		for e := range messages {
			event := e.(*Event)
			f(event)
		}
	}()
}

func defaultHandlers(event *Event) {
	client := event.Client

	switch event.Command {
	case "PING":
		client.SendRawf("PONG %s", event.Arguments[len(event.Arguments)-1])

	case "437":
		// client.currentNickname = client.currentNickname + "_"
		// client.SendRawf("NICK %s", client.currentNickname)

	case "433":
		if len(client.currentNickname) > 8 {
			client.currentNickname = "_" + client.currentNickname
		} else {
			client.currentNickname = client.currentNickname + "_"
		}
		client.SendRawf("NICK %s", client.currentNickname)

	case "PONG":
		ns, _ := strconv.ParseInt(event.Arguments[1], 10, 64)
		delta := time.Duration(time.Now().UnixNano() - ns)
		log.Printf("Lag: %v\n", delta)

	case "001":
		// client.currentNickname = event.Arguments[0]
	}
}

func New(nick, user string) *IRCClient {
	irc := &IRCClient{
		Nickname:       nick,
		Username:       user,
		readerExit:     make(chan bool),
		writerExit:     make(chan bool),
		pingerExit:     make(chan bool),
		endping:        make(chan bool),
		broadcast:      broadcast.NewBroadcaster(1024),
	}

	// Add default IRC client handlers
	irc.AddHandler(defaultHandlers)

	return irc
}

