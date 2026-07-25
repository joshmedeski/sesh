package model

// DefaultAliasAutoConnectDelay is the grace period used when
// [tui] alias_auto_connect_delay is not set.
const DefaultAliasAutoConnectDelay = "150ms"

type (
	Config struct {
		Cache                   bool                 `toml:"cache"`
		StrictMode              bool                 `toml:"strict_mode"`
		ImportPaths             []string             `toml:"import"`
		DefaultSessionConfig    DefaultSessionConfig `toml:"default_session"`
		Blacklist               []string             `toml:"blacklist"`
		SessionConfigs          []SessionConfig      `toml:"session"`
		SortOrder               []string             `toml:"sort_order"`
		WindowConfigs           []WindowConfig       `toml:"window"`
		WildcardConfigs         []WildcardConfig     `toml:"wildcard"`
		DirLength               int                  `toml:"dir_length"`
		GitNamerUseWorktreeRoot bool                 `toml:"git_namer_use_worktree_root"`
		GitDirLength            int                  `toml:"git_dir_length"`
		SeparatorAware          bool                 `toml:"separator_aware"`
		TmuxCommand             string               `toml:"tmux_command"`
		Frecency                FrecencyConfig       `toml:"frecency"`
		TUI                     TUIConfig            `toml:"tui"`
	}
	Evaluation struct {
		StrictMode bool `toml:"strict_mode"`
	}

	// FrecencyConfig overrides the commands used to drive the frecency
	// directory-jumping backend (zoxide by default). Each command may be
	// swapped for an alternative tool such as fasd, autojump, or memy.
	// Empty fields fall back to the zoxide defaults, so an absent [frecency]
	// table leaves behavior byte-identical to prior versions.
	FrecencyConfig struct {
		// ListCommand enumerates all tracked entries. Its output is parsed
		// one path per line; a leading numeric score is detected and used
		// when present (e.g. zoxide's `query --list --score`).
		ListCommand string `toml:"list_command"`
		// QueryCommand resolves a single input to a path. The `{}`
		// placeholder is replaced with the query string.
		QueryCommand string `toml:"query_command"`
		// AddCommand records a path to bump its frecency after connecting.
		// The `{}` placeholder is replaced with the path.
		AddCommand string `toml:"add_command"`
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
		// Alias is a short, exact-match shortcut for this session. Typing it in
		// the picker marks the session with a chip, and `sesh connect <alias>`
		// resolves to this session's name.
		Alias string `toml:"alias"`
		// AliasAutoConnect connects to this session as soon as its alias is
		// fully typed in the picker, without pressing enter. Opt-in per session
		// since it is only desirable for sessions you jump to constantly.
		AliasAutoConnect bool `toml:"alias_auto_connect"`
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
		ShowWindows bool   `toml:"show_windows"`
		Prompt      string `toml:"prompt"`
		Placeholder string `toml:"placeholder"`
		// AliasAutoConnectDelay is the grace period between typing an alias and
		// auto-connecting to it, giving longer aliases that share a prefix time
		// to be typed. Any duration string time.ParseDuration accepts.
		AliasAutoConnectDelay string `toml:"alias_auto_connect_delay"`
	}

	WildcardConfig struct {
		Pattern             string   `toml:"pattern"`
		StartupCommand      string   `toml:"startup_command"`
		DisableStartCommand bool     `toml:"disable_startup_command"`
		PreviewCommand      string   `toml:"preview_command"`
		Windows             []string `toml:"windows"`
	}
)
