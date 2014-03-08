
package main

import (
	"github.com/chooper/go-bot/botconf"
	"flag"
	"github.com/chooper/go-bot/handlers"
	"github.com/mikeclarke/go-irclib"
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
	client.AddHandler(handlers.URLHandler)
	
	// Connect to server
	current_server := config.Servers[0] // FIXME: Pick server at random
	err := client.Connect(current_server)
	if err != nil {
		// TODO: Don't crash - recover and connect to new server
		log.Fatalf("Error connecting to server %q: %s\n", current_server, err)
	}
	
	// Join channels
	for _, irc_chan := range config.Channels {
		log.Printf("%s: Joining channel %q\n", config.Botname, irc_chan)
		client.Join(irc_chan)
	}
	
	// Run loop
	client.Run()
}
