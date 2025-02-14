package cloner

import (
	"os"
	"testing"

	"github.com/joshmedeski/sesh/v2/model"
	"github.com/stretchr/testify/assert"
)

func TestListSessions(t *testing.T) {
	t.Run("get path with both cmd and repo", func(t *testing.T) {
		mockOpts := model.GitCloneOptions{Repo: "https://www.github.comtest/repo.git", CmdDir: "cmdDir", Dir: "dir"}
		actual := getPath(mockOpts)
		assert.Equal(t, "cmdDir/dir", actual)
	})
	t.Run("get path cmdDir", func(t *testing.T) {
		mockOpts := model.GitCloneOptions{Repo: "https://www.github.comtest/repo.git", CmdDir: "cmdDir"}
		actual := getPath(mockOpts)
		assert.Equal(t, "cmdDir/repo", actual)
	})
	t.Run("get path dir", func(t *testing.T) {
		mockOpts := model.GitCloneOptions{Repo: "https://www.github.comtest/repo.git", Dir: "dir"}
		actual := getPath(mockOpts)
		expected, _ := os.Getwd()
		expected = expected + "/" + "dir"
		assert.Equal(t, expected, actual)
	})
	t.Run("get path with no dir or cmdDir", func(t *testing.T) {
		mockOpts := model.GitCloneOptions{Repo: "https://www.github.comtest/repo.git"}
		actual := getPath(mockOpts)
		expected, _ := os.Getwd()
		expected = expected + "/" + "repo"
		assert.Equal(t, expected, actual)
	})
}
