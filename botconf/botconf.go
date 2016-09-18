package botconf

import (
	"log"
	"os"
	"strings"
)

type Config struct {
	Botname  string
	Nickname string
	Server   string
	Channels []string
}

func ReadConfig() *Config {
	botname := os.Getenv("BOTNAME")
	nickname := os.Getenv("BOTNAME")
	if botname == "" {
		log.Fatal("Missing BOTNAME env var")
	}

	server := os.Getenv("IRC_ADDRESS")
	channels := strings.Split(os.Getenv("IRC_CHANNEL"), ",")

	config := Config{botname, nickname, server, channels}
	return &config
}
