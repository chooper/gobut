package poller

import (
    //"log"
    "github.com/chooper/mcstatus-client-go/client"
    "os"
    "time"
)

// This poller will periodically poll Minecraft for the list of online players
// and send notifications on its `NotifyChan` whenever the player enters or leaves
// a new game.

type StateChange struct {
    Who         string
    Direction   string
}
    
type Notification struct {
    Changes     []StateChange
}

type Poller struct {
    Server      string
    NotifyChan  chan Notification
}

func contains(haystack client.PlayerList, needle string) bool {
    for _, a := range haystack { if a == needle { return true } }
    return false
}

func (p Poller) Loop() {
    known_players := make(client.PlayerList, 0)
    api := client.Connect(os.Getenv("MCSTATUS_API")) // TODO handle if not exists

    for {
        var notification Notification
        players, err := api.PlayersOnline(p.Server)
        if err != nil {
            // TODO log
            time.Sleep(10 * time.Second)
            continue
        }

        // Check for players who left
        for _, player := range known_players {
            if !contains(players, player) {
                // Player left the game
                change := StateChange{
                    Who:        player,
                    Direction:  "left",
                }
                notification.Changes = append(notification.Changes, change)
            }
        }

        // Check for players who've joined
        for _, player := range players {
            if !contains(known_players, player) {
                // Player joined the game
                change := StateChange{
                    Who:        player,
                    Direction:  "joined",
                }
                notification.Changes = append(notification.Changes, change)
            }
        }

        // Update known state for next iteration
        known_players = players

        p.NotifyChan <- notification
        time.Sleep(10 * time.Second)
    }
}

