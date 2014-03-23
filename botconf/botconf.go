
package botconf

import (
	"log"
	"os"
)

type Config struct {
	Botname		string
	Nickname	string
	Server		string
	Channels	[]string
}

func ReadConfig() *Config {
	botname := os.Getenv("BOTNAME")
	nickname := os.Getenv("BOTNAME")
	if botname == "" {
		log.Fatal("Missing BOTNAME env var")
	}

	server := os.Getenv("IRC_ADDRESS")

	// TODO: Allow multiple channels
	channels := make([]string, 1)
	channels[0] = os.Getenv("IRC_CHANNEL")

	config := Config{botname, nickname, server, channels}
	return &config
}
