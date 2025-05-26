package model

type (
	Config struct {
		ImportPaths          []string             `toml:"import"`
		DefaultSessionConfig DefaultSessionConfig `toml:"default_session"`
		MarkerConfig         MarkerConfig         `toml:"marker"`
		Blacklist            []string             `toml:"blacklist"`
		SessionConfigs       []SessionConfig      `toml:"session"`
	}

	DefaultSessionConfig struct {
		// TODO: mention breaking change in v2 release notes
		// StartupScript  string `toml:"startup_script"`
		StartupCommand string `toml:"startup_command"`
		Tmuxp          string `toml:"tmuxp"`
		Tmuxinator     string `toml:"tmuxinator"`
		PreviewCommand string `toml:"preview_command"`
	}

	SessionConfig struct {
		Name                string `toml:"name"`
		Path                string `toml:"path"`
		DisableStartCommand bool   `toml:"disable_startup_command"`
		DefaultSessionConfig
	}

	MarkerConfig struct {
		InactivityThreshold int `toml:"inactivity_threshold"` // Seconds before alert starts
		AlertLevel1Time     int `toml:"alert_level_1_time"`   // Seconds for level 1 alert (light)
		AlertLevel2Time     int `toml:"alert_level_2_time"`   // Seconds for level 2 alert (medium)
		AlertLevel3Time     int `toml:"alert_level_3_time"`   // Seconds for level 3 alert (urgent)
	}
)
