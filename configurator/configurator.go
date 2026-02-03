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
	GetConfig() (model.Config, error) // Since error is an interface, we use it here to return a single variable instead of multiple variables (configError holds 2 strings, human and err)
}

type RealConfigurator struct {
	os         oswrap.Os
	path       pathwrap.Path
	runtime    runtimewrap.Runtime
	configPath string
}

// Helper for consolidation of error into a single structure
type ConfigError struct {
	Err          string // Load the (DecodeError/StrictMissingError).Error() into this
	HumanDetails string // Load the (DecodeError/StrictMissingError).String() into this
}

func (ce *ConfigError) Error() string {
	return ce.Err // Return the error
}

func (ce *ConfigError) Human() string {
	return ce.HumanDetails // Return the string
}

func NewConfigurator(os oswrap.Os, path pathwrap.Path, runtime runtimewrap.Runtime) Configurator {
	return &RealConfigurator{os: os, path: path, runtime: runtime}
}

func NewConfiguratorWithPath(os oswrap.Os, path pathwrap.Path, runtime runtimewrap.Runtime, configPath string) Configurator {
	return &RealConfigurator{os: os, path: path, runtime: runtime, configPath: configPath}
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

func (c *RealConfigurator) parseConfigFile(file []byte) (model.Config, error) {
	config := model.Config{}
	evaluation := model.Evaluation{}

	_ = toml.Unmarshal(file, &evaluation)
	if evaluation.StrictMode {
		reader := strings.NewReader(string(file))
		d := toml.NewDecoder(reader)
		d.DisallowUnknownFields() // enable the strict mode
		err := d.Decode(&config)
		if err != nil {
			var details *toml.StrictMissingError
			if errors.As(err, &details) {
				return config, &ConfigError{Err: err.Error(), HumanDetails: details.String()}
			}
			return config, err
		}
	} else {
		err := toml.Unmarshal(file, &config)
		if err != nil {
			var derr *toml.DecodeError
			if errors.As(err, &derr) {
				return config, &ConfigError{Err: err.Error(), HumanDetails: derr.String()}
			}
		}
	}

	return config, nil
}

func (c *RealConfigurator) resolveImports(config *model.Config, homeDir string) error {
	for _, importPath := range config.ImportPaths {
		importFilePath, err := c.fullImportPath(homeDir, importPath)
		if err != nil {
			return fmt.Errorf("couldn't get full import path: %q", err)
		}

		importFile, err := c.os.ReadFile(importFilePath)
		if err != nil {
			return fmt.Errorf("couldn't read import file %s: %q", importFilePath, err)
		}

		importConfig := model.Config{}
		if err := toml.Unmarshal(importFile, &importConfig); err != nil {
			return fmt.Errorf("couldn't unmarshal import file %s: %q", importFilePath, err)
		}

		config.SessionConfigs = append(config.SessionConfigs, importConfig.SessionConfigs...)
		config.WindowConfigs = append(config.WindowConfigs, importConfig.WindowConfigs...)
		config.WildcardConfigs = append(config.WildcardConfigs, importConfig.WildcardConfigs...)
	}
	return nil
}

func (c *RealConfigurator) applyDefaults(config *model.Config) {
	if config.DirLength < 1 {
		config.DirLength = 1
	}
}

func (c *RealConfigurator) getConfigFileFromPath(configPath string) (model.Config, error) {
	file, err := c.os.ReadFile(configPath)
	if err != nil {
		return model.Config{}, fmt.Errorf("couldn't read config file %q: %w", configPath, err)
	}

	config, err := c.parseConfigFile(file)
	if err != nil {
		return model.Config{}, err
	}

	userHomeDir, err := c.os.UserHomeDir()
	if err != nil {
		return config, fmt.Errorf("couldn't get user home dir: %q", err)
	}

	if err := c.resolveImports(&config, userHomeDir); err != nil {
		return config, err
	}

	c.applyDefaults(&config)
	return config, nil
}

func (c *RealConfigurator) getConfigFileFromUserConfigDir() (model.Config, error) {
	userHomeDir, err := c.os.UserHomeDir()
	if err != nil {
		return model.Config{}, fmt.Errorf("couldn't get user config dir: %q", err)
	}

	// Check XDG_CONFIG_HOME first, fall back to $HOME/.config
	userConfigDir := c.os.Getenv("XDG_CONFIG_HOME")
	if userConfigDir == "" {
		userConfigDir = c.path.Join(userHomeDir, ".config")
	}

	configFilePath := c.configFilePath(userConfigDir)
	file, _ := c.os.ReadFile(configFilePath)
	// TODO: add to debugging logs (Update, added details string)
	// if err != nil {
	// 	return config, "", fmt.Errorf("couldn't read config file: %q", err)
	// }

	config, err := c.parseConfigFile(file)
	if err != nil {
		return config, err
	}

	if err := c.resolveImports(&config, userHomeDir); err != nil {
		return config, err
	}

	c.applyDefaults(&config)
	return config, nil
}

func (c *RealConfigurator) GetConfig() (model.Config, error) {
	if c.configPath != "" {
		return c.getConfigFileFromPath(c.configPath)
	}
	config, err := c.getConfigFileFromUserConfigDir()
	if err != nil {
		return model.Config{}, err
	}
	return config, nil
}
