
package sp

import (
    "github.com/mikeclarke/go-irclib"
	"github.com/chooper/steamstatus-api/poller"
	"strings"
	"os"
	"log"
)

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

func RunPoller(c *irc.IRCClient, irc_chan string) {
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

