package startup

import (
	"fmt"
	"testing"

	"github.com/joshmedeski/sesh/v2/home"
	"github.com/joshmedeski/sesh/v2/lister"
	"github.com/joshmedeski/sesh/v2/model"
	"github.com/joshmedeski/sesh/v2/replacer"
	"github.com/joshmedeski/sesh/v2/tmux"
	"github.com/stretchr/testify/assert"
)

func TestExecTargetsSessionForWindowScriptsAndStartupCommand(t *testing.T) {
	mLister := lister.NewMockLister(t)
	mTmux := tmux.NewMockTmux(t)
	mHome := home.NewMockHome(t)
	mReplacer := replacer.NewMockReplacer(t)

	const sessionName = "myapp"
	const sessionPath = "/repo/myapp"
	mHome.EXPECT().ExpandPath(sessionPath).Return(sessionPath, nil)
	mTmux.EXPECT().NewWindowInSession(model.TmuxWindowOpts{
		Name:          "myapp-logs",
		StartDir:      sessionPath,
		TargetSession: sessionName,
	}).Return("myapp:1", nil)
	mTmux.EXPECT().SendKeys(sessionName+":myapp-logs", "echo ready").Return("", nil)
	mTmux.EXPECT().NextWindowInSession(sessionName).Return("", nil)
	mLister.EXPECT().FindConfigSession(sessionName).Return(model.SeshSession{}, false)
	mLister.EXPECT().FindConfigWildcard(sessionPath).Return(model.WildcardConfig{}, false)
	mReplacer.EXPECT().Replace("nvim", map[string]string{"{}": sessionPath}).Return("nvim")
	mTmux.EXPECT().SendKeys(sessionName+":", "nvim").Return("", nil)

	s := &RealStartup{
		lister:   mLister,
		tmux:     mTmux,
		config:   model.Config{WindowConfigs: []model.WindowConfig{{Name: "myapp-logs", StartupScript: "echo ready"}}, DefaultSessionConfig: model.DefaultSessionConfig{StartupCommand: "nvim"}},
		home:     mHome,
		replacer: mReplacer,
	}

	message, err := s.Exec(model.SeshSession{
		Name:        sessionName,
		Path:        sessionPath,
		WindowNames: []string{"myapp-logs"},
	})

	assert.NoError(t, err)
	assert.Equal(t, "executing startup command: nvim", message)
}

func TestCreatePanesSelectsFirstFocusedPane(t *testing.T) {
	mTmux := tmux.NewMockTmux(t)
	const target = "myapp:dev"
	window := model.WindowConfig{Name: "dev", Path: "/repo/myapp", Panes: []string{"logs", "server", "shell"}, Layout: "tiled"}
	panes := model.SeshPaneMap{
		lister.ConfigKey("logs"):   {Name: "logs", StartupCommand: "tail -f log"},
		lister.ConfigKey("server"): {Name: "server", StartupCommand: "make run", Focus: true},
		lister.ConfigKey("shell"):  {Name: "shell", Focus: true},
	}
	for i, name := range window.Panes {
		mTmux.EXPECT().NewPaneInWindow(model.TmuxPaneOpts{
			TargetWindow: target,
			StartDir:     window.Path,
			Command:      panes[lister.ConfigKey(name)].StartupCommand,
			Layout:       "tiled",
		}).Return(fmt.Sprintf("%%%d", i+1), nil).Once()
	}
	mTmux.EXPECT().SelectPaneByID("%2").Return("", nil).Once()

	s := &RealStartup{tmux: mTmux}
	_, err := s.createPanes(target, window, panes)
	assert.NoError(t, err)
}

func TestCreatePanesWithoutFocusKeepsMainPane(t *testing.T) {
	mTmux := tmux.NewMockTmux(t)
	window := model.WindowConfig{Name: "dev", Panes: []string{"logs"}}
	panes := model.SeshPaneMap{lister.ConfigKey("logs"): {Name: "logs"}}
	mTmux.EXPECT().NewPaneInWindow(model.TmuxPaneOpts{TargetWindow: "myapp:dev"}).Return("%1", nil).Once()

	s := &RealStartup{tmux: mTmux}
	_, err := s.createPanes("myapp:dev", window, panes)
	assert.NoError(t, err)
	mTmux.AssertNotCalled(t, "SelectPaneByID", "%1")
}

func TestCreatePanesZoomsFocusedPane(t *testing.T) {
	mTmux := tmux.NewMockTmux(t)
	window := model.WindowConfig{Name: "dev", Panes: []string{"logs"}, Zoom: true}
	panes := model.SeshPaneMap{lister.ConfigKey("logs"): {Name: "logs", Focus: true}}
	mTmux.EXPECT().NewPaneInWindow(model.TmuxPaneOpts{TargetWindow: "myapp:dev"}).Return("%1", nil).Once()
	mTmux.EXPECT().SelectPaneByID("%1").Return("", nil).Once()
	mTmux.EXPECT().ZoomPane("%1").Return("", nil).Once()

	s := &RealStartup{tmux: mTmux}
	_, err := s.createPanes("myapp:dev", window, panes)
	assert.NoError(t, err)
}

func TestCreatePanesZoomedPaneWinsOverFocus(t *testing.T) {
	mTmux := tmux.NewMockTmux(t)
	window := model.WindowConfig{Name: "dev", Panes: []string{"logs", "server"}}
	panes := model.SeshPaneMap{
		lister.ConfigKey("logs"):   {Name: "logs", Focus: true},
		lister.ConfigKey("server"): {Name: "server", Zoom: true},
	}
	mTmux.EXPECT().NewPaneInWindow(model.TmuxPaneOpts{TargetWindow: "myapp:dev"}).Return("%1", nil).Once()
	mTmux.EXPECT().NewPaneInWindow(model.TmuxPaneOpts{TargetWindow: "myapp:dev"}).Return("%2", nil).Once()
	mTmux.EXPECT().SelectPaneByID("%1").Return("", nil).Once()
	mTmux.EXPECT().ZoomPane("%2").Return("", nil).Once()

	s := &RealStartup{tmux: mTmux}
	_, err := s.createPanes("myapp:dev", window, panes)
	assert.NoError(t, err)
}

func TestCreatePanesZoomsMainPaneWithoutFocus(t *testing.T) {
	mTmux := tmux.NewMockTmux(t)
	window := model.WindowConfig{Name: "dev", Panes: []string{"logs"}, Zoom: true}
	panes := model.SeshPaneMap{lister.ConfigKey("logs"): {Name: "logs"}}
	mTmux.EXPECT().NewPaneInWindow(model.TmuxPaneOpts{TargetWindow: "myapp:dev"}).Return("%1", nil).Once()
	mTmux.EXPECT().ZoomPane("myapp:dev").Return("", nil).Once()

	s := &RealStartup{tmux: mTmux}
	_, err := s.createPanes("myapp:dev", window, panes)
	assert.NoError(t, err)
}
