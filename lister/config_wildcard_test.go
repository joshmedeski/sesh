package lister

import (
	"log"
	"testing"

	"github.com/joshmedeski/sesh/v2/home"
	"github.com/joshmedeski/sesh/v2/model"
	"github.com/joshmedeski/sesh/v2/tmux"
	"github.com/joshmedeski/sesh/v2/tmuxinator"
	"github.com/joshmedeski/sesh/v2/zoxide"
	"github.com/stretchr/testify/assert"
)

func TestFindConfigWildcard(t *testing.T) {
	mockHome := new(home.MockHome)
	mockZoxide := new(zoxide.MockZoxide)
	mockTmux := new(tmux.MockTmux)
	mockTmuxinator := new(tmuxinator.MockTmuxinator)

	config := model.Config{
		WildcardConfigs: []model.WildcardConfig{
			{
				Pattern:        "~/projects/*",
				StartupCommand: "nvim",
			},
			{
				Pattern:        "~/work/*",
				StartupCommand: "make dev",
				PreviewCommand: "ls -la",
			},
		},
	}
	lister := NewLister(config, mockHome, mockTmux, mockZoxide, mockTmuxinator)
	realLister, ok := lister.(*RealLister)
	if !ok {
		log.Fatal("Cannot convert lister to *RealLister")
	}

	// Register pattern expansions used across multiple test cases
	mockHome.On("ExpandHome", "~/projects/*").Return("/Users/test/projects/*", nil)
	mockHome.On("ExpandHome", "~/work/*").Return("/Users/test/work/*", nil)

	t.Run("should match path against wildcard pattern", func(t *testing.T) {
		mockHome.On("ExpandHome", "/Users/test/projects/myapp").Return("/Users/test/projects/myapp", nil)
		wc, found := realLister.FindConfigWildcard("/Users/test/projects/myapp")
		assert.True(t, found)
		assert.Equal(t, "nvim", wc.StartupCommand)
		assert.Equal(t, "~/projects/*", wc.Pattern)
	})

	t.Run("should return false when no pattern matches", func(t *testing.T) {
		mockHome.On("ExpandHome", "/Users/test/other/myapp").Return("/Users/test/other/myapp", nil)
		wc, found := realLister.FindConfigWildcard("/Users/test/other/myapp")
		assert.False(t, found)
		assert.Equal(t, model.WildcardConfig{}, wc)
	})

	t.Run("should match second wildcard pattern", func(t *testing.T) {
		mockHome.On("ExpandHome", "/Users/test/work/project").Return("/Users/test/work/project", nil)
		wc, found := realLister.FindConfigWildcard("/Users/test/work/project")
		assert.True(t, found)
		assert.Equal(t, "make dev", wc.StartupCommand)
		assert.Equal(t, "ls -la", wc.PreviewCommand)
	})

	t.Run("should not match nested paths (single-level glob)", func(t *testing.T) {
		mockHome.On("ExpandHome", "/Users/test/projects/foo/bar").Return("/Users/test/projects/foo/bar", nil)
		_, found := realLister.FindConfigWildcard("/Users/test/projects/foo/bar")
		assert.False(t, found)
	})
}
