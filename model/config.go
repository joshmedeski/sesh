package model

// DefaultAliasAutoConnectDelay is the grace period used when
// [tui] alias_auto_connect_delay is not set.
const DefaultAliasAutoConnectDelay = "150ms"

// DefaultAliasFilterPrefix is the sigil that enters alias-filter mode when
// [tui] alias_filter_prefix is not set.
const DefaultAliasFilterPrefix = "/"

// DefaultWindowNameFormat is the tmux format used for window names when
// [tui] window_name_format is not set.
const DefaultWindowNameFormat = "#{window_name}"

const (
	// DefaultPreviewWidth is the share of the terminal, in percent, guaranteed
	// to the picker's preview pane when [tui] preview_width is not set.
	DefaultPreviewWidth = 60
	// MinPreviewWidth and MaxPreviewWidth bound preview_width. A pane narrower
	// or wider than these leaves too little room for the other half.
	MinPreviewWidth = 10
	MaxPreviewWidth = 90
	// DefaultPreviewMinWidth is the narrowest terminal, in columns, that still
	// gets a preview pane when [tui] preview_min_width is not set.
	DefaultPreviewMinWidth = 100
)

// The dividers [tui] preview_border accepts between the picker's session list
// and its preview pane. PreviewBorderNone draws no divider at all.
const (
	PreviewBorderNone   = "none"
	PreviewBorderLine   = "line"
	PreviewBorderThick  = "thick"
	PreviewBorderDouble = "double"
	// DefaultPreviewBorder is the divider used when [tui] preview_border is
	// not set.
	DefaultPreviewBorder = PreviewBorderLine
)

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
		// Alias is a short, exact-match shortcut for this session. Typing it in
		// the picker marks the session with a chip, and `sesh connect <alias>`
		// resolves to this session's name.
		Alias string `toml:"alias"`
		// AliasAutoConnect connects to this session as soon as its alias is
		// fully typed in the picker, without pressing enter. Opt-in per session
		// since it is only desirable for sessions you jump to constantly.
		AliasAutoConnect bool `toml:"alias_auto_connect"`
		// Icon replaces the source glyph this session gets in the picker with
		// any string — a nerd font glyph or an emoji. Picker-only: `sesh list`
		// consumers trim a known-width glyph, so a custom one would break them.
		Icon string `toml:"icon"`
		DefaultSessionConfig
	}

	WindowConfig struct {
		Name          string `toml:"name"`
		StartupScript string `toml:"startup_script"`
		Path          string `toml:"path"`
	}

	TUIConfig struct {
		// TODO: keybindings and more
		ShowIcons   bool `toml:"show_icons"`
		ShowWindows bool `toml:"show_windows"`
		// WindowNameFormat is the tmux format rendered for each window when
		// ShowWindows is enabled. Empty means DefaultWindowNameFormat.
		WindowNameFormat string `toml:"window_name_format"`
		Prompt           string `toml:"prompt"`
		Placeholder      string `toml:"placeholder"`
		// AliasAutoConnectDelay is the grace period between typing an alias and
		// auto-connecting to it, giving longer aliases that share a prefix time
		// to be typed. Any duration string time.ParseDuration accepts.
		AliasAutoConnectDelay string `toml:"alias_auto_connect_delay"`
		// AliasFilterPrefix is the single character that, typed first in the
		// picker, narrows the list to aliased sessions. It is a pointer so an
		// explicit empty string (disable the mode) is distinguishable from an
		// absent key (use DefaultAliasFilterPrefix).
		AliasFilterPrefix *string `toml:"alias_filter_prefix"`
		// Preview shows a preview of the highlighted session beside the list,
		// using the same output as `sesh preview`. Opt-in: the pane costs a
		// command per cursor move, and not everyone wants a split.
		Preview bool `toml:"preview"`
		// PreviewWidth is the share of the terminal, in percent, guaranteed to
		// the preview pane. The list keeps its column cap, so anything left
		// over goes to the preview on top of this. Zero means unset.
		PreviewWidth int `toml:"preview_width"`
		// PreviewMinWidth is the narrowest terminal that still gets a preview
		// pane; below it the picker renders list-only. Zero means unset.
		PreviewMinWidth int `toml:"preview_min_width"`
		// PreviewBorder is the divider drawn between the list and the preview
		// pane: "line" (default), "thick", "double", or "none" for no divider.
		// Empty means unset.
		PreviewBorder string `toml:"preview_border"`
	}

	WildcardConfig struct {
		Pattern             string   `toml:"pattern"`
		StartupCommand      string   `toml:"startup_command"`
		DisableStartCommand bool     `toml:"disable_startup_command"`
		PreviewCommand      string   `toml:"preview_command"`
		Windows             []string `toml:"windows"`
		// Icon replaces the source glyph every session under this pattern gets
		// in the picker. See SessionConfig.Icon.
		Icon string `toml:"icon"`
	}

	// WorktreeConfig maps a GitHub "org/repo" to a local repository so
	// `sesh worktree connect <number>` knows where to add worktrees, how to
	// name their branches, and what to run on connect.
	WorktreeConfig struct {
		Repo           string `toml:"repo"`            // GitHub "org/repo"
		Path           string `toml:"path"`            // local repo root (supports ~)
		WorktreeDir    string `toml:"worktree_dir"`    // default ".wk"; relative to Path or absolute
		BranchTemplate string `toml:"branch_template"` // default "{number}"
		BaseBranch     string `toml:"base_branch"`     // default "origin/main"
		Fetch          *bool  `toml:"fetch"`           // default true (nil => true)
		StartupCommand string `toml:"startup_command"` // runs when connecting to a worktree that already existed
		CreateCommand  string `toml:"create_command"`  // runs instead, on the connect that creates the worktree
	}

	// BrowserConfig configures reading the active browser tab's URL so
	// `sesh worktree connect --browser` can derive the target issue/PR.
	// macOS-only (osascript). An empty Application disables the feature.
	BrowserConfig struct {
		Application string `toml:"application"` // browser app name, e.g. "Helium"
		URLCommand  string `toml:"url_command"` // AppleScript fragment; default "URL of active tab of front window"
	}
)
