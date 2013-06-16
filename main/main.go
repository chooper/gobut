// bot -- Copyright (c) 2013 Charles Hooper

package main

import (
	"flag"
	"fmt"
	handlers "github.com/chooper/go-bot/handlers"
	"github.com/mikeclarke/go-irclib"
	"log"
	"github.com/kylelemons/go-gypsy/yaml"
)

func readConfig(conf_file string) *yaml.File {
	config, err := yaml.ReadFile(conf_file);
	if err != nil {
		log.Fatalf("readConfig(%q): %s", conf_file, err)
	}
	return config;
}

func main() {
	flag.Parse()

	// Read the config
	// TODO: Read whole config into a Map we can pass around
	config := readConfig("config.yml")
	botname, _ := config.Get("global.botname")

	network_config, _ := yaml.Child(config.Root, ".network")
	nickname_node, _ := yaml.Child(network_config, ".nick")
	servers_node, _ := yaml.Child(network_config, ".servers")
	channels_node, _ := yaml.Child(network_config, ".channels")
	
	nickname := nickname_node.(yaml.Scalar).String()
	servers_node_list := servers_node.(yaml.List)
	channels_node_list := channels_node.(yaml.List)

	// Config: Copy servers into an arrray
	servers_len := servers_node_list.Len()
	servers := make([]string, servers_len)
	for idx := range servers {
		// God this is ugly; am I even doing this right?
		servers[idx] = servers_node_list.Item(idx).(yaml.Scalar).String()
	}
		
	// Config: Copy channels into an array
	channels_len := channels_node_list.Len()
	channels := make([]string, channels_len)
	for idx := range channels {
		channels[idx] = channels_node_list[idx].(yaml.Scalar).String()
	}
	// End config reading

	// Set up IRC client
	client := irc.New(nickname, botname)
	// Register handlers
	client.AddHandler(handlers.DebugHandler)
	client.AddHandler(handlers.EchoHandler)
	client.AddHandler(handlers.NamesHandler)
	client.AddHandler(handlers.PartHandler)
	client.AddHandler(handlers.QuitHandler)
	client.AddHandler(handlers.JoinHandler)
	client.AddHandler(handlers.ModeHandler)
	
	// Connect to server
	current_server := servers[0] // FIXME: Pick server at random
	err := client.Connect(current_server)
	if err != nil {
		// TODO: Don't crash - recover and connect to new server
		log.Fatalf("Error connecting to server %q: %s\n", current_server, err)
	}
	
	// Join channels
	for _, irc_chan := range channels {
		fmt.Printf("%s: Joining channel %q\n", botname, irc_chan)
		client.Join(irc_chan)
	}
	
	// Run loop
	client.Run()
}
