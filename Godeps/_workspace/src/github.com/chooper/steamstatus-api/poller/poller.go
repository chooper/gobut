
package poller

import (
    //"log"
    "github.com/chooper/steamstatus-api/client"
    "os"
    "time"
)

// This poller will periodically poll Steam for each given player's status
// and send notifications on its `NotifyChan` whenever the player enters a
// new game.
//
// Note that notifications are NOT sent when the player leaves a game,
// although that state change IS tracked.

type KnownState map[string]string

type StateChange struct {
    Who         string
    OldState    string
    NewState    string
}
    
type Notification struct {
    Changes     []StateChange
}

type Poller struct {
    Usernames   []string
    NotifyChan  chan Notification
}

func (p Poller) Loop() {
    state := make(KnownState)
    api := client.Connect(os.Getenv("STEAMSTATUS_API")) // TODO handle if not exists

    for {
        var notification Notification
        // TODO: Query API instead of requesting this directly?
        //profiles := profiles.FetchProfiles(p.Usernames)
        profiles, err := api.Profiles(p.Usernames)
        if err != nil {
            // TODO log
            time.Sleep(10 * time.Second)
            continue
        }

        for _, profile := range profiles {
            // Don't notify when not in a game (even if player just left one)
            if profile.InGame == "" {
                state[profile.PersonaName] = profile.InGame
                continue
            }

            if profile.InGame != state[profile.PersonaName] {
                // Track the change
                change := StateChange{
                    Who:        profile.PersonaName,
                    OldState:   state[profile.PersonaName],
                    NewState:   profile.InGame,
                }

                // Update state
                state[profile.PersonaName] = profile.InGame

                // Don't enqueue any notifications about leaving a game
                if profile.InGame == "" {
                    continue
                } else {
                    notification.Changes = append(notification.Changes, change)
                }
            }
        }
        p.NotifyChan <- notification
        time.Sleep(10 * time.Second)
    }
}

