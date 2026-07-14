package mkdirer

import (
	"errors"
	"os"
	"testing"

	"github.com/joshmedeski/sesh/v2/connector"
	"github.com/joshmedeski/sesh/v2/home"
	"github.com/joshmedeski/sesh/v2/model"
	"github.com/joshmedeski/sesh/v2/oswrap"
	"github.com/stretchr/testify/assert"
)

func TestMkdir(t *testing.T) {
	t.Run("creates a relative path directory and connects to it", func(t *testing.T) {
		mockOs := new(oswrap.MockOs)
		mockHome := new(home.MockHome)
		mockConnector := new(connector.MockConnector)

		mockHome.EXPECT().ExpandPath("my-project").Return("my-project", nil)
		mockOs.EXPECT().MkdirAll("my-project", os.FileMode(0o755)).Return(nil)
		mockConnector.EXPECT().Connect("my-project", model.ConnectOpts{}).Return("session", nil)

		m := NewMkdirer(mockOs, mockHome, mockConnector)
		name, err := m.Mkdir("my-project", model.ConnectOpts{})

		assert.NoError(t, err)
		assert.Equal(t, "session", name)
	})

	t.Run("creates an absolute path directory and connects to it", func(t *testing.T) {
		mockOs := new(oswrap.MockOs)
		mockHome := new(home.MockHome)
		mockConnector := new(connector.MockConnector)

		mockHome.EXPECT().ExpandPath("/tmp/my-project").Return("/tmp/my-project", nil)
		mockOs.EXPECT().MkdirAll("/tmp/my-project", os.FileMode(0o755)).Return(nil)
		mockConnector.EXPECT().Connect("/tmp/my-project", model.ConnectOpts{}).Return("session", nil)

		m := NewMkdirer(mockOs, mockHome, mockConnector)
		name, err := m.Mkdir("/tmp/my-project", model.ConnectOpts{})

		assert.NoError(t, err)
		assert.Equal(t, "session", name)
	})

	t.Run("expands a tilde path before creating the directory", func(t *testing.T) {
		mockOs := new(oswrap.MockOs)
		mockHome := new(home.MockHome)
		mockConnector := new(connector.MockConnector)

		mockHome.EXPECT().ExpandPath("~/my-project").Return("/home/user/my-project", nil)
		mockOs.EXPECT().MkdirAll("/home/user/my-project", os.FileMode(0o755)).Return(nil)
		mockConnector.EXPECT().Connect("/home/user/my-project", model.ConnectOpts{}).Return("session", nil)

		m := NewMkdirer(mockOs, mockHome, mockConnector)
		name, err := m.Mkdir("~/my-project", model.ConnectOpts{})

		assert.NoError(t, err)
		assert.Equal(t, "session", name)
	})

	t.Run("returns error when ExpandPath fails", func(t *testing.T) {
		mockOs := new(oswrap.MockOs)
		mockHome := new(home.MockHome)
		mockConnector := new(connector.MockConnector)

		mockHome.EXPECT().ExpandPath("my-project").Return("", errors.New("no home"))

		m := NewMkdirer(mockOs, mockHome, mockConnector)
		_, err := m.Mkdir("my-project", model.ConnectOpts{})

		assert.Error(t, err)
	})

	t.Run("returns error when MkdirAll fails", func(t *testing.T) {
		mockOs := new(oswrap.MockOs)
		mockHome := new(home.MockHome)
		mockConnector := new(connector.MockConnector)

		mockHome.EXPECT().ExpandPath("my-project").Return("my-project", nil)
		mockOs.EXPECT().MkdirAll("my-project", os.FileMode(0o755)).Return(errors.New("permission denied"))

		m := NewMkdirer(mockOs, mockHome, mockConnector)
		_, err := m.Mkdir("my-project", model.ConnectOpts{})

		assert.Error(t, err)
	})

	t.Run("returns error when Connect fails", func(t *testing.T) {
		mockOs := new(oswrap.MockOs)
		mockHome := new(home.MockHome)
		mockConnector := new(connector.MockConnector)

		mockHome.EXPECT().ExpandPath("my-project").Return("my-project", nil)
		mockOs.EXPECT().MkdirAll("my-project", os.FileMode(0o755)).Return(nil)
		mockConnector.EXPECT().Connect("my-project", model.ConnectOpts{}).Return("", errors.New("connect failed"))

		m := NewMkdirer(mockOs, mockHome, mockConnector)
		_, err := m.Mkdir("my-project", model.ConnectOpts{})

		assert.Error(t, err)
	})
}
