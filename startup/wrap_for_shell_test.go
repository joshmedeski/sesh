package startup

import (
	"testing"

	"github.com/joshmedeski/sesh/v2/oswrap"
	"github.com/stretchr/testify/assert"
)

func TestWrapForShell(t *testing.T) {
	t.Run("returns empty for empty command", func(t *testing.T) {
		mockOs := new(oswrap.MockOs)
		s := &RealStartup{os: mockOs}
		assert.Equal(t, "", s.WrapForShell(""))
		mockOs.AssertNotCalled(t, "Getenv")
	})

	t.Run("wraps command with $SHELL from env (zsh)", func(t *testing.T) {
		mockOs := new(oswrap.MockOs)
		mockOs.On("Getenv", "SHELL").Return("/bin/zsh")
		s := &RealStartup{os: mockOs}

		got := s.WrapForShell("nvim")
		want := `'/bin/zsh' -i -c 'nvim; exec /bin/zsh -i -f'`
		assert.Equal(t, want, got)
	})

	t.Run("wraps command with $SHELL from env (bash)", func(t *testing.T) {
		mockOs := new(oswrap.MockOs)
		mockOs.On("Getenv", "SHELL").Return("/bin/bash")
		s := &RealStartup{os: mockOs}

		got := s.WrapForShell("nvim")
		want := `'/bin/bash' -i -c 'nvim; exec /bin/bash -i --norc --noprofile'`
		assert.Equal(t, want, got)
	})

	t.Run("falls back to /bin/sh when $SHELL unset", func(t *testing.T) {
		mockOs := new(oswrap.MockOs)
		mockOs.On("Getenv", "SHELL").Return("")
		s := &RealStartup{os: mockOs}

		got := s.WrapForShell("ls")
		want := `'/bin/sh' -i -c 'ls; exec /bin/sh -i'`
		assert.Equal(t, want, got)
	})

	t.Run("escapes single quotes in command", func(t *testing.T) {
		mockOs := new(oswrap.MockOs)
		mockOs.On("Getenv", "SHELL").Return("/bin/zsh")
		s := &RealStartup{os: mockOs}

		got := s.WrapForShell("echo 'hi'")
		want := `'/bin/zsh' -i -c 'echo '\''hi'\''; exec /bin/zsh -i -f'`
		assert.Equal(t, want, got)
	})
}
