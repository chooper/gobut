
package main

import (
	"github.com/chooper/gobut/botconf"
	"flag"
	"github.com/chooper/gobut/handlers"
	"github.com/mikeclarke/go-irclib"
	"github.com/chooper/steamstatus-api/poller"
	"strings"
	"os"
	"log"
)

func main() {
	flag.Parse()

	// Read the config
	config := botconf.ReadConfig()

	// Set up IRC client
	client := irc.New(config.Nickname, config.Botname)

	// Register handlers
	client.AddHandler(handlers.DebugHandler)
	client.AddHandler(handlers.EchoHandler)
	client.AddHandler(handlers.NamesHandler)
	client.AddHandler(handlers.PartHandler)
	client.AddHandler(handlers.QuitHandler)
	client.AddHandler(handlers.JoinHandler)
	client.AddHandler(handlers.ModeHandler)
	client.AddHandler(handlers.FuckYeahHandler)
	client.AddHandler(handlers.URLHandler)
	
	// Connect to server
	err := client.Connect(config.Server)
	if err != nil {
		// TODO: Don't crash - recover and connect to new server
		log.Fatalf("Error connecting to server %q: %s\n", config.Server, err)
	}
	
	// Join channels
	var irc_chan string
	for _, irc_chan = range config.Channels {
		log.Printf("%s: Joining channel %q\n", config.Botname, irc_chan)
		client.Join(irc_chan)
	}

	// Set up steam poller
	launch_poller(client, irc_chan)

	// Run loop
	client.Run()
}

func poller_listener(c *irc.IRCClient, irc_chan string, ch chan poller.Notification) {
	for {
	    select {
	        case n := <- ch:
	            for _, delta := range n.Changes {
	                c.Privmsgf(irc_chan, "%s has started playing %s", delta.Who, delta.NewState)
	            }
	    }
	}
}

func launch_poller(c *irc.IRCClient, irc_chan string) {
	var usernames []string
	if usernames = strings.Split(os.Getenv("POLL_USERNAMES"), ","); len(usernames) == 0 {
		log.Fatalf("Must set POLL_USERNAMES env var!")
	}

	ch := make(chan poller.Notification)
	p := &poller.Poller{
		Usernames:	  usernames,
		NotifyChan:	 ch,
	}

	// Start the poller
	go p.Loop()
	go poller_listener(c, irc_chan, ch)
}

