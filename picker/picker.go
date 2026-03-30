package picker

import (
	"errors"
	"fmt"

	tea "charm.land/bubbletea/v2"

	"github.com/joshmedeski/sesh/v2/model"
)

const (
	defaultPrompt      = "> "
	defaultPlaceholder = "Filter sessions..."
)

type PickerOptions struct {
	ShowIcons      *bool
	SeparatorAware *bool
	Prompt         *string
	Placeholder    *string
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

func (p *RealPicker) Pick(fetchFunc FetchFunc, opts PickerOptions) (string, error) {
	m := New(fetchFunc, *opts.ShowIcons, *opts.SeparatorAware, *opts.Prompt, *opts.Placeholder)
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
