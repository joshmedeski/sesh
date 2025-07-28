package namer

import (
	"testing"

	"github.com/joshmedeski/sesh/v2/git"
	"github.com/joshmedeski/sesh/v2/home"
	"github.com/joshmedeski/sesh/v2/model"
	"github.com/joshmedeski/sesh/v2/pathwrap"
	"github.com/stretchr/testify/assert"
)

func TestDirName(t *testing.T) {
	mockPathwrap := new(pathwrap.MockPath)
	mockGit := new(git.MockGit)
	mockHome := new(home.MockHome)

	tests := []struct {
		name     string
		path     string
		dirLen   int
		expected string
	}{
		{"dir_length 1", "/Users/john/projects/sesh", 1, "sesh"},
		{"dir_length 2", "/Users/john/projects/sesh", 2, "projects/sesh"},
		{"dir_length 3", "/Users/john/projects/sesh", 3, "john/projects/sesh"},
		{"folder1/path with dir_length 2", "/folder1/path", 2, "folder1/path"},
		{"folder2/path with dir_length 2", "/folder2/path", 2, "folder2/path"},
		{"short path with dir_length 5", "/one/two", 5, "one/two"},
		{"single dir with dir_length 3", "/path", 3, "path"},
		{"zero dir_length", "/some/path", 0, "path"},
		{"negative dir_length", "/some/path", -1, "path"},
		{"root path", "/", 1, "/"},
		{"trailing slash", "/Users/john/projects/sesh/", 2, "projects/sesh"},
		{"relative path", "projects/sesh/subdir", 2, "sesh/subdir"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			config := model.Config{DirLength: test.dirLen}
			n := &RealNamer{pathwrap: mockPathwrap, git: mockGit, home: mockHome, config: config}

			// Set up Base mock for dir_length <= 1 cases
			if test.dirLen <= 1 {
				mockPathwrap.On("Base", test.path).Return(test.expected)
			}

			result, err := dirName(n, test.path)
			assert.NoError(t, err)
			assert.Equal(t, test.expected, result)
		})
	}
}
