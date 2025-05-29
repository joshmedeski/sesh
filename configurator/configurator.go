package configurator

import (
	"errors"
	"fmt"
	"strings"

	"github.com/joshmedeski/sesh/v2/model"
	"github.com/joshmedeski/sesh/v2/oswrap"
	"github.com/joshmedeski/sesh/v2/pathwrap"
	"github.com/joshmedeski/sesh/v2/runtimewrap"
	"github.com/pelletier/go-toml/v2"
)

type Configurator interface {
	GetConfig() (model.Config, string, error)
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

func (c *RealConfigurator) fullImportPath(homeDir, importPath string) (string, error) {
	if !strings.HasPrefix(importPath, "~") {
		return c.path.Abs(importPath)
	}

	return c.path.Join(homeDir, importPath[1:]), nil
}

func (c *RealConfigurator) getConfigFileFromUserConfigDir() (model.Config, string, error) {
	config := model.Config{}
	userHomeDir, err := c.os.UserHomeDir()
	if err != nil {
		return config, "", fmt.Errorf("couldn't get user config dir: %q", err)
	}
	userConfigDir := c.path.Join(userHomeDir, ".config")
	configFilePath := c.configFilePath(userConfigDir)
	file, _ := c.os.ReadFile(configFilePath)
	// TODO: add to debugging logs
	// if err != nil {
	// 	return config, fmt.Errorf("couldn't read config file: %q", err)
	// }
	_ = toml.Unmarshal(file, &config)
	if config.EvalSettings.Strict == "yes" {
		reader := strings.NewReader(string(file))
		d := toml.NewDecoder(reader)
		d.DisallowUnknownFields()
		err = d.Decode(&config)
		var details *toml.StrictMissingError
		if err != nil {
			if !errors.As(err, &details) {
				panic(fmt.Sprintf("err should have been a *toml.StrictMissingError, but got %s (%T)", err, err))
			}
			fmt.Println(details.String())
		}
		if details != nil {
			return config, details.String(), err
		}
	} else {
		err = toml.Unmarshal(file, &config)
		if err != nil {
			var derr *toml.DecodeError
			if errors.As(err, &derr) {
				return config, derr.String(), err
			}
		}
	}

	for _, importPath := range config.ImportPaths {
		importFilePath, err := c.fullImportPath(userHomeDir, importPath)
		if err != nil {
			return config, "", fmt.Errorf("couldn't get full import path: %q", err)
		}

		importFile, err := c.os.ReadFile(importFilePath)
		if err != nil {
			return config, "", fmt.Errorf("couldn't read import file %s: %q", importFilePath, err)
		}

		importConfig := model.Config{}
		if err := toml.Unmarshal(importFile, &importConfig); err != nil {
			return config, "", fmt.Errorf("couldn't unmarshal import file %s: %q", importFilePath, err)
		}

		config.SessionConfigs = append(config.SessionConfigs, importConfig.SessionConfigs...)
	}

	return config, "", nil
}

func (c *RealConfigurator) GetConfig() (model.Config, string, error) {
	config, details, err := c.getConfigFileFromUserConfigDir()
	if err != nil {
		return model.Config{}, details, err
	}
	return config, "", nil
}
