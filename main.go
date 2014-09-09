
package main

import (
    "github.com/chooper/gobut/botconf"
    "flag"
    "github.com/chooper/gobut/handlers"
    sp "github.com/chooper/gobut/steam-poller"
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
    client.AddHandler(handlers.FuckYeahHandler)
    client.AddHandler(handlers.URLHandler)
    client.AddHandler(handlers.TopSharersHandler)
    client.AddHandler(handlers.RandomURLHandler)
    client.AddHandler(handlers.SearchURLHandler)
    client.AddHandler(handlers.CountURLsHandler)
    
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
    sp.RunPoller(client, irc_chan)

    // Run loop
    client.Run()
}

