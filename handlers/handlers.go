// handlers -- Copyright (c) 2013 Charles Hooper

package handlers

import (
	irc "github.com/mikeclarke/go-irclib"
	"fmt"
	"regexp"
	"strings"
)

func DebugHandler(event *irc.Event) {
	message_RE := regexp.MustCompile(`(.*)`)
	matches := message_RE.FindStringSubmatch(event.Raw)
	
	if len(matches) < 1 {
		return
	}

	fmt.Printf("-> %q", event)
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
