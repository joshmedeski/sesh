package seshcli

import (
	"errors"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/joshmedeski/sesh/v2/model"
)

func NewCloneCommand(base *BaseDeps) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "clone <repo> [-- <extra clone flags>]",
		Aliases: []string{"cl"},
		Short:   "Clone a repo with git or jj and connect to it as a session",
		Args:    cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("please provide url to clone")
			}
			repo := args[0]

			deps, err := buildDeps(cmd, base)
			if err != nil {
				return err
			}

			cmdDir, _ := cmd.Flags().GetString("cmdDir")
			dir, _ := cmd.Flags().GetString("dir")
			vcsFlag, _ := cmd.Flags().GetString("vcs")
			vcs := cloneVcs(vcsFlag, deps.Config.Vcs)
			switchFlag, _ := cmd.Flags().GetBool("switch")
			command, _ := cmd.Flags().GetString("command")
			background, _ := cmd.Flags().GetBool("background")

			var gitFlags []string
			dashIdx := cmd.ArgsLenAtDash()
			if dashIdx != -1 {
				gitFlags = args[dashIdx:]
			}

			opts := model.GitCloneOptions{
				CmdDir:   cmdDir,
				Repo:     repo,
				Dir:      dir,
				GitFlags: gitFlags,
				Vcs:      vcs,
				Connect:  model.ConnectOpts{Switch: switchFlag, Command: command, Background: background},
			}
			session, err := deps.Cloner.Clone(opts)
			if err != nil {
				return err
			}
			if background {
				fmt.Println(session)
			}
			return nil
		},
	}

	cmd.Flags().StringP("cmdDir", "c", "", "The directory to clone into")
	cmd.Flags().StringP("dir", "d", "", "The name of the directory the clone creates")
	cmd.Flags().String("vcs", "", "Clone with {git|jj} (overrides the vcs config)")
	cmd.Flags().BoolP("switch", "s", false, "Switch the session (rather than attach). This is useful for actions triggered outside the terminal.")
	cmd.Flags().String("command", "", "Execute a command when connecting to the new session.")
	cmd.Flags().BoolP("background", "b", false, "create without selecting it or moving the client")

	return cmd
}

func cloneVcs(flag string, configured string) string {
	if flag != "" {
		return flag
	}
	return configured
}
