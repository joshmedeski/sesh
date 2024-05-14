package configurator

import (
	"fmt"

	"github.com/joshmedeski/sesh/model"
	"github.com/joshmedeski/sesh/oswrap"
	"github.com/joshmedeski/sesh/pathwrap"
	"github.com/joshmedeski/sesh/runtimewrap"
	"github.com/pelletier/go-toml/v2"
)

type Configurator interface {
	GetConfig() (model.Config, error)
}

type RealConfigurator struct {
	os      oswrap.Os
	path    pathwrap.Path
	runtime runtimewrap.Runtime
}

func NewConfigurator(os oswrap.Os, path pathwrap.Path, runtime runtimewrap.Runtime) Configurator {
	return &RealConfigurator{os, path, runtime}
}

func (c *RealConfigurator) configFilePath(rootDir string) string {
	return c.path.Join(rootDir, "sesh", "sesh.toml")
}

func (c *RealConfigurator) getConfigFileFromUserConfigDir() (model.Config, error) {
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

func (c *RealConfigurator) GetConfig() (model.Config, error) {
	config, err := c.getConfigFileFromUserConfigDir()
	if err != nil {
		return model.Config{}, err
	}
	return config, nil
}
