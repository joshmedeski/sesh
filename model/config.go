package model

type (
	Config struct {
		Cache                bool                 `toml:"cache"`
		StrictMode           bool                 `toml:"strict_mode"`
		ImportPaths          []string             `toml:"import"`
		DefaultSessionConfig DefaultSessionConfig `toml:"default_session"`
		Blacklist            []string             `toml:"blacklist"`
		SessionConfigs       []SessionConfig      `toml:"session"`
		SortOrder            []string             `toml:"sort_order"`
		WindowConfigs        []WindowConfig       `toml:"window"`
		WildcardConfigs      []WildcardConfig     `toml:"wildcard"`
		DirLength            int                  `toml:"dir_length"`
		SeparatorAware       bool                 `toml:"separator_aware"`
		TmuxCommand          string               `toml:"tmux_command"`
		TUI                  TUIConfig            `toml:"tui"`
		Github               GithubConfig         `toml:"github"`
	}
	Evaluation struct {
		StrictMode bool `toml:"strict_mode"`
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
		Name          string `toml:"name"`
		StartupScript string `toml:"startup_script"`
		Path          string `toml:"path"`
	}

	TUIConfig struct {
		// TODO: keybindings and more
		ShowIcons   bool   `toml:"show_icons"`
		Prompt      string `toml:"prompt"`
		Placeholder string `toml:"placeholder"`
	}

	// GithubConfig holds settings for GitHub integration in the status bar.
	GithubConfig struct {
		// IssueTTL is the status cache lifetime in seconds. A pointer so an absent
		// section (nil → default 60) is distinguishable from an explicit 0 (disable
		// caching, always fetch live).
		IssueTTL *int `toml:"issue_ttl"`
	}

	WildcardConfig struct {
		Pattern             string   `toml:"pattern"`
		StartupCommand      string   `toml:"startup_command"`
		DisableStartCommand bool     `toml:"disable_startup_command"`
		PreviewCommand      string   `toml:"preview_command"`
		Windows             []string `toml:"windows"`
	}
)

// EffectiveTTL returns the cache TTL in seconds: 60 when unset, otherwise the
// configured value (0 means caching is disabled).
func (g GithubConfig) EffectiveTTL() int {
	if g.IssueTTL == nil {
		return 60
	}
	return *g.IssueTTL
}
