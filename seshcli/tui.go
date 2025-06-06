package seshcli

import (
	"log/slog"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/joshmedeski/sesh/v2/tui"
	cli "github.com/urfave/cli/v2"
)

func Tui(t tui.Tui) *cli.Command {
	return &cli.Command{
		Name:                   "ui",
		Aliases:                []string{"u"},
		Usage:                  "View the user interface",
		UseShortOptionHandling: true,
		Flags:                  []cli.Flag{},
		Action: func(cCtx *cli.Context) error {
			model := t.NewModel()
			p := tea.NewProgram(model)
			if _, err := p.Run(); err != nil {
				slog.Error("Whoops, something went wrong", err.Error())
				os.Exit(1)
			}
			return nil
		},
	}
}
