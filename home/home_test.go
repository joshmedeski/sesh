package home

import (
	"errors"
	"testing"

	"github.com/joshmedeski/sesh/v2/oswrap"
	"github.com/stretchr/testify/assert"
)

func TestExpandPath(t *testing.T) {
	t.Run("returns plain path unchanged", func(t *testing.T) {
		mockOs := new(oswrap.MockOs)
		mockOs.On("ExpandEnv", "/foo/bar").Return("/foo/bar")
		mockOs.On("UserHomeDir").Return("/home/user", nil)

		got, err := NewHome(mockOs).ExpandPath("/foo/bar")
		assert.NoError(t, err)
		assert.Equal(t, "/foo/bar", got)
	})

	t.Run("expands leading tilde to home dir", func(t *testing.T) {
		mockOs := new(oswrap.MockOs)
		mockOs.On("ExpandEnv", "~/projects").Return("~/projects")
		mockOs.On("UserHomeDir").Return("/home/user", nil)

		got, err := NewHome(mockOs).ExpandPath("~/projects")
		assert.NoError(t, err)
		assert.Equal(t, "/home/user/projects", got)
	})

	t.Run("expands $VAR syntax", func(t *testing.T) {
		mockOs := new(oswrap.MockOs)
		mockOs.On("ExpandEnv", "$WORK/app").Return("/opt/work/app")
		mockOs.On("UserHomeDir").Return("/home/user", nil)

		got, err := NewHome(mockOs).ExpandPath("$WORK/app")
		assert.NoError(t, err)
		assert.Equal(t, "/opt/work/app", got)
	})

	t.Run("expands ${VAR} syntax", func(t *testing.T) {
		mockOs := new(oswrap.MockOs)
		mockOs.On("ExpandEnv", "${WORK}/app").Return("/opt/work/app")
		mockOs.On("UserHomeDir").Return("/home/user", nil)

		got, err := NewHome(mockOs).ExpandPath("${WORK}/app")
		assert.NoError(t, err)
		assert.Equal(t, "/opt/work/app", got)
	})

	t.Run("expands env var then tilde when combined", func(t *testing.T) {
		// env expansion happens first; if result starts with ~, tilde expansion follows
		mockOs := new(oswrap.MockOs)
		mockOs.On("ExpandEnv", "$PREFIX/sub").Return("~/sub")
		mockOs.On("UserHomeDir").Return("/home/user", nil)

		got, err := NewHome(mockOs).ExpandPath("$PREFIX/sub")
		assert.NoError(t, err)
		assert.Equal(t, "/home/user/sub", got)
	})

	t.Run("expands env var inside tilde-prefixed path", func(t *testing.T) {
		mockOs := new(oswrap.MockOs)
		mockOs.On("ExpandEnv", "~/$PROJECT").Return("~/myapp")
		mockOs.On("UserHomeDir").Return("/home/user", nil)

		got, err := NewHome(mockOs).ExpandPath("~/$PROJECT")
		assert.NoError(t, err)
		assert.Equal(t, "/home/user/myapp", got)
	})

	t.Run("does not expand tilde mid-path", func(t *testing.T) {
		mockOs := new(oswrap.MockOs)
		mockOs.On("ExpandEnv", "/foo/~/bar").Return("/foo/~/bar")
		mockOs.On("UserHomeDir").Return("/home/user", nil)

		got, err := NewHome(mockOs).ExpandPath("/foo/~/bar")
		assert.NoError(t, err)
		assert.Equal(t, "/foo/~/bar", got)
	})

	t.Run("returns error when UserHomeDir fails", func(t *testing.T) {
		mockOs := new(oswrap.MockOs)
		mockOs.On("ExpandEnv", "~/x").Return("~/x")
		mockOs.On("UserHomeDir").Return("", errors.New("no home"))

		_, err := NewHome(mockOs).ExpandPath("~/x")
		assert.Error(t, err)
	})
}

func TestShortenHome(t *testing.T) {
	t.Run("replaces home dir prefix with tilde", func(t *testing.T) {
		mockOs := new(oswrap.MockOs)
		mockOs.On("UserHomeDir").Return("/home/user", nil)

		got, err := NewHome(mockOs).ShortenHome("/home/user/projects")
		assert.NoError(t, err)
		assert.Equal(t, "~/projects", got)
	})

	t.Run("returns path unchanged when home dir is not a prefix", func(t *testing.T) {
		mockOs := new(oswrap.MockOs)
		mockOs.On("UserHomeDir").Return("/home/user", nil)

		got, err := NewHome(mockOs).ShortenHome("/var/log")
		assert.NoError(t, err)
		assert.Equal(t, "/var/log", got)
	})

	t.Run("returns error when UserHomeDir fails", func(t *testing.T) {
		mockOs := new(oswrap.MockOs)
		mockOs.On("UserHomeDir").Return("", errors.New("no home"))

		_, err := NewHome(mockOs).ShortenHome("/whatever")
		assert.Error(t, err)
	})
}
