package core

import (
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/joshmedeski/sesh/v2/execwrap"
)

func TestRunShellReturnsStdoutOnly(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("uses sh")
	}
	r := NewCommandRunner(execwrap.NewExec())

	out, err := r.RunShell("echo out; echo progress >&2")
	require.NoError(t, err)
	assert.Equal(t, "out\n", string(out))

	_, err = r.RunShell("echo boom >&2; exit 3")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "boom")
}
