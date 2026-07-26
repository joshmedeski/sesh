package model

// DefaultAliasAutoConnectDelay is the grace period used when
// [tui] alias_auto_connect_delay is not set.
const DefaultAliasAutoConnectDelay = "150ms"

// DefaultAliasFilterPrefix is the sigil that enters alias-filter mode when
// [tui] alias_filter_prefix is not set.
const DefaultAliasFilterPrefix = "/"

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
	}
)
