package main

import (
	"flag"
	"github.com/chooper/gobut/botconf"
	"github.com/chooper/gobut/handlers"
	sp "github.com/chooper/gobut/steam-poller"
	"github.com/mikeclarke/go-irclib"
	"log"
	"time"
)

func main() {
	flag.Parse()

	// Read the config
	config := botconf.ReadConfig()

	// Set up IRC client
	client := irc.New(config.Nickname, config.Botname)

	// Register handlers
	client.AddHandler(handlers.EchoHandler)
	client.AddHandler(handlers.NamesHandler)
	client.AddHandler(handlers.PartHandler)
	client.AddHandler(handlers.QuitHandler)
	client.AddHandler(handlers.JoinHandler)
	client.AddHandler(handlers.AutoOpHandler)
	client.AddHandler(handlers.URLHandler)
	client.AddHandler(handlers.TopSharersHandler)
	client.AddHandler(handlers.RandomURLHandler)
	client.AddHandler(handlers.SearchURLHandler)
	client.AddHandler(handlers.CountURLsHandler)
	client.AddHandler(handlers.FrostDateHandler)

	// Connect to server
	err := client.Connect(config.Server)
	if err != nil {
		// TODO: Don't crash - recover and connect to new server
		log.Fatalf("Error connecting to server %q: %s\n", config.Server, err)
	}

	// Join channels
	var irc_chan string
	go func() {
		for {
			for _, irc_chan = range config.Channels {
				log.Printf("%s: Joining channel %q\n", config.Botname, irc_chan)
				client.Join(irc_chan)
				client.SendRawf("NAMES %s", irc_chan)
			}
			time.Sleep(time.Duration(10) * time.Second)
		}
	}()

	// Set up steam poller for each channel
	for _, irc_chan = range config.Channels {
		log.Printf("%s: Setting up steam poller for channel %q\n", config.Botname, irc_chan)
		sp.RunPoller(client, irc_chan)
	}

	// Run loop
	client.Run()
}
