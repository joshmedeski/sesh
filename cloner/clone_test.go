package cloner

import (
	"errors"
	"os"
	"testing"

	"github.com/joshmedeski/sesh/v2/connector"
	"github.com/joshmedeski/sesh/v2/git"
	"github.com/joshmedeski/sesh/v2/jj"
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

func TestClone(t *testing.T) {
	t.Run("clones with jj into the computed path", func(t *testing.T) {
		mockConnector := new(connector.MockConnector)
		mockGit := new(git.MockGit)
		mockJj := new(jj.MockJj)
		c := NewCloner(mockConnector, mockGit, mockJj)
		mockJj.On("Clone", "https://github.com/joshmedeski/sesh", "/code/sesh-jj", "--no-colocate").Return(nil)
		mockConnector.On("Connect", "/code/sesh-jj", model.ConnectOpts{}).Return("", nil)

		_, err := c.Clone(model.GitCloneOptions{Repo: "https://github.com/joshmedeski/sesh", CmdDir: "/code", Dir: "sesh-jj", GitFlags: []string{"--no-colocate"}, Vcs: "jj"})

		assert.NoError(t, err)
		mockJj.AssertExpectations(t)
		mockConnector.AssertExpectations(t)
		mockGit.AssertNotCalled(t, "Clone")
	})

	t.Run("does not connect when jj clone fails", func(t *testing.T) {
		mockConnector := new(connector.MockConnector)
		mockJj := new(jj.MockJj)
		c := NewCloner(mockConnector, new(git.MockGit), mockJj)
		mockJj.On("Clone", "https://github.com/joshmedeski/sesh", "/code/sesh").Return(errors.New("boom"))

		_, err := c.Clone(model.GitCloneOptions{Repo: "https://github.com/joshmedeski/sesh", CmdDir: "/code", Vcs: "jj"})

		assert.Error(t, err)
		mockConnector.AssertNotCalled(t, "Connect")
	})

	t.Run("rejects an unknown vcs", func(t *testing.T) {
		mockConnector := new(connector.MockConnector)
		c := NewCloner(mockConnector, new(git.MockGit), new(jj.MockJj))

		_, err := c.Clone(model.GitCloneOptions{Repo: "https://github.com/joshmedeski/sesh", Vcs: "hg"})

		assert.ErrorContains(t, err, `unsupported vcs "hg"`)
		mockConnector.AssertNotCalled(t, "Connect")
	})

	t.Run("clones with git by default", func(t *testing.T) {
		mockConnector := new(connector.MockConnector)
		mockGit := new(git.MockGit)
		mockJj := new(jj.MockJj)
		c := NewCloner(mockConnector, mockGit, mockJj)
		mockGit.On("Clone", "https://github.com/joshmedeski/sesh", "/code", "sesh-git", "--depth", "1").Return("", nil)
		mockConnector.On("Connect", "/code/sesh-git", model.ConnectOpts{}).Return("", nil)

		_, err := c.Clone(model.GitCloneOptions{Repo: "https://github.com/joshmedeski/sesh", CmdDir: "/code", Dir: "sesh-git", GitFlags: []string{"--depth", "1"}})

		assert.NoError(t, err)
		mockGit.AssertExpectations(t)
		mockConnector.AssertExpectations(t)
		mockJj.AssertNotCalled(t, "Clone")
	})

	t.Run("does not connect when git clone fails", func(t *testing.T) {
		mockConnector := new(connector.MockConnector)
		mockGit := new(git.MockGit)
		c := NewCloner(mockConnector, mockGit, new(jj.MockJj))
		mockGit.On("Clone", "https://github.com/joshmedeski/sesh", "/code", "").Return("", errors.New("boom"))

		_, err := c.Clone(model.GitCloneOptions{Repo: "https://github.com/joshmedeski/sesh", CmdDir: "/code", Vcs: "git"})

		assert.Error(t, err)
		mockConnector.AssertNotCalled(t, "Connect")
	})
}
