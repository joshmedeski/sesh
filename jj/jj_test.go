package jj

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/joshmedeski/sesh/v2/oswrap"
	"github.com/joshmedeski/sesh/v2/shell"
	"github.com/stretchr/testify/assert"
)

func TestRoots(t *testing.T) {
	root := t.TempDir()
	repo := filepath.Join(root, "project")
	assert.NoError(t, os.MkdirAll(filepath.Join(repo, ".jj", "repo"), 0o755))
	assert.NoError(t, os.MkdirAll(filepath.Join(repo, "src", "pkg"), 0o755))

	absWorkspace := filepath.Join(root, "feat")
	assert.NoError(t, os.MkdirAll(filepath.Join(absWorkspace, ".jj"), 0o755))
	assert.NoError(t, os.WriteFile(filepath.Join(absWorkspace, ".jj", "repo"), []byte(filepath.Join(repo, ".jj", "repo")), 0o644))

	relWorkspace := filepath.Join(repo, "workspaces", "fix")
	assert.NoError(t, os.MkdirAll(filepath.Join(relWorkspace, ".jj"), 0o755))
	assert.NoError(t, os.WriteFile(filepath.Join(relWorkspace, ".jj", "repo"), []byte("../../../.jj/repo\n"), 0o644))

	j := NewJj(oswrap.NewOs(), nil)

	tests := []struct {
		name      string
		path      string
		isJj      bool
		workspace string
		repo      string
	}{
		{"repo root", repo, true, repo, repo},
		{"nested subdir", filepath.Join(repo, "src", "pkg"), true, repo, repo},
		{"workspace with absolute repo pointer", absWorkspace, true, absWorkspace, repo},
		{"workspace with relative repo pointer", relWorkspace, true, relWorkspace, repo},
		{"not a jj repo", root, false, "", ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			isJj, workspace, repoRoot := j.Roots(tt.path)
			assert.Equal(t, tt.isJj, isJj)
			assert.Equal(t, tt.workspace, workspace)
			assert.Equal(t, tt.repo, repoRoot)
		})
	}
}

func TestClone(t *testing.T) {
	mockShell := new(shell.MockShell)
	mockShell.On("CmdWithOutput", "jj", "git", "clone", "--colocate", "https://github.com/joshmedeski/sesh", "/code/sesh").Return("", nil)
	j := NewJj(nil, mockShell)

	assert.NoError(t, j.Clone("https://github.com/joshmedeski/sesh", "/code/sesh", "--colocate"))
	mockShell.AssertExpectations(t)
}
