package model

type (
	Config struct {
		ImportPaths          []string             `toml:"import"`
		DefaultSessionConfig DefaultSessionConfig `toml:"default_session"`
		SessionConfigs       []SessionConfig      `toml:"session"`
	}

	DefaultSessionConfig struct {
		// TODO: mention breaking change in v2 release notes
		// StartupScript  string `toml:"startup_script"`
		StartupCommand string `toml:"startup_command"`
		Tmuxp          string `toml:"tmuxp"`
		Tmuxinator     string `toml:"tmuxinator"`
		LsCommand      string `toml:"ls_command"`
	}

	SessionConfig struct {
		Name                string `toml:"name"`
		Path                string `toml:"path"`
		DisableStartCommand bool   `toml:"disable_startup_command"`
		DefaultSessionConfig
	}
)
