// handlers -- Copyright (c) 2013 Charles Hooper

package handlers

import (
	irc "github.com/mikeclarke/go-irclib"
	"log"
	"regexp"
	"strings"
)

func DebugHandler(event *irc.Event) {
	log.Println("-> ", event)
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
	// TODO: Update state
	log.Println("XXX In channel = ", event.Arguments[3:])
}

func PartHandler(event *irc.Event){ 
	if event.Command != "PART" {
		return
	}
	// TODO: Update state
	log.Println("User left channel = ", event.Arguments)
}

func JoinHandler(event *irc.Event) {
	if event.Command != "JOIN" {
		return
	}
	// TODO: Update state
	log.Println("User joined channel = ", event.Arguments)
}

func QuitHandler(event *irc.Event) {
	if event.Command != "QUIT" {
		return
	}
	// TODO: Update state
	log.Println("User quit = ", event.Arguments)
}

func ModeHandler(event *irc.Event) {
	if event.Command != "MODE" {
		return
	}
	// TODO: Update state
	log.Println("Mode change = ", event.Arguments)
}
