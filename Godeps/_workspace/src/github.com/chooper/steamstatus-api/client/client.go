
package client

import (
    "encoding/json"
    "io/ioutil"
    "log"
    "net/http"
    "github.com/chooper/steamstatus-api/profiles"
    "strings"
)

type Api struct {
    Uri     string
}

func Connect(uri string) (*Api) {
    return &Api{
        Uri:    uri,
    }
}

func (api *Api) Profiles(usernames []string) ([]profiles.ProfileData, error) {
    uri := api.Uri + "/status?usernames=" + strings.Join(usernames, ",")

    response, err := http.Get(uri)
    if err != nil {
        log.Print(err)
        return []profiles.ProfileData{}, err
    }

    defer response.Body.Close()
    body, err := ioutil.ReadAll(response.Body)
    if err != nil {
        log.Print(err)
        return []profiles.ProfileData{}, err
    }

    var profs []profiles.ProfileData
    if err := json.Unmarshal([]byte(body), &profs); err != nil {
        log.Print(err)
        return []profiles.ProfileData{}, err
    }

    return profs, nil
}

