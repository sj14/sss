package controller

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
)

func LoadConfig(configPath string) (Config, error) {
	var config Config

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return config, err
	}

	if configPath == "" {
		configPath = filepath.Join(homeDir, ".config", "sss", "config.toml")
	}

	md, err := toml.DecodeFile(configPath, &config)
	if err != nil {
		return config, err
	}

	if undecoded := md.Undecoded(); len(undecoded) > 0 {
		return config, fmt.Errorf("unknown fields in config: %v", undecoded)
	}

	return config, nil
}
