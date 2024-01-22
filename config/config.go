package config

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"runtime"

	"github.com/pelletier/go-toml/v2"
)

type (
	// SessionName struct {
	// 	PrependParentDir bool `toml:"prepend_parent_dir"`
	// }
	//
	// Session struct {
	// 	Name SessionName
	// }
	Script struct {
		SessionPath string `toml:"session_path"`
		ScriptPath  string `toml:"script_path"`
	}
	Config struct {
		// Session              Session
		StartupScripts       []Script `toml:"startup_scripts"`
		DefaultStartupScript string   `toml:"default_startup_script"`
	}
)

func getUserConfigDir() (string, error) {
	switch runtime.GOOS {
	case "darwin":
		// typically ~/Library/Application Support, but we want to use ~/.config
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		return path.Join(homeDir, ".config"), nil
	default:
		return os.UserConfigDir()
	}
}

func ParseConfigFile() Config {
	config := Config{}
	configDir, err := getUserConfigDir()
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
