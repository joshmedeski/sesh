package config

import (
	"fmt"

	"github.com/joshmedeski/sesh/model"
	"github.com/joshmedeski/sesh/oswrap"
	"github.com/joshmedeski/sesh/pathwrap"
	"github.com/joshmedeski/sesh/runtimewrap"
	"github.com/pelletier/go-toml/v2"
)

type Config interface {
	GetConfig() (model.Config, error)
}

type RealConfig struct {
	os      oswrap.Os
	path    pathwrap.Path
	runtime runtimewrap.Runtime
}

func NewConfig(os oswrap.Os, path pathwrap.Path, runtime runtimewrap.Runtime) Config {
	return &RealConfig{os, path, runtime}
}

func (c *RealConfig) configFilePath(rootDir string) string {
	return c.path.Join(rootDir, "sesh", "sesh.toml")
}

func (c *RealConfig) getConfigFileFromUserConfigDir() (model.Config, error) {
	config := model.Config{}

	userConfigDir, err := c.os.UserConfigDir()
	if err != nil {
		return config, fmt.Errorf("couldn't get user config dir: %q", err)
	}
	configFilePath := c.configFilePath(userConfigDir)
	file, err := c.os.ReadFile(configFilePath)
	if err != nil {
		return config, fmt.Errorf("couldn't read config file: %q", err)
	}
	err = toml.Unmarshal(file, &config)
	if err != nil {
		return config, fmt.Errorf("couldn't unmarshal config file: %q", err)
	}
	return config, nil

	// TODO: look for config file in `~/.config`

	// switch c.runtime.GOOS() {
	// case "darwin":
	// 	// TODO: support both
	// 	// typically ~/Library/Application Support, but we want to use ~/.config
	// 	homeDir, err := os.UserHomeDir()
	// 	if err != nil {
	// 		return model.Config{}, err
	// 	}
	// 	return path.Join(homeDir, ".config"), nil
	// default:
	// 	return os.UserConfigDir()
	// }
}

func (c *RealConfig) GetConfig() (model.Config, error) {
	config, err := c.getConfigFileFromUserConfigDir()
	if err != nil {
		return model.Config{}, err
	}
	return config, nil
}
