package seshcli

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/joshmedeski/sesh/v2/marker"
	cli "github.com/urfave/cli/v2"
)

func Mark(marker marker.Marker) *cli.Command {
	return &cli.Command{
		Name:    "mark",
		Aliases: []string{"m"},
		Usage:   "Mark current session:window for priority listing",
		Action: func(cCtx *cli.Context) error {
			session, window, err := getCurrentSessionWindow()
			if err != nil {
				return fmt.Errorf("failed to get current session:window: %w", err)
			}

			if err := marker.Mark(session, window); err != nil {
				return fmt.Errorf("failed to mark session %s:%s: %w", session, window, err)
			}

			fmt.Printf("Marked %s:%s\n", session, window)
			return nil
		},
	}
}

func Unmark(marker marker.Marker) *cli.Command {
	return &cli.Command{
		Name:    "unmark",
		Aliases: []string{"u"},
		Usage:   "Unmark current session:window",
		Action: func(cCtx *cli.Context) error {
			session, window, err := getCurrentSessionWindow()
			if err != nil {
				return fmt.Errorf("failed to get current session:window: %w", err)
			}

			if err := marker.Unmark(session, window); err != nil {
				return fmt.Errorf("failed to unmark session %s:%s: %w", session, window, err)
			}

			fmt.Printf("Unmarked %s:%s\n", session, window)
			return nil
		},
	}
}

func getCurrentSessionWindow() (string, string, error) {
	sessionCmd := exec.Command("tmux", "display-message", "-p", "#S")
	sessionOut, err := sessionCmd.Output()
	if err != nil {
		return "", "", fmt.Errorf("failed to get current session: %w", err)
	}
	session := strings.TrimSpace(string(sessionOut))

	windowCmd := exec.Command("tmux", "display-message", "-p", "#I")
	windowOut, err := windowCmd.Output()
	if err != nil {
		return "", "", fmt.Errorf("failed to get current window: %w", err)
	}
	window := strings.TrimSpace(string(windowOut))

	return session, window, nil
}