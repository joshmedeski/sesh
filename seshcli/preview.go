package seshcli

import (
	"context"
	"errors"
	"fmt"
	"io"
	"time"

	"github.com/spf13/cobra"
)

const previewClearScreen = "\033[2J\033[H"

type previewFunc func(string) (string, error)

func writePreview(w io.Writer, preview previewFunc, name string, clear bool) error {
	output, err := preview(name)
	if err != nil {
		return err
	}

	if clear {
		output = previewClearScreen + output
	}

	_, err = fmt.Fprint(w, output)
	return err
}

func watchPreview(ctx context.Context, w io.Writer, preview previewFunc, name string, ticks <-chan time.Time) error {
	clear := false
	for {
		if err := writePreview(w, preview, name, clear); err != nil {
			return err
		}
		clear = true

		select {
		case <-ctx.Done():
			return nil
		case _, ok := <-ticks:
			if !ok {
				return nil
			}
		}
	}
}

func NewPreviewCommand(base *BaseDeps) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "preview",
		Aliases: []string{"p"},
		Short:   "Preview a session or directory",
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) != 1 {
				return errors.New("session name or directory is required")
			}

			deps, err := buildDeps(cmd, base)
			if err != nil {
				return err
			}

			name := args[0]
			watch, _ := cmd.Flags().GetBool("watch")
			interval, _ := cmd.Flags().GetDuration("interval")

			if !watch {
				return writePreview(cmd.OutOrStdout(), deps.Previewer.Preview, name, false)
			}
			if interval <= 0 {
				return errors.New("preview interval must be greater than zero")
			}

			_, sessionExists := deps.Lister.FindTmuxSession(deps.Icon.RemoveIcon(name))
			if !sessionExists {
				return writePreview(cmd.OutOrStdout(), deps.Previewer.Preview, name, false)
			}

			ticker := time.NewTicker(interval)
			defer ticker.Stop()

			return watchPreview(cmd.Context(), cmd.OutOrStdout(), deps.Previewer.Preview, name, ticker.C)
		},
	}

	cmd.Flags().BoolP("watch", "w", false, "continuously refresh the preview")
	cmd.Flags().Duration("interval", 500*time.Millisecond, "refresh interval when watching")

	return cmd
}
