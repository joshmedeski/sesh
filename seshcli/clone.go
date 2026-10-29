package seshcli

import (
	"errors"

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

			var gitFlags []string
			dashIdx := cmd.ArgsLenAtDash()
			if dashIdx != -1 {
				gitFlags = args[dashIdx:]
			}

			opts := model.GitCloneOptions{CmdDir: cmdDir, Repo: repo, Dir: dir, GitFlags: gitFlags, Vcs: vcs}
			if _, err := deps.Cloner.Clone(opts); err != nil {
				return err
			} else {
				return nil
			}
		},
	}

	cmd.Flags().StringP("cmdDir", "c", "", "The directory to clone into")
	cmd.Flags().StringP("dir", "d", "", "The name of the directory the clone creates")
	cmd.Flags().String("vcs", "", "Clone with {git|jj} (overrides the vcs config)")

	return cmd
}

func cloneVcs(flag string, configured string) string {
	if flag != "" {
		return flag
	}
	return configured
}
