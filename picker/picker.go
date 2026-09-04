package picker

import (
	"errors"
	"fmt"
	"strings"
	"time"

	tea "charm.land/bubbletea/v2"

	"github.com/joshmedeski/sesh/v2/home"
	"github.com/joshmedeski/sesh/v2/model"
	"github.com/joshmedeski/sesh/v2/previewer"
)

const (
	defaultPrompt              = "> "
	defaultPlaceholder         = "Filter sessions..."
	defaultWorktreePlaceholder = "Filter worktrees..."
)

type PickerOptions struct {
	ShowIcons               *bool
	ShowWindows             *bool
	SeparatorAware          *bool
	Prompt                  *string
	Placeholder             *string
	DisableAliasAutoConnect *bool
	Preview                 *bool
	Query                   *string
}

// WorktreePickerOptions overrides the configured defaults for a single
// invocation of the worktree picker. A nil field means "not overridden".
type WorktreePickerOptions struct {
	ShowIcons   *bool
	Prompt      *string
	Placeholder *string
	Query       *string
	AutoRefresh *bool
}

type Picker interface {
	Pick(fetchFunc FetchFunc, opts PickerOptions) (string, error)
	// PickWorktree picks a worktree by issue number or title. The second return
	// is false when nothing was chosen — the picker was quit, or had no rows.
	PickWorktree(fetchFunc WorktreeFetchFunc, opts WorktreePickerOptions) (model.WorktreeEntry, bool, error)
}

type RealPicker struct {
	config    model.Config
	previewer previewer.Previewer
	// home expands the paths on [[session]] blocks so a configured icon can be
	// matched against the absolute path a session is listed with.
	home home.Home
	// wildcards matches a session path to a [[wildcard]] block, for icons
	// declared on a pattern rather than a single session.
	wildcards WildcardFinder
}

func NewPicker(config model.Config, previewer previewer.Previewer, home home.Home, wildcards WildcardFinder) Picker {
	return &RealPicker{config: config, previewer: previewer, home: home, wildcards: wildcards}
}

// buildAliases collects the aliases defined on [[session]] blocks, keyed by
// lowercased alias so lookups are case-insensitive. Duplicate aliases are
// rejected at config load, so last-one-wins here is unreachable in practice.
func buildAliases(sessions []model.SessionConfig) map[string]Alias {
	aliases := make(map[string]Alias)
	for _, session := range sessions {
		if session.Alias == "" || session.Name == "" {
			continue
		}
		aliases[strings.ToLower(session.Alias)] = Alias{
			Alias:       session.Alias,
			Target:      session.Name,
			AutoConnect: session.AliasAutoConnect,
		}
	}
	return aliases
}

// aliasAutoConnectDelay parses the configured delay, falling back to the
// default rather than failing the picker on a malformed value (the configurator
// already rejects those at load time).
func aliasAutoConnectDelay(configured string) time.Duration {
	if d, err := time.ParseDuration(configured); err == nil {
		return d
	}
	d, _ := time.ParseDuration(model.DefaultAliasAutoConnectDelay)
	return d
}

// previewWidth resolves the percent of the terminal guaranteed to the preview
// pane, clamping rather than failing the picker on an out-of-range value (the
// JSON schema flags those in the editor, but nothing rejects them at load).
func previewWidth(configured int) int {
	if configured <= 0 {
		return model.DefaultPreviewWidth
	}
	if configured < model.MinPreviewWidth {
		return model.MinPreviewWidth
	}
	if configured > model.MaxPreviewWidth {
		return model.MaxPreviewWidth
	}
	return configured
}

// previewMinWidth resolves the narrowest terminal that still gets a preview
// pane. Zero means the key is absent.
func previewMinWidth(configured int) int {
	if configured <= 0 {
		return model.DefaultPreviewMinWidth
	}
	return configured
}

// previewBorder resolves the divider drawn between the list and the preview
// pane, falling back to the default on an empty or unrecognized value rather
// than failing the picker (the JSON schema flags those in the editor, but
// nothing rejects them at load).
func previewBorder(configured string) string {
	switch configured {
	case model.PreviewBorderNone, model.PreviewBorderLine, model.PreviewBorderThick, model.PreviewBorderDouble:
		return configured
	default:
		return model.DefaultPreviewBorder
	}
}

// aliasFilterPrefix resolves the sigil that enters alias-filter mode. A nil
// value means the key is absent (the configurator normally fills it in, but the
// picker can be constructed with a bare config), while an explicit empty string
// disables the mode.
func aliasFilterPrefix(configured *string) string {
	if configured == nil {
		return model.DefaultAliasFilterPrefix
	}
	return *configured
}

