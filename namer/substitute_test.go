package namer

import (
	"testing"

	"github.com/joshmedeski/sesh/v2/git"
	"github.com/joshmedeski/sesh/v2/home"
	"github.com/joshmedeski/sesh/v2/model"
	"github.com/joshmedeski/sesh/v2/pathwrap"
	"github.com/stretchr/testify/assert"
)

func TestApplySubstitutions(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		rules    []model.NameSubstitution
		expected string
	}{
		{
			name:     "no rules leaves the name untouched",
			input:    "~/c/dotfiles/.config/nvim",
			rules:    nil,
			expected: "~/c/dotfiles/.config/nvim",
		},
		{
			name:     "literal prefix replacement",
			input:    "~/c/dotfiles/.config/nvim",
			rules:    []model.NameSubstitution{{Find: "~/c/dotfiles/.config/", Replace: ""}},
			expected: "nvim",
		},
		{
			name:     "literal replaces every occurrence",
			input:    "~/work/work/api",
			rules:    []model.NameSubstitution{{Find: "work/", Replace: ""}},
			expected: "~/api",
		},
		{
			name:  "rules apply in order and compose",
			input: "~/c/projects/api",
			rules: []model.NameSubstitution{
				{Find: "~/c/projects/", Replace: "proj:"},
				{Find: "proj:", Replace: "🚀 "},
			},
			expected: "🚀 api",
		},
		{
			name:  "an empty find is skipped",
			input: "~/c/api",
			rules: []model.NameSubstitution{
				{Find: "", Replace: "ignored"},
				{Find: "~/c/", Replace: ""},
			},
			expected: "api",
		},
		{
			name:     "literal find treats dots as literal, not wildcards",
			input:    "~/c/dotXconfig/nvim",
			rules:    []model.NameSubstitution{{Find: "~/c/dot.config/", Replace: ""}},
			expected: "~/c/dotXconfig/nvim",
		},
		{
			name:     "regex with a capture group",
			input:    "~/c/workspace/1_docs",
			rules:    []model.NameSubstitution{{Find: `.*/workspace/[0-9]+_(.*)`, Replace: "ws-$1", Regex: true}},
			expected: "ws-docs",
		},
		{
			name:     "regex metacharacters match when regex is enabled",
			input:    "~/c/dotXconfig/nvim",
			rules:    []model.NameSubstitution{{Find: `~/c/dot.config/`, Replace: "", Regex: true}},
			expected: "nvim",
		},
		{
			name:     "an invalid regex is skipped rather than applied",
			input:    "~/c/api",
			rules:    []model.NameSubstitution{{Find: `(unterminated`, Replace: "x", Regex: true}},
			expected: "~/c/api",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert.Equal(t, test.expected, applySubstitutions(test.input, test.rules))
		})
	}
}

func TestNameSubstitutionStrategy(t *testing.T) {
	t.Run("returns empty when no rules are configured", func(t *testing.T) {
		mockHome := new(home.MockHome)
		n := &RealNamer{home: mockHome, config: model.Config{}}

		name, err := nameSubstitution(n, "/home/john/c/nvim")

		assert.NoError(t, err)
		assert.Equal(t, "", name)
		mockHome.AssertNotCalled(t, "ShortenHome")
	})

	t.Run("returns the rewritten name when a rule matches", func(t *testing.T) {
		mockHome := new(home.MockHome)
		mockHome.On("ShortenHome", "/home/john/c/dotfiles/.config/nvim").
			Return("~/c/dotfiles/.config/nvim", nil)
		config := model.Config{NameSubstitutions: []model.NameSubstitution{
			{Find: "~/c/dotfiles/.config/", Replace: ""},
		}}
		n := &RealNamer{home: mockHome, config: config}

		name, err := nameSubstitution(n, "/home/john/c/dotfiles/.config/nvim")

		assert.NoError(t, err)
		assert.Equal(t, "nvim", name)
	})

	t.Run("falls through when no rule changes the path", func(t *testing.T) {
		mockHome := new(home.MockHome)
		mockHome.On("ShortenHome", "/home/john/c/api").Return("~/c/api", nil)
		config := model.Config{NameSubstitutions: []model.NameSubstitution{
			{Find: "~/nowhere/", Replace: ""},
		}}
		n := &RealNamer{home: mockHome, config: config}

		name, err := nameSubstitution(n, "/home/john/c/api")

		assert.NoError(t, err)
		assert.Equal(t, "", name)
	})
}

func TestNameUsesSubstitutionBeforeDirStrategy(t *testing.T) {
	mockPathwrap := new(pathwrap.MockPath)
	mockGit := new(git.MockGit)
	mockHome := new(home.MockHome)

	path := "/home/john/c/dotfiles/.config/nvim"
	mockPathwrap.On("EvalSymlinks", path).Return(path, nil)
	mockHome.On("ShortenHome", path).Return("~/c/dotfiles/.config/nvim", nil)

	config := model.Config{
		DirLength: 1,
		NameSubstitutions: []model.NameSubstitution{
			{Find: "~/c/dotfiles/.config/", Replace: "cfg "},
		},
	}
	n := &RealNamer{pathwrap: mockPathwrap, git: mockGit, home: mockHome, config: config}

	name, err := n.Name(path)

	assert.NoError(t, err)
	// convertToValidName collapses the space, so "cfg nvim" becomes "cfg_nvim".
	assert.Equal(t, "cfg_nvim", name)
	mockGit.AssertNotCalled(t, "ShowTopLevel")
}
