package model

type (
	Config struct {
		ImportPaths          []string             `toml:"import"`
		DefaultSessionConfig DefaultSessionConfig `toml:"default_session"`
		SessionConfigs       []SessionConfig      `toml:"session"`
	}

	DefaultSessionConfig struct {
		StartupScript  string `toml:"startup_script"`
		StartupCommand string `toml:"startup_command"`
		Tmuxp          string `toml:"tmuxp"`
		Tmuxinator     string `toml:"tmuxinator"`
	}

	SessionConfig struct {
		Name string `toml:"name"`
		Path string `toml:"path"`
		DefaultSessionConfig
	}
)
