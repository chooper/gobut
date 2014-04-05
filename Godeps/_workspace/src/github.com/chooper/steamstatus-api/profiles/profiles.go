
package profiles

import (
    "encoding/json"
    "io/ioutil"
    "log"
    "net/http"
    "regexp"
    "time"
)

type ProfileData struct {
    Url             string  `json:"url"`
    SteamId         string  `json:"steamid"`
    PersonaName     string  `json:"personaname"`
    Summary         string  `json:"summary"`
    InGame          string  `json:"ingame"`
}

type ProfileResponse struct {
    Profile         ProfileData
    Error           error
}

func ParseProfile(response_body *string) (ProfileData, error) {
    var profile ProfileData
    json_regex := regexp.MustCompile(`g_rgProfileData = (.*);`)
    json_matches := json_regex.FindStringSubmatch(*response_body)

    if len(json_matches) > 0 {
        if err := json.Unmarshal([]byte(json_matches[1]), &profile); err != nil {
            return ProfileData{}, err
        }
    }

    // Find out if user is in a game
    ingame_regex := regexp.MustCompile(`<div class="profile_in_game_header">(.*)</div>`)
    ingame_matches := ingame_regex.FindStringSubmatch(*response_body)

    var ingame bool = false
    if len(ingame_matches) > 0 && ingame_matches[1] == "Currently In-Game" {
        ingame = true

        // Find out which game
        gamename_regex := regexp.MustCompile(`<div class="profile_in_game_name">(.*)</div>`)
        gamename_matches := gamename_regex.FindStringSubmatch(*response_body)

        // Add the game name to ProfileData
        if ingame && len(gamename_matches) > 0 {
            profile.InGame = gamename_matches[1]
        }
    }
    return profile, nil
}

func GetProfile(username string) (ProfileData, error) {
    // Download the profile from steam
    profile_url := "http://steamcommunity.com/id/" + username + "/"
    response, err := http.Get(profile_url)
    defer response.Body.Close()
    if err != nil {
        return ProfileData{}, err
    }
    body, err := ioutil.ReadAll(response.Body)
    if err != nil {
        return ProfileData{}, err
    }
    response_body := string(body)

    // Parse profile data
    return ParseProfile(&response_body)
}

func FetchProfiles(usernames []string) []ProfileData {
    var profile_count int = len(usernames)
    var profiles = make([]ProfileData, profile_count)
    var profile_c = make(chan ProfileResponse)

    // Request multiple users at once
    for _, username := range usernames {
        go FetchProfile(username, profile_c)
    }

    // Wait for responses
    timeout := time.After(10 * time.Second)
    for idx := 0; idx < profile_count; idx++ {
        select {
        case response := <- profile_c:
            if response.Error == nil {
                profiles[idx] = response.Profile
            }
        case <- timeout:
            log.Print("FetchProfiles timed out!")
            break
        }
    }
    return profiles
}

func FetchProfile(username string, c chan ProfileResponse) {
    fanout_c := make(chan ProfileResponse, 1)
    for i := 0; i < 3; i++ {
        go func() {
            p, err := GetProfile(username)
            response := ProfileResponse{
                Profile: p,
                Error:  err,
            }
            select {
                case fanout_c <- response:
                default:
            }
        }()
    }
    select {
    case c <- <- fanout_c:
    }
}

