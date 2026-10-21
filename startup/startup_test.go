package startup

import (
	"errors"
	"fmt"
	"testing"

	"github.com/joshmedeski/sesh/v2/home"
	"github.com/joshmedeski/sesh/v2/lister"
	"github.com/joshmedeski/sesh/v2/model"
	"github.com/joshmedeski/sesh/v2/replacer"
	"github.com/joshmedeski/sesh/v2/tmux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
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
	mTmux.EXPECT().SendKeys("myapp:1", "echo ready").Return("", nil)
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

func TestExecCreatesConfiguredPanesWithWindowLayout(t *testing.T) {
	mLister := lister.NewMockLister(t)
	mTmux := tmux.NewMockTmux(t)
	mHome := home.NewMockHome(t)

	const sessionName = "myapp"
	const sessionPath = "/repo/myapp"
	mHome.EXPECT().ExpandPath(sessionPath).Return(sessionPath, nil)
	mTmux.EXPECT().NewWindowInSession(model.TmuxWindowOpts{
		Name: "dev", StartDir: sessionPath, TargetSession: sessionName,
	}).Return("myapp:1", nil)
	mTmux.EXPECT().SendKeys("myapp:1", "nvim").Return("", nil)
	mTmux.EXPECT().NewPaneInWindow(model.TmuxPaneOpts{
		TargetWindow: "myapp:1", StartDir: sessionPath, Command: "lazydocker", Layout: "main-vertical",
	}).Return("%1", nil).Once()
	mTmux.EXPECT().NextWindowInSession(sessionName).Return("", nil)
	mLister.EXPECT().FindConfigSession(sessionName).Return(model.SeshSession{}, false)
	mLister.EXPECT().FindConfigWildcard(sessionPath).Return(model.WildcardConfig{}, false)

	s := &RealStartup{
		lister: mLister,
		tmux:   mTmux,
		config: model.Config{
			WindowConfigs: []model.WindowConfig{{Name: "dev", StartupScript: "nvim", Panes: []string{"docker"}, Layout: "main-vertical"}},
			PaneConfigs:   []model.PaneConfig{{Name: "docker", StartupCommand: "lazydocker"}},
		},
		home: mHome,
	}

	_, err := s.Exec(model.SeshSession{Name: sessionName, Path: sessionPath, WindowNames: []string{"dev"}})
	assert.NoError(t, err)
}

func TestCreatePanesReportsUndefinedPaneAndCreatesTheRest(t *testing.T) {
	mTmux := tmux.NewMockTmux(t)
	window := model.WindowConfig{Name: "dev", Panes: []string{"missing", "logs", "server"}}
	panes := model.SeshPaneMap{
		lister.ConfigKey("logs"):   {Name: "logs", Focus: true},
		lister.ConfigKey("server"): {Name: "server"},
	}
	mTmux.EXPECT().NewPaneInWindow(model.TmuxPaneOpts{TargetWindow: "myapp:dev"}).Return("%1", nil).Once()
	mTmux.EXPECT().NewPaneInWindow(model.TmuxPaneOpts{TargetWindow: "myapp:dev", TargetPane: "%1"}).Return("%2", nil).Once()
	mTmux.EXPECT().SelectPaneByID("%1").Return("", nil).Once()

	s := &RealStartup{tmux: mTmux}
	err := s.createPanes("myapp:dev", window, panes)
	assert.EqualError(t, err, "window dev: pane missing is not defined in config")
}

func TestCreatePanesKeepsPaneWhenSendKeysOrLayoutFails(t *testing.T) {
	mTmux := tmux.NewMockTmux(t)
	window := model.WindowConfig{Name: "dev", Panes: []string{"logs", "server"}, Zoom: true}
	panes := model.SeshPaneMap{
		lister.ConfigKey("logs"):   {Name: "logs", StartupCommand: "tail", Focus: true},
		lister.ConfigKey("server"): {Name: "server"},
	}
	mTmux.EXPECT().NewPaneInWindow(model.TmuxPaneOpts{TargetWindow: "myapp:dev", Command: "tail"}).
		Return("%1", errors.New("send-keys failed")).Once()
	mTmux.EXPECT().NewPaneInWindow(model.TmuxPaneOpts{TargetWindow: "myapp:dev", TargetPane: "%1"}).Return("%2", nil).Once()
	mTmux.EXPECT().SelectPaneByID("%1").Return("", nil).Once()
	mTmux.EXPECT().ZoomPane("%1").Return("", nil).Once()

	s := &RealStartup{tmux: mTmux}
	err := s.createPanes("myapp:dev", window, panes)
	assert.EqualError(t, err, "window dev: pane logs: send-keys failed")
}

func TestCreatePanesSkipsPaneWithoutID(t *testing.T) {
	mTmux := tmux.NewMockTmux(t)
	window := model.WindowConfig{Name: "dev", Panes: []string{"logs", "server"}}
	panes := model.SeshPaneMap{
		lister.ConfigKey("logs"):   {Name: "logs", Focus: true},
		lister.ConfigKey("server"): {Name: "server"},
	}
	mTmux.EXPECT().NewPaneInWindow(model.TmuxPaneOpts{TargetWindow: "myapp:dev"}).Return("", errors.New("split failed")).Once()
	mTmux.EXPECT().NewPaneInWindow(model.TmuxPaneOpts{TargetWindow: "myapp:dev"}).Return("%2", nil).Once()

	s := &RealStartup{tmux: mTmux}
	err := s.createPanes("myapp:dev", window, panes)
	assert.EqualError(t, err, "window dev: pane logs: split failed")
	mTmux.AssertNotCalled(t, "SelectPaneByID", mock.Anything)
}

func TestCreatePanesRecordsFocusAndZoomFailures(t *testing.T) {
	mTmux := tmux.NewMockTmux(t)
	window := model.WindowConfig{Name: "dev", Panes: []string{"logs"}}
	panes := model.SeshPaneMap{lister.ConfigKey("logs"): {Name: "logs", Focus: true, Zoom: true}}
	mTmux.EXPECT().NewPaneInWindow(model.TmuxPaneOpts{TargetWindow: "myapp:dev"}).Return("%1", nil).Once()
	mTmux.EXPECT().SelectPaneByID("%1").Return("", errors.New("no pane")).Once()
	mTmux.EXPECT().ZoomPane("%1").Return("", errors.New("cannot zoom")).Once()

	s := &RealStartup{tmux: mTmux}
	err := s.createPanes("myapp:dev", window, panes)
	assert.EqualError(t, err, "window dev: focus pane: no pane\nwindow dev: zoom pane: cannot zoom")
}

func TestExecKeepsGoingAfterUndefinedWindow(t *testing.T) {
	mLister := lister.NewMockLister(t)
	mTmux := tmux.NewMockTmux(t)
	mHome := home.NewMockHome(t)
	mReplacer := replacer.NewMockReplacer(t)

	const sessionName = "myapp"
	const sessionPath = "/repo/myapp"
	mHome.EXPECT().ExpandPath(sessionPath).Return(sessionPath, nil)
	mTmux.EXPECT().NewWindowInSession(model.TmuxWindowOpts{
		Name: "logs", StartDir: sessionPath, TargetSession: sessionName,
	}).Return("myapp:1", nil)
	mTmux.EXPECT().SendKeys("myapp:1", "").Return("", nil)
	mTmux.EXPECT().NextWindowInSession(sessionName).Return("", nil)
	mLister.EXPECT().FindConfigSession(sessionName).Return(model.SeshSession{}, false)
	mLister.EXPECT().FindConfigWildcard(sessionPath).Return(model.WildcardConfig{}, false)
	mReplacer.EXPECT().Replace("nvim", map[string]string{"{}": sessionPath}).Return("nvim")
	mTmux.EXPECT().SendKeys(sessionName+":", "nvim").Return("", nil)

	s := &RealStartup{
		lister: mLister, tmux: mTmux, home: mHome, replacer: mReplacer,
		config: model.Config{
			WindowConfigs:        []model.WindowConfig{{Name: "logs"}},
			DefaultSessionConfig: model.DefaultSessionConfig{StartupCommand: "nvim"},
		},
	}

	message, err := s.Exec(model.SeshSession{
		Name: sessionName, Path: sessionPath, WindowNames: []string{"missing", "logs"},
	})

	assert.EqualError(t, err, "window missing is not defined in config")
	assert.Equal(t, "executing startup command: nvim", message)
}

func TestExecKeepsGoingAfterWindowCreationFails(t *testing.T) {
	mLister := lister.NewMockLister(t)
	mTmux := tmux.NewMockTmux(t)
	mHome := home.NewMockHome(t)

	const sessionName = "myapp"
	const sessionPath = "/repo/myapp"
	mHome.EXPECT().ExpandPath(sessionPath).Return(sessionPath, nil).Twice()
	mTmux.EXPECT().NewWindowInSession(model.TmuxWindowOpts{
		Name: "bad", StartDir: sessionPath, TargetSession: sessionName,
	}).Return("", errors.New("duplicate window"))
	mTmux.EXPECT().NewWindowInSession(model.TmuxWindowOpts{
		Name: "dev", StartDir: sessionPath, TargetSession: sessionName,
	}).Return("myapp:2", nil)
	mTmux.EXPECT().SendKeys("myapp:2", "nvim").Return("", errors.New("send-keys failed"))
	mTmux.EXPECT().NewPaneInWindow(model.TmuxPaneOpts{TargetWindow: "myapp:2", StartDir: sessionPath}).Return("%1", nil)
	mTmux.EXPECT().NextWindowInSession(sessionName).Return("", nil)
	mLister.EXPECT().FindConfigSession(sessionName).Return(model.SeshSession{}, false)
	mLister.EXPECT().FindConfigWildcard(sessionPath).Return(model.WildcardConfig{}, false)

	s := &RealStartup{
		lister: mLister, tmux: mTmux, home: mHome,
		config: model.Config{
			WindowConfigs: []model.WindowConfig{{Name: "bad"}, {Name: "dev", StartupScript: "nvim", Panes: []string{"p", "nope"}}},
			PaneConfigs:   []model.PaneConfig{{Name: "p"}},
		},
	}

	message, err := s.Exec(model.SeshSession{
		Name: sessionName, Path: sessionPath, WindowNames: []string{"bad", "dev"},
	})

	assert.Empty(t, message)
	assert.EqualError(t, err, "window bad: duplicate window\n"+
		"window dev: send-keys failed\n"+
		"window dev: pane nope is not defined in config")
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
		prevID := ""
		if i > 0 {
			prevID = fmt.Sprintf("%%%d", i)
		}
		mTmux.EXPECT().NewPaneInWindow(model.TmuxPaneOpts{
			TargetWindow: target,
			TargetPane:   prevID,
			StartDir:     window.Path,
			Command:      panes[lister.ConfigKey(name)].StartupCommand,
			Layout:       "tiled",
		}).Return(fmt.Sprintf("%%%d", i+1), nil).Once()
	}
	mTmux.EXPECT().SelectPaneByID("%2").Return("", nil).Once()

	s := &RealStartup{tmux: mTmux}
	err := s.createPanes(target, window, panes)
	assert.NoError(t, err)
}

