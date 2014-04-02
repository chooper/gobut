
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
	var usernames []string
	if usernames = strings.Split(os.Getenv("POLL_USERNAMES"), ","); len(usernames) == 0 {
		log.Fatalf("Must set POLL_USERNAMES env var!")
	}

	c := make(chan poller.Notification)
	p := &poller.Poller{
		Usernames:	  usernames,
		NotifyChan:	 c,
	}

	// Start the poller
	go p.Loop()

	// Start the poller listener
	go func() {
		for {
			select {
				case n := <- c:
					for _, delta := range n.Changes {
						client.Privmsgf(irc_chan, "%s has started playing %s", delta.Who, delta.NewState)
					}
			}
		}
	}()

	// Run loop
	client.Run()
}

