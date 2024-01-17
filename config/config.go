package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/pelletier/go-toml/v2"
)

type SessionName struct {
	IncludeRootDir bool `toml:"include_root_dir"`
}
type Session struct {
	Name SessionName
}
type Config struct {
	Session Session
}

func ParseConfigFile() Config {
	configDir, err := os.UserConfigDir()
	config := Config{}
	if err != nil {
		fmt.Printf(
			"Error determining the user config directory: %s\nUsing default config instead",
			err,
		)
		return config
	}
	configPath := filepath.Join(configDir, "sesh", "sesh.toml")
	data, err := os.ReadFile(configPath)
	if err != nil {
		return config
	}
	err = toml.Unmarshal(data, &config)
	if err != nil {
		fmt.Printf(
			"Error parsing config file: %s\nUsing default config instead",
			err,
		)
		return config
	}
	return config
}