func TestCreatePanesWithoutFocusKeepsMainPane(t *testing.T) {
	mTmux := tmux.NewMockTmux(t)
	window := model.WindowConfig{Name: "dev", Panes: []string{"logs"}}
	panes := model.SeshPaneMap{lister.ConfigKey("logs"): {Name: "logs"}}
	mTmux.EXPECT().NewPaneInWindow(model.TmuxPaneOpts{TargetWindow: "myapp:dev"}).Return("%1", nil).Once()

	s := &RealStartup{tmux: mTmux}
	err := s.createPanes("myapp:dev", window, panes)
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
	err := s.createPanes("myapp:dev", window, panes)
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
	mTmux.EXPECT().NewPaneInWindow(model.TmuxPaneOpts{TargetWindow: "myapp:dev", TargetPane: "%1"}).Return("%2", nil).Once()
	mTmux.EXPECT().SelectPaneByID("%1").Return("", nil).Once()
	mTmux.EXPECT().ZoomPane("%2").Return("", nil).Once()

	s := &RealStartup{tmux: mTmux}
	err := s.createPanes("myapp:dev", window, panes)
	assert.NoError(t, err)
}

func TestCreatePanesZoomsMainPaneWithoutFocus(t *testing.T) {
	mTmux := tmux.NewMockTmux(t)
	window := model.WindowConfig{Name: "dev", Panes: []string{"logs"}, Zoom: true}
	panes := model.SeshPaneMap{lister.ConfigKey("logs"): {Name: "logs"}}
	mTmux.EXPECT().NewPaneInWindow(model.TmuxPaneOpts{TargetWindow: "myapp:dev"}).Return("%1", nil).Once()
	mTmux.EXPECT().ZoomPane("myapp:dev").Return("", nil).Once()

	s := &RealStartup{tmux: mTmux}
	err := s.createPanes("myapp:dev", window, panes)
	assert.NoError(t, err)
}
