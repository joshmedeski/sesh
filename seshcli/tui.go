package seshcli

import (
	"log/slog"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/joshmedeski/sesh/v2/tui"
	"github.com/spf13/cobra"
)

func NewTuiCommand(t tui.Tui) *cobra.Command {
	return &cobra.Command{
		Use:     "ui",
		Aliases: []string{"u"},
		Short:   "View the user interface",
		RunE: func(cmd *cobra.Command, args []string) error {
			model := t.NewModel()
			p := tea.NewProgram(model, tea.WithAltScreen())
			if _, err := p.Run(); err != nil {
				slog.Error("Whoops, something went wrong", "err", err.Error())
				os.Exit(1)
			}
			return nil
		},
	}
}
