package configurator

import (
	"testing"

	"github.com/joshmedeski/sesh/v2/model"
	"github.com/pelletier/go-toml/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWorktreeConfigParsing(t *testing.T) {
	input := `
terminal = "wezterm"

[[worktree]]
name = "nutiliti/nutiliti"
path = "~/c/nu"
worktree_dir = "w"
branch_template = "jam/{number}-1"
base_branch = "origin/main"
fetch = false
startup_command = "nu_setup"

[[worktree]]
name = "joshmedeski/sesh"
path = "~/c/sesh"
worktree_dir = "w"
`
	var c model.Config
	require.NoError(t, toml.Unmarshal([]byte(input), &c))

	assert.Equal(t, "wezterm", c.Terminal)
	require.Len(t, c.WorktreeConfigs, 2)

	nu := c.WorktreeConfigs[0]
	assert.Equal(t, "nutiliti/nutiliti", nu.Name)
	assert.Equal(t, "~/c/nu", nu.Path)
	assert.Equal(t, "w", nu.WorktreeDir)
	assert.Equal(t, "jam/{number}-1", nu.BranchTemplate)
	assert.Equal(t, "origin/main", nu.BaseBranch)
	require.NotNil(t, nu.Fetch)
	assert.False(t, *nu.Fetch)
	assert.Equal(t, "nu_setup", nu.StartupCommand)

	sesh := c.WorktreeConfigs[1]
	assert.Equal(t, "joshmedeski/sesh", sesh.Name)
	assert.Nil(t, sesh.Fetch) // absent => nil => treated as true downstream
}
