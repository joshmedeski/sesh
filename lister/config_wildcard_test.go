package lister

import (
	"log"
	"testing"

	"github.com/Wingsdh/cc-sesh/v2/home"
	"github.com/Wingsdh/cc-sesh/v2/model"
	"github.com/Wingsdh/cc-sesh/v2/tmux"
	"github.com/Wingsdh/cc-sesh/v2/tmuxinator"
	"github.com/Wingsdh/cc-sesh/v2/zoxide"
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
			{
				Pattern:        "~/deep/**",
				StartupCommand: "deep-cmd",
			},
		},
	}
	lister := NewLister(config, mockHome, mockTmux, mockZoxide, mockTmuxinator)
	realLister, ok := lister.(*RealLister)
	if !ok {
		log.Fatal("Cannot convert lister to *RealLister")
	}

	// Register pattern expansions used across multiple test cases
	mockHome.On("ExpandPath","~/projects/*").Return("/Users/test/projects/*", nil)
	mockHome.On("ExpandPath","~/work/*").Return("/Users/test/work/*", nil)
	mockHome.On("ExpandPath","~/deep/**").Return("/Users/test/deep/**", nil)

	t.Run("should match path against wildcard pattern", func(t *testing.T) {
		mockHome.On("ExpandPath","/Users/test/projects/myapp").Return("/Users/test/projects/myapp", nil)
		wc, found := realLister.FindConfigWildcard("/Users/test/projects/myapp")
		assert.True(t, found)
		assert.Equal(t, "nvim", wc.StartupCommand)
		assert.Equal(t, "~/projects/*", wc.Pattern)
	})

	t.Run("should return false when no pattern matches", func(t *testing.T) {
		mockHome.On("ExpandPath","/Users/test/other/myapp").Return("/Users/test/other/myapp", nil)
		wc, found := realLister.FindConfigWildcard("/Users/test/other/myapp")
		assert.False(t, found)
		assert.Equal(t, model.WildcardConfig{}, wc)
	})

	t.Run("should match second wildcard pattern", func(t *testing.T) {
		mockHome.On("ExpandPath","/Users/test/work/project").Return("/Users/test/work/project", nil)
		wc, found := realLister.FindConfigWildcard("/Users/test/work/project")
		assert.True(t, found)
		assert.Equal(t, "make dev", wc.StartupCommand)
		assert.Equal(t, "ls -la", wc.PreviewCommand)
	})

	t.Run("should not match nested paths (single-level glob)", func(t *testing.T) {
		mockHome.On("ExpandPath","/Users/test/projects/foo/bar").Return("/Users/test/projects/foo/bar", nil)
		_, found := realLister.FindConfigWildcard("/Users/test/projects/foo/bar")
		assert.False(t, found)
	})

	t.Run("should match nested paths with ** pattern", func(t *testing.T) {
		mockHome.On("ExpandPath","/Users/test/deep/foo/bar/baz").Return("/Users/test/deep/foo/bar/baz", nil)
		wc, found := realLister.FindConfigWildcard("/Users/test/deep/foo/bar/baz")
		assert.True(t, found)
		assert.Equal(t, "deep-cmd", wc.StartupCommand)
	})

	t.Run("should not match unrelated paths with ** pattern", func(t *testing.T) {
		mockHome.On("ExpandPath","/Users/test/other/foo").Return("/Users/test/other/foo", nil)
		_, found := realLister.FindConfigWildcard("/Users/test/other/foo")
		assert.False(t, found)
	})

	t.Run("should not match the prefix directory itself with **", func(t *testing.T) {
		mockHome.On("ExpandPath","/Users/test/deep").Return("/Users/test/deep", nil)
		_, found := realLister.FindConfigWildcard("/Users/test/deep")
		assert.False(t, found)
	})

	t.Run("should not match the prefix directory with trailing slash and **", func(t *testing.T) {
		mockHome.On("ExpandPath","/Users/test/deep/").Return("/Users/test/deep/", nil)
		_, found := realLister.FindConfigWildcard("/Users/test/deep/")
		assert.False(t, found)
	})
}
