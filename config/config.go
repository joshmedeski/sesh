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
	Script struct {
		SessionPath string `toml:"session_path"`
		ScriptPath  string `toml:"script_path"`
	}
	ExtendedConfig struct {
		Path string `toml:"path"`
	}
	Config struct {
		ExtendedConfigs      []ExtendedConfig `toml:"extended_configs"`
		StartupScripts       []Script         `toml:"startup_scripts"`
		DefaultStartupScript string           `toml:"default_startup_script"`
	}
)

type ConfigDirectoryFetcher interface {
	GetUserConfigDir() (string, error)
}

type DefaultConfigDirectoryFetcher struct{}

var _ ConfigDirectoryFetcher = (*DefaultConfigDirectoryFetcher)(nil)

func (d *DefaultConfigDirectoryFetcher) GetUserConfigDir() (string, error) {
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

func parseConfigFromFile(configPath string, config *Config) error {
	data, err := os.ReadFile(configPath)
	if err != nil {
		return fmt.Errorf("Error reading config file: %s", err)
	}
	err = toml.Unmarshal(data, config)
	if err != nil {
		return fmt.Errorf("Error parsing config file: %s", err)
	}

	if len(config.ExtendedConfigs) > 0 {
		for _, item := range config.ExtendedConfigs {
			extendedConfig := Config{}
			if err := parseConfigFromFile(item.Path, &extendedConfig); err != nil {
				return fmt.Errorf("Error parsing extended config file: %s", err)
			}
			config.StartupScripts = append(config.StartupScripts, extendedConfig.StartupScripts...)
		}
	}

	return nil
}

func ParseConfigFile(fetcher ConfigDirectoryFetcher) Config {
	config := Config{}
	configDir, err := fetcher.GetUserConfigDir()
	if err != nil {
		fmt.Printf(
			"Error determining the user config directory: %s\nUsing default config instead",
			err,
		)
		return config
	}
	configPath := filepath.Join(configDir, "sesh", "sesh.toml")

	if err := parseConfigFromFile(configPath, &config); err != nil {
		fmt.Printf(
			"Error parsing config file: %s\nUsing default config instead",
			err,
		)
	}
	return config
}
