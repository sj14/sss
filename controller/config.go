package controller

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
)

func DefaultConfigPath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(homeDir, ".config", "sss", "config.toml"), nil
}

func LoadConfig(configPath string) (Config, error) {
	var (
		config Config
		err    error
	)

	if configPath == "" {
		configPath, err = DefaultConfigPath()
		if err != nil {
			return config, err
		}

		// prevent failing when the default config does not exist
		_, err := os.Stat(configPath)
		if os.IsNotExist(err) {
			return config, nil
		}
		if err != nil {
			return config, err
		}
	}

	// TODO: check file permissions (e.g. is it group readable, similar as ssh is doing)

	md, err := toml.DecodeFile(configPath, &config)
	if err != nil {
		return config, err
	}

	if undecoded := md.Undecoded(); len(undecoded) > 0 {
		return config, fmt.Errorf("unknown fields in config: %v", undecoded)
	}

	return config, nil
}
