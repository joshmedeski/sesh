package config

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"runtime"

	"github.com/joshmedeski/sesh/dir"
	"github.com/pelletier/go-toml/v2"
)

type (
	Script struct {
		SessionPath string `toml:"session_path"`
		ScriptPath  string `toml:"script_path"`
	}
	Config struct {
		ImportPaths          []string `toml:"import"`
		StartupScripts       []Script `toml:"startup_scripts"`
		DefaultStartupScript string   `toml:"default_startup_script"`
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

	if len(config.ImportPaths) > 0 {
		for _, path := range config.ImportPaths {
			importConfig := Config{}
			importConfigPath := dir.FullPath(path)
			if err := parseConfigFromFile(importConfigPath, &importConfig); err != nil {
				return fmt.Errorf("Error parsing import config file: %s", err)
			}
			config.StartupScripts = append(config.StartupScripts, importConfig.StartupScripts...)
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
