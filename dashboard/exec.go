package dashboard

import (
	"github.com/joshmedeski/sesh/v2/dashboard/core"
	"github.com/joshmedeski/sesh/v2/execwrap"
)

// CommandRunner is defined in dashboard/core and re-exported here so existing
// callers keep working.
type CommandRunner = core.CommandRunner

// NewCommandRunner returns a CommandRunner that executes commands via exec.
func NewCommandRunner(exec execwrap.Exec) CommandRunner {
	return core.NewCommandRunner(exec)
}
