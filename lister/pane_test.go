package lister

import (
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/joshmedeski/sesh/v2/home"
	"github.com/joshmedeski/sesh/v2/model"
	"github.com/joshmedeski/sesh/v2/oswrap"
	"github.com/joshmedeski/sesh/v2/tmux"
	"github.com/joshmedeski/sesh/v2/tmuxinator"
	"github.com/joshmedeski/sesh/v2/zoxide"
	"github.com/stretchr/testify/assert"
)

func makeTmuxPane(windowIndex int, windowName string, paneIndex int, paneTitle, paneCommand, panePath, paneID string) *model.TmuxPane {
	return &model.TmuxPane{
		WindowIndex: windowIndex,
		WindowName:  windowName,
		PaneIndex:   paneIndex,
		PaneTitle:   paneTitle,
		PaneCommand: paneCommand,
		PanePath:    panePath,
		PaneID:      paneID,
	}
}

func TestListTmuxPanes(t *testing.T) {
	hostname, _ := os.Hostname()

	t.Run("should list tmux panes", func(t *testing.T) {
		mockTmux := new(tmux.MockTmux)
		mockTmux.On("ListTmuxPanes").Return([]*model.TmuxPane{
			makeTmuxPane(0, "editor", 0, hostname, "vim", "/home/user/project", "%0"),
			makeTmuxPane(0, "editor", 1, hostname, "zsh", "/home/user/project", "%1"),
			makeTmuxPane(1, "tests", 0, hostname, "go", "/home/user/project", "%2"),
		}, nil)

		lister := NewLister(new(oswrap.MockOs), model.Config{}, new(home.MockHome), mockTmux, new(zoxide.MockZoxide), new(tmuxinator.MockTmuxinator))
		realLister, ok := lister.(*RealLister)
		if !ok {
			log.Fatal("Cannot convert lister to *RealLister")
		}

		sessions, err := listTmuxPanes(realLister)
		assert.Nil(t, err)
		assert.Equal(t, 3, len(sessions.OrderedIndex))
		assert.Equal(t, "tmux-pane:editor/%0", sessions.OrderedIndex[0])
		assert.Equal(t, "editor/vim", sessions.Directory["tmux-pane:editor/%0"].Name)
		assert.Equal(t, "tmux-pane", sessions.Directory["tmux-pane:editor/%0"].Src)
		assert.Equal(t, "tmux-pane:editor/%1", sessions.OrderedIndex[1])
		assert.Equal(t, "editor/zsh", sessions.Directory["tmux-pane:editor/%1"].Name)
		assert.Equal(t, "tmux-pane:tests/%2", sessions.OrderedIndex[2])
		assert.Equal(t, "tests/go", sessions.Directory["tmux-pane:tests/%2"].Name)
	})

	t.Run("should use pane title when explicitly set", func(t *testing.T) {
		mockTmux := new(tmux.MockTmux)
		mockTmux.On("ListTmuxPanes").Return([]*model.TmuxPane{
			makeTmuxPane(0, "editor", 0, "my-editor", "vim", "/home/user/project", "%0"),
			makeTmuxPane(0, "editor", 1, hostname, "zsh", "/home/user/project", "%1"),
		}, nil)

		lister := NewLister(new(oswrap.MockOs), model.Config{}, new(home.MockHome), mockTmux, new(zoxide.MockZoxide), new(tmuxinator.MockTmuxinator))
		realLister, ok := lister.(*RealLister)
		if !ok {
			log.Fatal("Cannot convert lister to *RealLister")
		}

		sessions, err := listTmuxPanes(realLister)
		assert.Nil(t, err)
		assert.Equal(t, "editor/my-editor", sessions.Directory["tmux-pane:editor/%0"].Name)
		assert.Equal(t, "editor/zsh", sessions.Directory["tmux-pane:editor/%1"].Name)
	})

	t.Run("should disambiguate duplicate names with suffixes", func(t *testing.T) {
		mockTmux := new(tmux.MockTmux)
		mockTmux.On("ListTmuxPanes").Return([]*model.TmuxPane{
			makeTmuxPane(0, "editor", 0, hostname, "zsh", "/home/user/project", "%0"),
			makeTmuxPane(0, "editor", 1, hostname, "zsh", "/home/user/project", "%1"),
			makeTmuxPane(1, "tests", 0, hostname, "zsh", "/tmp", "%2"),
		}, nil)

		lister := NewLister(new(oswrap.MockOs), model.Config{}, new(home.MockHome), mockTmux, new(zoxide.MockZoxide), new(tmuxinator.MockTmuxinator))
		realLister, ok := lister.(*RealLister)
		if !ok {
			log.Fatal("Cannot convert lister to *RealLister")
		}

		sessions, err := listTmuxPanes(realLister)
		assert.Nil(t, err)
		assert.Equal(t, "editor/zsh.0", sessions.Directory["tmux-pane:editor/%0"].Name)
		assert.Equal(t, "editor/zsh.1", sessions.Directory["tmux-pane:editor/%1"].Name)
		assert.Equal(t, "tests/zsh", sessions.Directory["tmux-pane:tests/%2"].Name)
	})

	t.Run("should return empty sessions on error", func(t *testing.T) {
		mockTmux := new(tmux.MockTmux)
		mockTmux.On("ListTmuxPanes").Return(nil, fmt.Errorf("some error"))

		lister := NewLister(new(oswrap.MockOs), model.Config{}, new(home.MockHome), mockTmux, new(zoxide.MockZoxide), new(tmuxinator.MockTmuxinator))
		realLister, ok := lister.(*RealLister)
		if !ok {
			log.Fatal("Cannot convert lister to *RealLister")
		}

		sessions, err := listTmuxPanes(realLister)
		assert.Equal(t, model.SeshSessions{}, sessions)
		assert.NotNil(t, err)
		assert.Contains(t, err.Error(), "couldn't list tmux panes")
	})
}
