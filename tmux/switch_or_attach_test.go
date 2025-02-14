package tmux

import (
	"errors"
	"testing"

	"github.com/joshmedeski/sesh/v2/model"
	"github.com/joshmedeski/sesh/v2/oswrap"
	"github.com/joshmedeski/sesh/v2/shell"
	"github.com/stretchr/testify/assert"
	mock "github.com/stretchr/testify/mock"
)

func TestSwitchOrAttach(t *testing.T) {
	mockOs := new(oswrap.MockOs)
	mockShell := new(shell.MockShell)
	tmux := NewTmux(mockOs, mockShell)

	t.Run("switches because of option", func(t *testing.T) {
		mockOs.ExpectedCalls = nil
		mockShell.ExpectedCalls = nil
		mockShell.On("Cmd", "tmux", "switch-client", "-t", mock.Anything).Return("", nil)
		response, error := tmux.SwitchOrAttach("dotfiles", model.ConnectOpts{Switch: true})
		assert.Equal(t, "switching to tmux session: dotfiles", response)
		assert.Equal(t, nil, error)
	})

	t.Run("switches when attached", func(t *testing.T) {
		mockOs.ExpectedCalls = nil
		mockShell.ExpectedCalls = nil
		mockOs.On("Getenv", "TMUX").Return("/private/tmp/tmux-501/default,72439,4")
		mockShell.On("Cmd", "tmux", "switch-client", "-t", mock.Anything).Return("", nil)
		response, error := tmux.SwitchOrAttach("dotfiles", model.ConnectOpts{Switch: false})
		assert.Equal(t, "switching to tmux session: dotfiles", response)
		assert.Equal(t, nil, error)
	})

	t.Run("errors when switching to a missing session", func(t *testing.T) {
		mockOs.ExpectedCalls = nil
		mockShell.ExpectedCalls = nil
		mockOs.On("Getenv", "TMUX").Return("/private/tmp/tmux-501/default,72439,4")
		mockShell.On("Cmd", "tmux", "switch-client", "-t", mock.Anything).Return("", errors.New("can't find session: dotfiles"))
		response, err := tmux.SwitchOrAttach("dotfiles", model.ConnectOpts{Switch: false})
		assert.Equal(t, "", response)
		assert.EqualError(t, err, "failed to switch to tmux session: can't find session: dotfiles")
	})

	t.Run("attaches", func(t *testing.T) {
		mockOs.ExpectedCalls = nil
		mockShell.ExpectedCalls = nil
		mockOs.On("Getenv", "TMUX").Return("")
		mockShell.On("Cmd", "tmux", "attach-session", "-t", mock.Anything).Return("", nil)
		response, error := tmux.SwitchOrAttach("dotfiles", model.ConnectOpts{Switch: false})
		assert.Equal(t, "attaching to tmux session: dotfiles", response)
		assert.Equal(t, nil, error)
	})
}
