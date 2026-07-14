package model

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
		WorktreeConfigs         []WorktreeConfig     `toml:"worktree"`
		DirLength               int                  `toml:"dir_length"`
		GitNamerUseWorktreeRoot bool                 `toml:"git_namer_use_worktree_root"`
		GitDirLength            int                  `toml:"git_dir_length"`
		SeparatorAware          bool                 `toml:"separator_aware"`
		TmuxCommand             string               `toml:"tmux_command"`
		Terminal                string               `toml:"terminal"`
		Browser                 BrowserConfig        `toml:"browser"`
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

	WildcardConfig struct {
		Pattern             string   `toml:"pattern"`
		StartupCommand      string   `toml:"startup_command"`
		DisableStartCommand bool     `toml:"disable_startup_command"`
		PreviewCommand      string   `toml:"preview_command"`
		Windows             []string `toml:"windows"`
	}

	// WorktreeConfig maps a GitHub "org/repo" to a local repository so
	// `sesh worktree create <number>` knows where to add worktrees, how to
	// name their branches, and what to run on connect.
	WorktreeConfig struct {
		Name           string `toml:"name"`            // GitHub "org/repo"
		Path           string `toml:"path"`            // local repo root (supports ~)
		WorktreeDir    string `toml:"worktree_dir"`    // default ".wk"; relative to Path or absolute
		BranchTemplate string `toml:"branch_template"` // default "{number}"
		BaseBranch     string `toml:"base_branch"`     // default "origin/main"
		Fetch          *bool  `toml:"fetch"`           // default true (nil => true)
		StartupCommand string `toml:"startup_command"` // runs on connect
	}

	// BrowserConfig configures reading the active browser tab's URL so
	// `sesh worktree create --browser` can derive the target issue/PR.
	// macOS-only (osascript). An empty Application disables the feature.
	BrowserConfig struct {
		Application string `toml:"application"` // browser app name, e.g. "Helium"
		URLCommand  string `toml:"url_command"` // AppleScript fragment; default "URL of active tab of front window"
	}
)
