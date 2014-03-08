
package botconf

import (
	"log"
	"github.com/kylelemons/go-gypsy/yaml"
)

type Config struct {
	Botname		string
	Nickname	string
	Servers		[]string
	Channels	[]string
}

func ReadConfig(conf_file string) *Config {
	file, err := yaml.ReadFile(conf_file);
	if err != nil {
		log.Fatalf("readConfig(%q): %s", file, err)
	}
	botname, _ := file.Get("global.botname")
	network_config, _ := yaml.Child(file.Root, ".network")
	nickname_node, _ := yaml.Child(network_config, ".nick")
	servers_node, _ := yaml.Child(network_config, ".servers")
	channels_node, _ := yaml.Child(network_config, ".channels")
	
	nickname := nickname_node.(yaml.Scalar).String()
	servers_node_list := servers_node.(yaml.List)
	channels_node_list := channels_node.(yaml.List)

	// Copy servers into an arrray
	servers_len := servers_node_list.Len()
	servers := make([]string, servers_len)
	for idx := range servers {
		// God this is ugly; am I even doing this right?
		servers[idx] = servers_node_list.Item(idx).(yaml.Scalar).String()
	}
		
	// Copy channels into an array
	channels_len := channels_node_list.Len()
	channels := make([]string, channels_len)
	for idx := range channels {
		channels[idx] = channels_node_list[idx].(yaml.Scalar).String()
	}
	
	config := Config{botname, nickname, servers, channels}
	return &config
}
