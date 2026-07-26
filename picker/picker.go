package picker

import (
	"errors"
	"fmt"
	"strings"
	"time"

	tea "charm.land/bubbletea/v2"

	"github.com/joshmedeski/sesh/v2/model"
)

const (
	defaultPrompt      = "> "
	defaultPlaceholder = "Filter sessions..."
)

type PickerOptions struct {
	ShowIcons               *bool
	ShowWindows             *bool
	SeparatorAware          *bool
	Prompt                  *string
	Placeholder             *string
	DisableAliasAutoConnect *bool
}

type Picker interface {
	Pick(fetchFunc FetchFunc, opts PickerOptions) (string, error)
}

type RealPicker struct {
	config model.Config
}

func NewPicker(config model.Config) Picker {
	return &RealPicker{config: config}
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

	m := New(fetchFunc, Options{
		ShowIcons:               showIcons,
		ShowWindows:             showWindows,
		SeparatorAware:          p.config.SeparatorAware,
		Prompt:                  prompt,
		Placeholder:             placeholder,
		Aliases:                 buildAliases(p.config.SessionConfigs),
		AliasFilterPrefix:       aliasFilterPrefix(p.config.TUI.AliasFilterPrefix),
		AliasAutoConnectDelay:   aliasAutoConnectDelay(p.config.TUI.AliasAutoConnectDelay),
		DisableAliasAutoConnect: disableAliasAutoConnect,
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
