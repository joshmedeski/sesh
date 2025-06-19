package seshcli

import (
	"errors"

	"github.com/spf13/cobra"

	"github.com/joshmedeski/sesh/v2/cloner"
	"github.com/joshmedeski/sesh/v2/model"
)

func NewCloneCommand(c cloner.Cloner) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "clone",
		Aliases: []string{"cl"},
		Short:   "Clone a git repo and connect to it as a session",
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) != 1 {
				return errors.New("please provide url to clone")
			}
			repo := args[0]

			cmdDir, _ := cmd.Flags().GetString("cmdDir")
			dir, _ := cmd.Flags().GetString("dir")

			opts := model.GitCloneOptions{CmdDir: cmdDir, Repo: repo, Dir: dir}
			if _, err := c.Clone(opts); err != nil {
				return err
			} else {
				return nil
			}
		},
	}

	cmd.Flags().StringP("cmdDir", "c", "", "The directory to run the git command in")
	cmd.Flags().StringP("dir", "d", "", "The name of the directory that git is creating")

	return cmd
}
