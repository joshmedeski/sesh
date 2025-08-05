package namer

import (
	"fmt"
	"strings"
	"testing"

	"github.com/joshmedeski/sesh/v2/git"
	"github.com/joshmedeski/sesh/v2/home"
	"github.com/joshmedeski/sesh/v2/model"
	"github.com/joshmedeski/sesh/v2/pathwrap"
	"github.com/stretchr/testify/assert"
)

func TestGitBareName(t *testing.T) {
	mockPathwrap := new(pathwrap.MockPath)
	mockGit := new(git.MockGit)
	mockHome := new(home.MockHome)

	tests := []struct {
		name      string
		path      string
		dirLength int
		isGit     bool
		commonDir string
		gitErr    error
		expected  string
		shouldErr bool
	}{
		{
			name:      "git bare repository",
			path:      "/Users/user/projects/sesh/feature-branch",
			dirLength: 1,
			isGit:     true,
			commonDir: "/Users/user/projects/sesh/.bare",
			expected:  "sesh/feature-branch",
		},
		{
			name:      "non-bare git repository",
			path:      "/Users/user/projects/normal-repo",
			dirLength: 1,
			isGit:     true,
			commonDir: "/Users/user/projects/normal-repo/.git",
			expected:  "",
		},
		{
			name:      "non-git directory",
			path:      "/Users/user/projects/non-git",
			dirLength: 1,
			isGit:     false,
			commonDir: "",
			gitErr:    fmt.Errorf("not a git repository"),
			expected:  "",
		},
		{
			name:      "respects dir_length configuration",
			path:      "/Users/user/projects/project/feature-branch",
			dirLength: 2,
			isGit:     true,
			commonDir: "/Users/user/projects/project/.bare",
			expected:  "projects/project/feature-branch",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			config := model.Config{DirLength: test.dirLength}
			n := &RealNamer{pathwrap: mockPathwrap, git: mockGit, home: mockHome, config: config}

			mockGit.On("GitCommonDir", test.path).Return(test.isGit, test.commonDir, test.gitErr)

			if test.isGit && test.commonDir != "" && strings.HasSuffix(test.commonDir, "/.bare") {
				topLevelDir := strings.TrimSuffix(test.commonDir, "/.bare")
				if test.dirLength <= 1 {
					mockPathwrap.On("Base", topLevelDir).Return("sesh")
				}
			}

			result, err := gitBareName(n, test.path)
			if test.shouldErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, test.expected, result)
			}
		})
	}
}

func TestGitRootName(t *testing.T) {
	mockPathwrap := new(pathwrap.MockPath)
	mockGit := new(git.MockGit)
	mockHome := new(home.MockHome)

	tests := []struct {
		name        string
		path        string
		dirLength   int
		isGit       bool
		topLevelDir string
		gitErr      error
		expected    string
		shouldErr   bool
	}{
		{
			name:        "git repository root",
			path:        "/Users/user/projects/sesh",
			dirLength:   1,
			isGit:       true,
			topLevelDir: "/Users/user/projects/sesh",
			expected:    "sesh",
		},
		{
			name:      "non-git directory",
			path:      "/Users/user/projects/non-git",
			dirLength: 1,
			isGit:     false,
			gitErr:    fmt.Errorf("not a git repository"),
			expected:  "",
		},
		{
			name:        "empty top level",
			path:        "/Users/user/projects/empty",
			dirLength:   1,
			isGit:       true,
			topLevelDir: "",
			expected:    "",
		},
		{
			name:        "respects dir_length configuration",
			path:        "/Users/user/projects/my-very-long-project-name",
			dirLength:   3,
			isGit:       true,
			topLevelDir: "/Users/user/projects/my-very-long-project-name",
			expected:    "user/projects/my-very-long-project-name",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			config := model.Config{DirLength: test.dirLength}
			n := &RealNamer{pathwrap: mockPathwrap, git: mockGit, home: mockHome, config: config}

			mockGit.On("ShowTopLevel", test.path).Return(test.isGit, test.topLevelDir, test.gitErr)
			if test.dirLength <= 1 && test.topLevelDir != "" {
				mockPathwrap.On("Base", test.topLevelDir).Return(test.expected)
			}

			result, err := gitRootName(n, test.path)
			if test.shouldErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, test.expected, result)
			}
		})
	}
}

func TestGitName(t *testing.T) {
	mockPathwrap := new(pathwrap.MockPath)
	mockGit := new(git.MockGit)
	mockHome := new(home.MockHome)

	tests := []struct {
		name        string
		path        string
		dirLength   int
		isGit       bool
		topLevelDir string
		gitErr      error
		expected    string
		shouldErr   bool
	}{
		{
			name:        "git repository with subdirectory",
			path:        "/Users/user/projects/sesh/cmd/sesh",
			dirLength:   1,
			isGit:       true,
			topLevelDir: "/Users/user/projects/sesh",
			expected:    "sesh/cmd/sesh",
		},
		{
			name:        "git repository root directory",
			path:        "/Users/user/projects/sesh",
			dirLength:   1,
			isGit:       true,
			topLevelDir: "/Users/user/projects/sesh",
			expected:    "sesh",
		},
		{
			name:      "non-git directory",
			path:      "/Users/user/projects/non-git",
			dirLength: 1,
			isGit:     false,
			gitErr:    fmt.Errorf("not a git repository"),
			expected:  "",
		},
		{
			name:        "empty top level",
			path:        "/Users/user/projects/empty",
			dirLength:   1,
			isGit:       true,
			topLevelDir: "",
			expected:    "",
		},
		{
			name:        "respects dir_length configuration",
			path:        "/Users/user/projects/project/src/main",
			dirLength:   2,
			isGit:       true,
			topLevelDir: "/Users/user/projects/project",
			expected:    "projects/project/src/main",
		},
		{
			name:        "handles deep directory structure",
			path:        "/Users/user/projects/sesh/internal/lister/tmux/session",
			dirLength:   1,
			isGit:       true,
			topLevelDir: "/Users/user/projects/sesh",
			expected:    "sesh/internal/lister/tmux/session",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			config := model.Config{DirLength: test.dirLength}
			n := &RealNamer{pathwrap: mockPathwrap, git: mockGit, home: mockHome, config: config}

			mockGit.On("ShowTopLevel", test.path).Return(test.isGit, test.topLevelDir, test.gitErr)
			if test.dirLength <= 1 && test.topLevelDir != "" {
				expectedBase := "sesh"
				if test.topLevelDir == "/Users/user/projects/project" {
					expectedBase = "project"
				}
				mockPathwrap.On("Base", test.topLevelDir).Return(expectedBase)
			}

			result, err := gitName(n, test.path)
			if test.shouldErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, test.expected, result)
			}
		})
	}
}
