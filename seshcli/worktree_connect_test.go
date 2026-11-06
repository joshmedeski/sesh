package seshcli

import (
	"io"
	"testing"

	"github.com/stretchr/testify/require"
)

// The connect command accepts exactly one source for the issue/PR reference:
// a positional reference, --browser, or --clipboard.
func TestWorktreeConnectRejectsCombinedSources(t *testing.T) {
	tests := []struct {
		name string
		args []string
	}{
		{"no reference and no source flag", []string{}},
		{"browser and clipboard together", []string{"--browser", "--clipboard"}},
		{"browser with a positional reference", []string{"--browser", "409"}},
		{"clipboard with a positional reference", []string{"--clipboard", "409"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Every invalid shape is rejected before a dependency is built,
			// so an empty BaseDeps is enough to reach the validation.
			cmd := newWorktreeConnectCommand(&BaseDeps{})
			cmd.SetArgs(tt.args)
			cmd.SetOut(io.Discard)
			cmd.SetErr(io.Discard)

			err := cmd.Execute()
			require.ErrorContains(t, err, "provide an issue/PR reference")
			if len(tt.args) > 1 {
				require.ErrorContains(t, err, ", not a combination")
			}
		})
	}
}
