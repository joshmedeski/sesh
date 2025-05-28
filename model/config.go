package model

type (
	Config struct {
		ImportPaths          []string             `toml:"import"`
		DefaultSessionConfig DefaultSessionConfig `toml:"default_session"`
		Blacklist            []string             `toml:"blacklist"`
		SessionConfigs       []SessionConfig      `toml:"session"`
    SortOrder            []string             `toml:"sort_order"`
		WindowConfigs        []WindowConfig       `toml:"window"`
	}

	DefaultSessionConfig struct {
		// TODO: mention breaking change in v2 release notes
		// StartupScript  string `toml:"startup_script"`
		StartupCommand string   `toml:"startup_command"`
		Tmuxp          string   `toml:"tmuxp"`
		Tmuxinator     string   `toml:"tmuxinator"`
		PreviewCommand string   `toml:"preview_command"`
		Windows        []string `toml:"windows"`
	}

	SessionConfig struct {
		Name                string `toml:"name"`
		Path                string `toml:"path"`
		DisableStartCommand bool   `toml:"disable_startup_command"`
		DefaultSessionConfig
	}

	WindowConfig struct {
		Name               string `toml:"name"`
		StartupScript      string `toml:"startup_script"`
		Path               string `toml:"path"`
	}
)
