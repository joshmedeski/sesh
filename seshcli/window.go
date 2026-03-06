package seshcli

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

func NewWindowCommand(base *BaseDeps) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "window",
		Aliases: []string{"w"},
		Short:   "List or switch/create windows in a tmux session",
		RunE: func(cmd *cobra.Command, args []string) error {
			targetSession, _ := cmd.Flags().GetString("session")
			jsonOutput, _ := cmd.Flags().GetBool("json")
			_ = jsonOutput

			if targetSession == "" {
				if !base.Tmux.IsAttached() {
					return fmt.Errorf("not inside a tmux session, use --session to specify one")
				}
			} else {
				sessions, err := base.Tmux.ListSessions()
				if err != nil {
					return err
				}
				found := false
				for _, s := range sessions {
					if s.Name == targetSession {
						found = true
						break
					}
				}
				if !found {
					return fmt.Errorf("session '%s' not found", targetSession)
				}
			}

			if len(args) == 0 {
				windows, err := base.Tmux.ListWindows(targetSession)
				if err != nil {
					return err
				}
				for _, w := range windows {
					fmt.Println(w.Name)
				}
				return nil
			}

			name := strings.Join(args, " ")

			windows, err := base.Tmux.ListWindows(targetSession)
			if err != nil {
				return err
			}
			for _, w := range windows {
				if w.Name == name {
					target := name
					if targetSession != "" {
						target = fmt.Sprintf("%s:%s", targetSession, name)
					}
					if _, err := base.Tmux.SelectWindow(target); err != nil {
						return fmt.Errorf("failed to select window '%s': %w", name, err)
					}
					return nil
				}
			}

			expanded, err := base.Home.ExpandHome(name)
			if err != nil {
				return err
			}
			isDir, absPath := base.Dir.Dir(expanded)
			if !isDir {
				return fmt.Errorf("'%s' is not an existing window or valid directory", name)
			}

			windowName := filepath.Base(absPath)
			newWindowArgs := []string{"new-window", "-n", windowName, "-c", absPath}
			if targetSession != "" {
				newWindowArgs = append(newWindowArgs, "-t", targetSession)
			}
			if _, err := base.Shell.Cmd("tmux", newWindowArgs...); err != nil {
				return fmt.Errorf("failed to create window: %w", err)
			}
			return nil
		},
	}

	cmd.Flags().StringP("session", "s", "", "target session (default: current attached session)")
	cmd.Flags().BoolP("json", "j", false, "output as json (list mode only)")

	return cmd
}
