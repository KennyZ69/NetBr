package pkg

import (
	"encoding/json"
	"os"
	"os/user"
	"path/filepath"
)

// config holds the user's data used by Kenny's networking tools
type Config struct {
	NetIfi  string `json:"interface"`
	LocalIP string `json:"local_ip"`
	Mac     string `json:"mac"`
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

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var conf Config
	if err = json.Unmarshal(data, &conf); err != nil {
		return nil, err
	}

	return &conf, err
}
