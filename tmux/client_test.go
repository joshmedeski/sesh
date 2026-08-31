package tmux

import (
	"errors"
	"testing"

	"github.com/joshmedeski/sesh/v2/model"
	"github.com/joshmedeski/sesh/v2/oswrap"
	"github.com/joshmedeski/sesh/v2/shell"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCurrentSessionID(t *testing.T) {
	tests := []struct {
		name     string
		tmuxEnv  string
		expected string
	}{
		{"popup and pane alike carry the session id", "/private/tmp/tmux-501/default,3746,31", "$31"},
		{"outside tmux", "", ""},
		{"truncated by a wrapper", "/private/tmp/tmux-501/default,3746", ""},
		{"empty session field", "/private/tmp/tmux-501/default,3746,", ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, currentSessionID(tt.tmuxEnv))
		})
	}
}

// clientLines renders what `tmux list-clients -F clientFormat` would print.
func clientLines(lines ...string) []string {
	return append(lines, "")
}

func TestResolveClient(t *testing.T) {
	const (
		onSession1   = "/dev/ttys001\t/dev/ttys001\t$1\t100"
		onSession2   = "/dev/ttys002\t/dev/ttys002\t$2\t200"
		alsoSession1 = "/dev/ttys003\t/dev/ttys003\t$1\t300"
	)

	tests := []struct {
		name       string
		seshClient string
		tmuxEnv    string
		clients    []string
		clientsErr error
		expected   string
	}{
		{
			name:       "SESH_CLIENT wins outright",
			seshClient: "/dev/ttys009",
			tmuxEnv:    "/tmp/default,3746,1",
			clients:    clientLines(onSession1),
			expected:   "/dev/ttys009",
		},
		{
			// The popup case: tmux would guess ttys002 as most recently
			// active server-wide, but the popup was launched from $1.
			name:     "prefers the client attached to the calling session",
			tmuxEnv:  "/tmp/default,3746,1",
			clients:  clientLines(onSession1, onSession2),
			expected: "/dev/ttys001",
		},
		{
			name:     "picks the most recently active of several on that session",
			tmuxEnv:  "/tmp/default,3746,1",
			clients:  clientLines(onSession1, alsoSession1),
			expected: "/dev/ttys003",
		},
		{
			name:     "falls back to most recently active when the session has no client",
			tmuxEnv:  "/tmp/default,3746,7",
			clients:  clientLines(onSession1, onSession2),
			expected: "/dev/ttys002",
		},
		{
			name:     "falls back to most recently active outside tmux",
			tmuxEnv:  "",
			clients:  clientLines(onSession1, onSession2),
			expected: "/dev/ttys002",
		},
		{
			name:     "no clients attached",
			tmuxEnv:  "/tmp/default,3746,1",
			clients:  clientLines(),
			expected: "",
		},
		{
			name:       "list-clients failed",
			tmuxEnv:    "/tmp/default,3746,1",
			clientsErr: errors.New("no server running"),
			expected:   "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockOs := oswrap.NewMockOs(t)
			mockShell := shell.NewMockShell(t)
			mockOs.EXPECT().Getenv("SESH_CLIENT").Return(tt.seshClient)
			if tt.seshClient == "" {
				mockShell.EXPECT().ListCmd("tmux", "list-clients", "-F", clientFormat).
					Return(tt.clients, tt.clientsErr)
				mockOs.EXPECT().Getenv("TMUX").Return(tt.tmuxEnv).Maybe()
			}
			tm := NewTmux(mockOs, mockShell, "")
			assert.Equal(t, tt.expected, tm.ResolveClient())
		})
	}
}

func TestSwitchOrAttachNamesTheResolvedClient(t *testing.T) {
	mockOs := oswrap.NewMockOs(t)
	mockShell := shell.NewMockShell(t)
	mockOs.EXPECT().Getenv("SESH_CLIENT").Return("")
	mockOs.EXPECT().Getenv("TMUX").Return("/tmp/default,3746,1")
	mockShell.EXPECT().ListCmd("tmux", "list-clients", "-F", clientFormat).
		Return([]string{
			"/dev/ttys001\t/dev/ttys001\t$1\t100",
			"/dev/ttys002\t/dev/ttys002\t$2\t200",
		}, nil)
	mockShell.EXPECT().Cmd("tmux", "switch-client", "-c", "/dev/ttys001", "-t", "dotfiles").
		Return("", nil)

	tm := NewTmux(mockOs, mockShell, "")
	response, err := tm.SwitchOrAttach("dotfiles", model.ConnectOpts{Switch: true})
	require.NoError(t, err)
	assert.Equal(t, "switching to tmux session: dotfiles", response)
}