func (p *RealPicker) Pick(fetchFunc FetchFunc, opts PickerOptions) (string, error) {
	showIcons := false
	if opts.ShowIcons != nil {
		showIcons = *opts.ShowIcons
	} else {
		showIcons = p.config.TUI.ShowIcons
	}

	showWindows := p.config.TUI.ShowWindows
	if opts.ShowWindows != nil {
		showWindows = *opts.ShowWindows
	}

	prompt := defaultPrompt
	if opts.Prompt != nil {
		prompt = *opts.Prompt
	} else if p.config.TUI.Prompt != "" {
		prompt = p.config.TUI.Prompt
	}

	placeholder := defaultPlaceholder
	if opts.Placeholder != nil {
		placeholder = *opts.Placeholder
	} else if p.config.TUI.Placeholder != "" {
		placeholder = p.config.TUI.Placeholder
	}

	disableAliasAutoConnect := false
	if opts.DisableAliasAutoConnect != nil {
		disableAliasAutoConnect = *opts.DisableAliasAutoConnect
	}

	preview := p.config.TUI.Preview
	if opts.Preview != nil {
		preview = *opts.Preview
	}

	// A pre-filled query is per-invocation only: a configured default would
	// silently hide sessions on every launch.
	query := ""
	if opts.Query != nil {
		query = *opts.Query
	}

	// The model reaches the previewer through a function so the TUI stays
	// unaware of the package, mirroring how sessions arrive via FetchFunc.
	var previewFunc PreviewFunc
	if p.previewer != nil {
		previewFunc = p.previewer.Preview
	}

	m := New(fetchFunc, Options{
		ShowIcons:               showIcons,
		ShowWindows:             showWindows,
		SeparatorAware:          p.config.SeparatorAware,
		Prompt:                  prompt,
		Placeholder:             placeholder,
		Query:                   query,
		Aliases:                 buildAliases(p.config.SessionConfigs),
		AliasFilterPrefix:       aliasFilterPrefix(p.config.TUI.AliasFilterPrefix),
		AliasAutoConnectDelay:   aliasAutoConnectDelay(p.config.TUI.AliasAutoConnectDelay),
		DisableAliasAutoConnect: disableAliasAutoConnect,
		Icon:                    buildIconResolver(p.config, p.home, p.wildcards),
		IconWidth:               iconColWidth(p.config),
		Preview:                 preview,
		PreviewWidth:            previewWidth(p.config.TUI.PreviewWidth),
		PreviewMinWidth:         previewMinWidth(p.config.TUI.PreviewMinWidth),
		PreviewBorder:           p.config.TUI.PreviewBorder,
		PreviewFunc:             previewFunc,
		GroupSeparator:          p.config.TUI.GroupSeparator,
	})
	prog := tea.NewProgram(m)
	result, err := prog.Run()
	if err != nil {
		return "", fmt.Errorf("picker error: %w", err)
	}
	pickerModel, ok := result.(Model)
	if !ok {
		return "", errors.New("unexpected model type")
	}
	if pickerModel.LoadErr() != nil {
		return "", fmt.Errorf("couldn't list sessions: %w", pickerModel.LoadErr())
	}
	if pickerModel.Quit() {
		return "", nil
	}
	return pickerModel.Chosen(), nil
}

func (p *RealPicker) PickWorktree(fetchFunc WorktreeFetchFunc, opts WorktreePickerOptions) (model.WorktreeEntry, bool, error) {
	// The configured prompt and placeholder apply to both pickers: they say how
	// the user wants a filter input to look, which the kind of thing being
	// filtered doesn't change. Only the fallback placeholder differs, since the
	// session picker's would name the wrong thing.
	prompt := defaultPrompt
	if opts.Prompt != nil {
		prompt = *opts.Prompt
	} else if p.config.TUI.Prompt != "" {
		prompt = p.config.TUI.Prompt
	}

	placeholder := defaultWorktreePlaceholder
	if opts.Placeholder != nil {
		placeholder = *opts.Placeholder
	} else if p.config.TUI.Placeholder != "" {
		placeholder = p.config.TUI.Placeholder
	}

	// A pre-filled query is per-invocation only: a configured default would
	// silently hide worktrees on every launch.
	query := ""
	if opts.Query != nil {
		query = *opts.Query
	}

	// show_icons is read as the statement that the font has nerd font glyphs,
	// which is what decides whether the badges get their half circles.
	showIcons := p.config.TUI.ShowIcons
	if opts.ShowIcons != nil {
		showIcons = *opts.ShowIcons
	}

	// Refreshing behind the cached rows is the default: the alternative is a
	// picker that shows a title edited yesterday as it was yesterday, and only
	// tells someone who already knew to press ctrl+r.
	autoRefresh := true
	if opts.AutoRefresh != nil {
		autoRefresh = *opts.AutoRefresh
	}

	m := NewWorktreeModel(fetchFunc, WorktreeOptions{
		Prompt:      prompt,
		Placeholder: placeholder,
		ShowIcons:   showIcons,
		Query:       query,
		AutoRefresh: autoRefresh,
	})
	prog := tea.NewProgram(m)
	result, err := prog.Run()
	if err != nil {
		return model.WorktreeEntry{}, false, fmt.Errorf("picker error: %w", err)
	}
	pickerModel, ok := result.(WorktreeModel)
	if !ok {
		return model.WorktreeEntry{}, false, errors.New("unexpected model type")
	}
	if pickerModel.LoadErr() != nil {
		return model.WorktreeEntry{}, false, fmt.Errorf("couldn't list worktrees: %w", pickerModel.LoadErr())
	}
	entry, picked := pickerModel.Chosen()
	return entry, picked, nil
}
