
package handlers

import (
	"encoding/json"
	"io/ioutil"
	irc "github.com/mikeclarke/go-irclib"
	"log"
	"net/http"
	"net/url"
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

type Urinfo struct {
	Uri	string `json:"uri"`
	Title string `json:"title"`
	Headers map[string]string `json:"headers"`
}

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

func FuckYeahHandler(event *irc.Event) {
    if event.Command != "PRIVMSG" {
        return
    }

    message_RE := regexp.MustCompile(`fuck yeah ([\w\d]+)`)
    matches := message_RE.FindStringSubmatch(event.Arguments[1])

    if len(matches) < 1 {
        return
    }

    phrase := matches[0]
    uri := "http://fuckyeah.herokuapp.com/" + phrase + ".jpg"

    // let app try to cache the result
    cache := func() { r, _ := http.Get(uri); r.Body.Close() }
    go cache()

    event.Client.Privmsg(event.Arguments[0], uri)
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

	target_uri := matches[0]

	// TODO: Don't hardcode this
	api_uri, err := url.Parse("http://urinfo.herokuapp.com/fetch")
	if err != nil {
		log.Print(err)
	}
	
	api_query := api_uri.Query()
	api_query.Set("uri", target_uri)
	api_uri.RawQuery = api_query.Encode()

	response, err := http.Get(api_uri.String())
	if err != nil {
		log.Print(err)
		return
	}

	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Print(err)
		return
	}

	var info Urinfo
	if err = json.Unmarshal([]byte(body), &info); err != nil {
		log.Print(err)
		return
	} 

	go robutdb.SaveURL(info.Uri, info.Title, event.Prefix)
}

