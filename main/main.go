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

type Config struct {
	botname		string
	nickname	string
	servers		[]string
	channels	[]string
}

func readConfig(conf_file string) *Config {
	file, err := yaml.ReadFile(conf_file);
	if err != nil {
		log.Fatalf("readConfig(%q): %s", file, err)
	}
	//return config;
	botname, _ := file.Get("global.botname")

	network_config, _ := yaml.Child(file.Root, ".network")
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
	
	config := Config{botname, nickname, servers, channels}
	return &config
}

func main() {
	flag.Parse()

	// Read the config
	config := readConfig("config.yml")

	// Set up IRC client
	client := irc.New(config.nickname, config.botname)
	// Register handlers
	client.AddHandler(handlers.DebugHandler)
	client.AddHandler(handlers.EchoHandler)
	client.AddHandler(handlers.NamesHandler)
	client.AddHandler(handlers.PartHandler)
	client.AddHandler(handlers.QuitHandler)
	client.AddHandler(handlers.JoinHandler)
	client.AddHandler(handlers.ModeHandler)
	
	// Connect to server
	current_server := config.servers[0] // FIXME: Pick server at random
	err := client.Connect(current_server)
	if err != nil {
		// TODO: Don't crash - recover and connect to new server
		log.Fatalf("Error connecting to server %q: %s\n", current_server, err)
	}
	
	// Join channels
	for _, irc_chan := range config.channels {
		fmt.Printf("%s: Joining channel %q\n", config.botname, irc_chan)
		client.Join(irc_chan)
	}
	
	// Run loop
	client.Run()
}
