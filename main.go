package main

import (
	"flag"
	"github.com/chooper/go-irclib"
	"github.com/chooper/gobut/botconf"
	"github.com/chooper/gobut/handlers"
	"log"
)

func main() {
	flag.Parse()

	// Read the config
	config := botconf.ReadConfig()

	// Set up IRC client
	client := irc.New(config.Nickname, config.Botname)

	// Register handlers
	client.AddHandler(handlers.RegistrationHandler)
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
		log.Printf("Could not connect to server %q: %s\n", config.Server, err)
		client.Reconnect()
	}

	// Run main loop
	client.Run()
}
