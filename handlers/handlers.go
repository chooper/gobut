
package handlers

import (
	irc "github.com/mikeclarke/go-irclib"
	"log"
	"regexp"
	"strings"
    "github.com/chooper/go-bot/robutdb"
)

type ChannelState struct {
	Server	string
	Channel	string
	Users	[]string
}

type ServerState struct {
	Server		string
	Channels	map[string] ChannelState
}

var Servers = make(map[string] *ServerState)

func DebugHandler(event *irc.Event) {
	log.Println("-> ", event)
	
	if event.Command != "PRIVMSG" { 
		return
	}
	
	message_RE := regexp.MustCompile(`^\.dump(.*)$`)
	matches := message_RE.FindStringSubmatch(event.Arguments[1])
	
	if len(matches) < 1 {
		return
	}
	
	event.Client.Privmsgf(event.Arguments[0], "state: %q", Servers)
}

func EchoHandler(event *irc.Event) {
	if event.Command != "PRIVMSG" {
		return
	}

	message_RE := regexp.MustCompile(`^\.echo\s*(.*)$`)
	matches := message_RE.FindStringSubmatch(event.Arguments[1])
		
	if len(matches) < 1 {
		return
	}

	echo := strings.Join(strings.Fields(matches[0])[1:], " ")
	event.Client.Privmsg(event.Arguments[0], echo)
}

/*
2013/06/15 18:58:53 XXX In channel =  [GoTest @hoop]
2013/06/15 18:58:53 User joined channel =  [#sandbox]
2013/06/15 18:58:53 Mode change =  [GoTest +i]
2013/06/15 18:58:57 ->  &{:hoop!~hoop@X.X.X.X MODE #sandbox +o GoTest hoop!~hoop@X.X.X.X MODE [#sandbox +o GoTest] 0xc2000b1000 false}
2013/06/15 18:58:57 Mode change =  [#sandbox +o GoTest]
*/

func NamesHandler(event *irc.Event) {
	if event.Command != "353" {
		return
	}

	server := event.Client.Server
	irc_chan := event.Arguments[2]
	users := strings.Fields(event.Arguments[3])	// TODO: Track modes
	
	var server_state *ServerState
	var ok bool
	if server_state, ok = Servers[server]; !ok {
		server_state = new(ServerState)
		server_state.Server = server
		server_state.Channels = make(map[string] ChannelState)
	}
	
	server_state.Channels[irc_chan] = ChannelState{server, irc_chan, users}
	Servers[server] = server_state
}

func PartHandler(event *irc.Event){ 
	if event.Command != "PART" {
		return
	}
	event.Client.SendRawf("NAMES %s", event.Arguments[0])
}

func JoinHandler(event *irc.Event) {
	if event.Command != "JOIN" {
		return
	}
	event.Client.SendRawf("NAMES %s", event.Arguments[0])
}

func QuitHandler(event *irc.Event) {
	if event.Command != "QUIT" {
		return
	}
	event.Client.SendRawf("NAMES %s", event.Arguments[0])
}

func ModeHandler(event *irc.Event) {
	if event.Command != "MODE" {
		return
	}
	event.Client.SendRawf("NAMES %s", event.Arguments[0])
}

func URLHandler(event *irc.Event) {
	if event.Command != "PRIVMSG" {
		return
	}

	// http://blog.mattheworiordan.com/post/13174566389/url-regular-expression-for-links-with-or-without-the
	message_RE := regexp.MustCompile(`((([A-Za-z]{3,9}:(?:\/\/)?)(?:[\-;:&=\+\$,\w]+@)?[A-Za-z0-9\.\-]+|(?:www\.|[\-;:&=\+\$,\w]+@)[A-Za-z0-9\.\-]+)((?:\/[\+~%\/\.\w\-_]*)?\??(?:[\-\+=&;%@\.\w_]*)#?(?:[\.\!\/\\\w]*))?)`)
	matches := message_RE.FindStringSubmatch(event.Arguments[1])

	if len(matches) < 1 {
		return
	}

	url := matches[0]
    go robutdb.SaveURL(url)
}

