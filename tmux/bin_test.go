package tmux

import (
	"testing"

	"github.com/joshmedeski/sesh/v2/oswrap"
	"github.com/joshmedeski/sesh/v2/shell"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestConfiguredBinaryIsUsed(t *testing.T) {
	tests := []struct {
		name   string
		expect func(s *shell.MockShell)
		call   func(tm Tmux) error
	}{
		{
			name: "ListWindows",
			expect: func(s *shell.MockShell) {
				s.EXPECT().ListCmd("psmux", "list-windows", "-F", mock.Anything).Return(nil, nil)
			},
			call: func(tm Tmux) error {
				_, err := tm.ListWindows("")
				return err
			},
		},
		{
			name: "ListTmuxPanes",
			expect: func(s *shell.MockShell) {
				s.EXPECT().ListCmd("psmux", "list-panes", "-s", "-F", mock.Anything).Return(nil, nil)
			},
			call: func(tm Tmux) error {
				_, err := tm.ListTmuxPanes()
				return err
			},
		},
		{
			name: "SelectWindow",
			expect: func(s *shell.MockShell) {
				s.EXPECT().Cmd("psmux", "select-window", "-t", "editor").Return("", nil)
			},
			call: func(tm Tmux) error {
				_, err := tm.SelectWindow("editor")
				return err
			},
		},
		{
			name: "GetCurrentSession",
			expect: func(s *shell.MockShell) {
				s.EXPECT().Cmd("psmux", "display-message", "-p", "#{session_name}").Return("work", nil)
			},
			call: func(tm Tmux) error {
				_, err := tm.GetCurrentSession()
				return err
			},
		},
		{
			name: "SelectPane",
			expect: func(s *shell.MockShell) {
				s.EXPECT().Cmd("psmux", "display-message", "-p", "#{session_name}").Return("work", nil)
				s.EXPECT().Cmd("psmux", "select-window", "-t", "work:1").Return("", nil)
				s.EXPECT().Cmd("psmux", "select-pane", "-t", "work:1.2").Return("", nil)
			},
			call: func(tm Tmux) error {
				_, err := tm.SelectPane(1, 2)
				return err
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := shell.NewMockShell(t)
			tt.expect(s)
			tm := NewTmux(oswrap.NewMockOs(t), s, "psmux")
			require.NoError(t, tt.call(tm))
		})
	}
}
