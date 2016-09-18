package handlers

import (
	"encoding/json"
	"github.com/chooper/gobut/robutdb"
	irc "github.com/mikeclarke/go-irclib"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strings"
)

type ChannelState struct {
	Server  string
	Channel string
	Users   []string
}

type ServerState struct {
	Server   string
	Channels map[string]ChannelState
}

var Servers = make(map[string]*ServerState)

type Urinfo struct {
	Uri     string            `json:"uri"`
	Title   string            `json:"title"`
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

func NamesHandler(event *irc.Event) {
	if event.Command != "353" {
		return
	}

	server := event.Client.Server
	irc_chan := event.Arguments[2]
	users := strings.Fields(event.Arguments[3]) // TODO: Track modes

	var server_state *ServerState
	var ok bool
	if server_state, ok = Servers[server]; !ok {
		server_state = new(ServerState)
		server_state.Server = server
		server_state.Channels = make(map[string]ChannelState)
	}

	server_state.Channels[irc_chan] = ChannelState{server, irc_chan, users}
	Servers[server] = server_state

	// HACK: make gobut op everyone
	for _, user := range users {
		if string(user[0]) == "@" { // skip already-opped people
			continue
		}
		if string(user[0]) == "+" { // user is voiced
			user = user[1:]
		}
		nick := strings.Split(user, "!")[0]
		event.Client.SendRawf("MODE %s +o %s", irc_chan, nick)
	}
}

func PartHandler(event *irc.Event) {
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

func AutoOpHandler(event *irc.Event) {
	if event.Command != "JOIN" {
		return
	}
	nick := strings.Split(event.Prefix, "!")[0]
	event.Client.SendRawf("MODE %s +o %s", event.Arguments[0], nick)
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

func FrostDateHandler(event *irc.Event) {
	if event.Command != "PRIVMSG" {
		return
	}

	args := strings.Fields(event.Arguments[1])

	if args[0] != ".frost" {
		return
	}

	zip := args[1]
	uri := "http://www.almanac.com/gardening/frostdates/zipcode/" + zip

	log.Printf("FrostDateHandler invoked. Zip: %s", zip)

	response, err := http.Get(uri)
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
	response_RE := regexp.MustCompile(`<table><tr><th>Climate Station<\/th><th>Last Spring Frost \(50% Probability\)<\/th><th>First Fall Frost \(50% Probability\)<\/th><th>Growing Season<\/th><\/tr><tr><td>([^<]+)<\/td><td>([^<]+)<\/td><td>([^<]+)<\/td><td>([^<]+)<\/td><\/tr><\/table>`)
	matches := response_RE.FindStringSubmatch(string(body))
	if len(matches) < 1 {
		return
	}

	log.Printf("FrostDateHandler matches: %q", matches)
	event.Client.Privmsgf(event.Arguments[0], "Climate station: %s / First fall frost: %s / Last spring frost: %s / Growing season: %s days", matches[1], matches[3], matches[2], matches[4])
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

	base_uri := os.Getenv("URINFO_API")
	if base_uri == "" {
		return
	}
	api_uri, err := url.Parse(base_uri)
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

func TopSharersHandler(event *irc.Event) {
	if event.Command != "PRIVMSG" {
		return
	}

	message_RE := regexp.MustCompile(`^\.top\s*$`)
	matches := message_RE.FindStringSubmatch(event.Arguments[1])

	if len(matches) < 1 {
		return
	}

	// TODO Hand this off to a goroutine
	stats, err := robutdb.TopSharers()
	if err != nil {
		log.Print(err)
		return
	}

	event.Client.Privmsg(event.Arguments[0], "Top 5 URL sharers for the past week")
	for k := range stats {
		event.Client.Privmsgf(event.Arguments[0], "%s: %d urls", k, stats[k])
	}
}

func RandomURLHandler(event *irc.Event) {
	if event.Command != "PRIVMSG" {
		return
	}

	message_RE := regexp.MustCompile(`^\.random\s*$`)
	matches := message_RE.FindStringSubmatch(event.Arguments[1])

	if len(matches) < 1 {
		return
	}

	url, err := robutdb.RandomURL()
	if err != nil {
		log.Print(err)
		return
	}
	event.Client.Privmsg(event.Arguments[0], url)
}

func SearchURLHandler(event *irc.Event) {
	if event.Command != "PRIVMSG" {
		return
	}

	message_RE := regexp.MustCompile(`^\.search\s*(.*)$`)
	matches := message_RE.FindStringSubmatch(event.Arguments[1])

	if len(matches) < 1 {
		return
	}

	query := strings.Join(strings.Fields(matches[0])[1:], " ")

	url, err := robutdb.SearchURL(query)
	if err != nil {
		log.Print(err)
		event.Client.Privmsgf(event.Arguments[0], "I could not find any URLs matching '%s' this time", query)
		return
	}
	event.Client.Privmsg(event.Arguments[0], url)
}

func CountURLsHandler(event *irc.Event) {
	if event.Command != "PRIVMSG" {
		return
	}

	message_RE := regexp.MustCompile(`^\.stats\s*$`)
	matches := message_RE.FindStringSubmatch(event.Arguments[1])

	if len(matches) < 1 {
		return
	}

	// TODO Hand this off to a goroutine
	count, err := robutdb.CountURLs()
	if err != nil {
		log.Print(err)
		return
	}

	event.Client.Privmsgf(event.Arguments[0], "There are %d unique URLs", count)
}
