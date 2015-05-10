
package sp

import (
    "github.com/mikeclarke/go-irclib"
    "github.com/chooper/mcstatus-client-go/poller"
    "os"
)

func poller_listener(c *irc.IRCClient, irc_chan string, ch chan poller.Notification) {
    for {
        select {
            case n := <- ch:
                for _, delta := range n.Changes {
                    c.Privmsgf(irc_chan, "%s has %s playing Minecraft on %s", delta.Who, delta.Direction, os.Getenv("MINECRAFT_SERVER"))
                }
        }
    }
}

func RunPoller(c *irc.IRCClient, irc_chan string) {
    var server string
    server = os.Getenv("MINECRAFT_SERVER")

    ch := make(chan poller.Notification)
    p := &poller.Poller{
        Server:         server,
        NotifyChan:     ch,
    }

    // Start the poller
    go p.Loop()
    go poller_listener(c, irc_chan, ch)
}
