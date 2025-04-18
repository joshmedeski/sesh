package connector

import (
	"fmt"

	"github.com/joshmedeski/sesh/v2/lister"
	"github.com/joshmedeski/sesh/v2/model"
)

func configStrategy(c *RealConnector, name string) (model.Connection, error) {
	config, exists := c.lister.FindConfigSession(name)
	if !exists {
		return model.Connection{Found: false}, nil
	}

	windows := make(model.SeshWindowMap)
	for _, window := range c.config.WindowConfigs {
		key := lister.ConfigKey(window.Name)
		var path string = ""
		var err error = nil
		if window.Path != "" {
			path, err = c.home.ExpandHome(window.Path)
			if err != nil {
				return model.Connection{}, fmt.Errorf("couldn't expand home: %q", err)
			}
		}

		if window.StartupScript != "" && window.DisableStartScript {
			return model.Connection{}, fmt.Errorf("startup_script and disable_start_script are mutually exclusive")
		}

		windows[key] = model.WindowConfig{
			Name:               window.Name,
			Path:               path,
			StartupScript:      window.StartupScript,
			DisableStartScript: window.DisableStartScript,
		}
	}

	windowConfigs := []model.WindowConfig{}
	for _, window := range config.WindowNames {
		windowConfig, ok := windows[lister.ConfigKey(window)]
		if !ok {
			return model.Connection{}, fmt.Errorf("window %s is not defined in config", window)
		}
		if windowConfig.Path == "" {
			path, err := c.home.ExpandHome(config.Path)
			if err != nil {
				return model.Connection{}, fmt.Errorf("couldn't expand home: %q", err)
			}
			windowConfig.Path = path
		}
		windowConfigs = append(windowConfigs, windowConfig)
	}

	config.WindowConfigs = windowConfigs

	return model.Connection{
		Found:       true,
		Session:     config,
		New:         true,
		AddToZoxide: true,
	}, nil
}
