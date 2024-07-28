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
	DefaultSessionConfig struct {
		StartupScript  string `toml:"startup_script"`
		StartupCommand string `toml:"startup_command"`
		Tmuxp          string `toml:"tmuxp"`
		Tmuxinator     string `toml:"tmuxinator"`
	}

	SessionConfig struct {
		Name     string   `toml:"name"`
		Path     string   `toml:"path"`
		PathList []string `toml:"path_list"`
		DefaultSessionConfig
	}

	Config struct {
		ImportPaths          []string             `toml:"import"`
		DefaultSessionConfig DefaultSessionConfig `toml:"default_session"`
		SessionConfigs       []SessionConfig      `toml:"session"`
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
	file, err := os.ReadFile(configPath)
	if err != nil {
		return fmt.Errorf("parseConfigFromFile - error reading config file (%s): %v", configPath, err)
	}
	err = toml.Unmarshal(file, config)
	if err != nil {
		return fmt.Errorf(": %s", err)
	}
	if len(config.ImportPaths) > 0 {
		for _, path := range config.ImportPaths {
			importConfig := Config{}
			importConfigPath := dir.FullPath(path)
			if err := parseConfigFromFile(importConfigPath, &importConfig); err != nil {
				return fmt.Errorf("parse config from import file failed: %s", err)
			}
			config.SessionConfigs = append(config.SessionConfigs, importConfig.SessionConfigs...)
		}
	}
	return nil
}

// TODO: add error handling (return error)
func ParseConfigFile(fetcher ConfigDirectoryFetcher) Config {
	config := Config{}
	configDir, err := fetcher.GetUserConfigDir()
	if err != nil {
		return config
	}
	configPath := filepath.Join(configDir, "sesh", "sesh.toml")
	parseConfigFromFile(configPath, &config)
	return config
}
