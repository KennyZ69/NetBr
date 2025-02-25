package pkg

import (
	"encoding/json"
	"fmt"
	"net"
	"os"
	"os/user"
	"path/filepath"
)

// Config holds the user's data used by Kenny's networking tools
type Config struct {
	// NetIfi  string     `json:"interface"` // local network interface

	NetIfi  *net.Interface `json:"interface"` // local network interface
	LocalIP string         `json:"local_ip"`  // users (attackers) local IP
	Mac     string         `json:"mac"`       // users (attackers) actual MAC address
	CIDR    *net.IPNet     `json:"cidr"`      // the IP range on network
	Gateway net.IP         `json:"gateway"`   // the gateway IP

	// may add more fields throughout the development process

	path string
}

func confPath() string {
	usr, err := user.Current()
	// if the current user throws an error, it falls back to a local dir config file
	if err != nil {
		return "./config.json"
	}
	// else It creates the path for directory of the current app and its files
	return filepath.Join(usr.HomeDir, ".config", "netBr", "config.json")
}

func SaveConf(conf *Config) error {
	path := confPath()
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}

	data, err := json.Marshal(&conf)
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0644)
}

func LoadConf() (*Config, error) {
	path := confPath()

	_, err := os.Stat(path)
	// fmt.Println("Config exists: ", err == nil)
	if err != nil {
		return nil, fmt.Errorf("Config file does not exist\n")
	}

	data, err := ReadFile(path)
	if err != nil {
		return nil, err
	}

	var conf Config
	if err = json.Unmarshal(data, &conf); err != nil {
		return nil, err
	}

	return &conf, nil
}
