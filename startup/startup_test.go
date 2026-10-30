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
		TargetWindow: "myapp:1", StartDir: sessionPath,
	}).Return("%1", nil).Once()
	mTmux.EXPECT().SendKeys("%1", "lazydocker").Return("", nil).Once()
	mTmux.EXPECT().SelectLayout("myapp:1", "main-vertical").Return("", nil).Once()
	mTmux.EXPECT().NextWindowInSession(sessionName).Return("", nil)
	mLister.EXPECT().FindConfigSession(sessionName).Return(model.SeshSession{}, false)
	mLister.EXPECT().FindConfigWildcard(sessionPath).Return(model.WildcardConfig{}, false)

	s := &RealStartup{
		lister: mLister,
		tmux:   mTmux,
		config: model.Config{
			WindowConfigs: []model.WindowConfig{{Name: "dev", StartupScript: "nvim", WindowPanes: model.WindowPanes{Panes: []string{"docker"}, Layout: "main-vertical"}}},
			PaneConfigs:   []model.PaneConfig{{Name: "docker", StartupCommand: "lazydocker"}},
		},
		home: mHome,
	}

	_, err := s.Exec(model.SeshSession{Name: sessionName, Path: sessionPath, WindowNames: []string{"dev"}})
	assert.NoError(t, err)
}

func TestCreatePanesStopsAtUndefinedPane(t *testing.T) {
	mTmux := tmux.NewMockTmux(t)
	windowPanes := model.WindowPanes{Panes: []string{"logs", "missing", "server"}}
	panes := model.SeshPaneMap{
		lister.ConfigKey("logs"):   {Name: "logs", Focus: true},
		lister.ConfigKey("server"): {Name: "server"},
	}
	mTmux.EXPECT().NewPaneInWindow(model.TmuxPaneOpts{TargetWindow: "myapp:dev"}).Return("%1", nil).Once()

	s := &RealStartup{tmux: mTmux}
	err := s.createPanes("myapp:dev", "", windowPanes, panes)
	assert.EqualError(t, err, "pane missing is not defined in config")
	mTmux.AssertNotCalled(t, "SelectPaneByID", mock.Anything)
}

func TestCreatePanesStopsWhenSplitFails(t *testing.T) {
	mTmux := tmux.NewMockTmux(t)
	windowPanes := model.WindowPanes{Panes: []string{"logs", "server"}}
	panes := model.SeshPaneMap{
		lister.ConfigKey("logs"):   {Name: "logs"},
		lister.ConfigKey("server"): {Name: "server"},
	}
	mTmux.EXPECT().NewPaneInWindow(model.TmuxPaneOpts{TargetWindow: "myapp:dev"}).Return("", errors.New("split failed")).Once()

	s := &RealStartup{tmux: mTmux}
	err := s.createPanes("myapp:dev", "", windowPanes, panes)
	assert.EqualError(t, err, "pane logs: split failed")
}

func TestCreatePanesReportsFocusFailure(t *testing.T) {
	mTmux := tmux.NewMockTmux(t)
	windowPanes := model.WindowPanes{Panes: []string{"logs"}}
	panes := model.SeshPaneMap{lister.ConfigKey("logs"): {Name: "logs", Focus: true, Zoom: true}}
	mTmux.EXPECT().NewPaneInWindow(model.TmuxPaneOpts{TargetWindow: "myapp:dev"}).Return("%1", nil).Once()
	mTmux.EXPECT().SelectPaneByID("%1").Return("", errors.New("no pane")).Once()

	s := &RealStartup{tmux: mTmux}
	err := s.createPanes("myapp:dev", "", windowPanes, panes)
	assert.EqualError(t, err, "focus pane: no pane")
	mTmux.AssertNotCalled(t, "ZoomPane", mock.Anything)
}

func TestCreatePanesReportsZoomFailure(t *testing.T) {
	mTmux := tmux.NewMockTmux(t)
	windowPanes := model.WindowPanes{Panes: []string{"logs"}}
	panes := model.SeshPaneMap{lister.ConfigKey("logs"): {Name: "logs", Zoom: true}}
	mTmux.EXPECT().NewPaneInWindow(model.TmuxPaneOpts{TargetWindow: "myapp:dev"}).Return("%1", nil).Once()
	mTmux.EXPECT().ZoomPane("%1").Return("", errors.New("cannot zoom")).Once()

	s := &RealStartup{tmux: mTmux}
	err := s.createPanes("myapp:dev", "", windowPanes, panes)
	assert.EqualError(t, err, "zoom pane: cannot zoom")
}

func TestExecStopsAtUndefinedWindow(t *testing.T) {
	mTmux := tmux.NewMockTmux(t)
	s := &RealStartup{
		tmux:   mTmux,
		config: model.Config{WindowConfigs: []model.WindowConfig{{Name: "logs"}}},
	}

	message, err := s.Exec(model.SeshSession{
		Name: "myapp", Path: "/repo/myapp", WindowNames: []string{"missing", "logs"},
	})

	assert.EqualError(t, err, "window missing is not defined in config")
	assert.Empty(t, message)
	mTmux.AssertNotCalled(t, "NewWindowInSession", mock.Anything)
	mTmux.AssertNotCalled(t, "NextWindowInSession", mock.Anything)
	mTmux.AssertNotCalled(t, "SendKeys", mock.Anything, mock.Anything)
}

func TestExecPrefixesPaneErrorsWithWindow(t *testing.T) {
	mTmux := tmux.NewMockTmux(t)
	mHome := home.NewMockHome(t)
	mHome.EXPECT().ExpandPath("/repo/myapp").Return("/repo/myapp", nil)
	mTmux.EXPECT().NewWindowInSession(model.TmuxWindowOpts{
		Name: "dev", StartDir: "/repo/myapp", TargetSession: "myapp",
	}).Return("myapp:1", nil)
	mTmux.EXPECT().SendKeys("myapp:1", "").Return("", nil)

	s := &RealStartup{
		tmux: mTmux, home: mHome,
		config: model.Config{WindowConfigs: []model.WindowConfig{{Name: "dev", WindowPanes: model.WindowPanes{Panes: []string{"foo"}}}}},
	}

	_, err := s.Exec(model.SeshSession{Name: "myapp", Path: "/repo/myapp", WindowNames: []string{"dev"}})

	assert.EqualError(t, err, "window dev: pane foo is not defined in config")
}

func TestExecReportsStartupCommandFailure(t *testing.T) {
	mLister := lister.NewMockLister(t)
	mTmux := tmux.NewMockTmux(t)
	mReplacer := replacer.NewMockReplacer(t)
	mTmux.EXPECT().NextWindowInSession("myapp").Return("", nil)
	mLister.EXPECT().FindConfigSession("myapp").Return(model.SeshSession{}, false)
	mLister.EXPECT().FindConfigWildcard("/repo/myapp").Return(model.WildcardConfig{}, false)
	mReplacer.EXPECT().Replace("nvim", map[string]string{"{}": "/repo/myapp"}).Return("nvim")
	mTmux.EXPECT().SendKeys("myapp:", "nvim").Return("", errors.New("no server"))

	s := &RealStartup{
		lister: mLister, tmux: mTmux, replacer: mReplacer,
		config: model.Config{DefaultSessionConfig: model.DefaultSessionConfig{StartupCommand: "nvim"}},
	}

	message, err := s.Exec(model.SeshSession{Name: "myapp", Path: "/repo/myapp"})

	assert.EqualError(t, err, "failed to send startup command: no server")
	assert.Empty(t, message)
}

func TestCreatePanesSelectsFirstFocusedPane(t *testing.T) {
	mTmux := tmux.NewMockTmux(t)
	const target = "myapp:dev"
	windowPanes := model.WindowPanes{Panes: []string{"logs", "server", "shell"}, Layout: "tiled"}
	panes := model.SeshPaneMap{
		lister.ConfigKey("logs"):   {Name: "logs", StartupCommand: "tail -f log"},
		lister.ConfigKey("server"): {Name: "server", StartupCommand: "make run", Focus: true},
		lister.ConfigKey("shell"):  {Name: "shell", Focus: true},
	}
	for i, name := range windowPanes.Panes {
		prevID := ""
		if i > 0 {
			prevID = fmt.Sprintf("%%%d", i)
		}
		paneID := fmt.Sprintf("%%%d", i+1)
		mTmux.EXPECT().NewPaneInWindow(model.TmuxPaneOpts{
			TargetWindow: target,
			TargetPane:   prevID,
			StartDir:     "/repo/myapp",
		}).Return(paneID, nil).Once()
		if command := panes[lister.ConfigKey(name)].StartupCommand; command != "" {
			mTmux.EXPECT().SendKeys(paneID, command).Return("", nil).Once()
		}
	}
	mTmux.EXPECT().SelectLayout(target, "tiled").Return("", nil).Once()
	mTmux.EXPECT().SelectPaneByID("%2").Return("", nil).Once()

	s := &RealStartup{tmux: mTmux}
	err := s.createPanes(target, "/repo/myapp", windowPanes, panes)
	assert.NoError(t, err)
}

func TestCreatePanesAppliesCustomLayoutAfterAllPanes(t *testing.T) {
	mTmux := tmux.NewMockTmux(t)
	const target = "myapp:dev"
	const layout = "5d4b,80x24,0,0{40x24,0,0,0,39x24,41,0[39x12,41,0,1,39x11,41,13,2]}"
	windowPanes := model.WindowPanes{Panes: []string{"logs", "server"}, Layout: layout}
	panes := model.SeshPaneMap{
		lister.ConfigKey("logs"):   {Name: "logs"},
		lister.ConfigKey("server"): {Name: "server"},
	}
	var calls []string
	mTmux.EXPECT().NewPaneInWindow(model.TmuxPaneOpts{TargetWindow: target}).
		Run(func(model.TmuxPaneOpts) { calls = append(calls, "split") }).Return("%1", nil).Once()
	mTmux.EXPECT().NewPaneInWindow(model.TmuxPaneOpts{TargetWindow: target, TargetPane: "%1"}).
		Run(func(model.TmuxPaneOpts) { calls = append(calls, "split") }).Return("%2", nil).Once()
	mTmux.EXPECT().SelectLayout(target, layout).
		Run(func(string, string) { calls = append(calls, "layout") }).Return("", nil).Once()

	s := &RealStartup{tmux: mTmux}
	err := s.createPanes(target, "", windowPanes, panes)
	assert.NoError(t, err)
	assert.Equal(t, []string{"split", "split", "layout"}, calls)
}

func TestCreatePanesReportsLayoutFailure(t *testing.T) {
	mTmux := tmux.NewMockTmux(t)
	windowPanes := model.WindowPanes{Panes: []string{"logs"}, Layout: "bogus", Zoom: true}
	panes := model.SeshPaneMap{lister.ConfigKey("logs"): {Name: "logs", Focus: true}}
	mTmux.EXPECT().NewPaneInWindow(model.TmuxPaneOpts{TargetWindow: "myapp:dev"}).Return("%1", nil).Once()
	mTmux.EXPECT().SelectLayout("myapp:dev", "bogus").Return("", errors.New("invalid layout")).Once()

	s := &RealStartup{tmux: mTmux}
	err := s.createPanes("myapp:dev", "", windowPanes, panes)
	assert.EqualError(t, err, "apply layout bogus: invalid layout")
	mTmux.AssertNotCalled(t, "SelectPaneByID", mock.Anything)
	mTmux.AssertNotCalled(t, "ZoomPane", mock.Anything)
}

func TestCreatePanesReportsPaneCommandFailure(t *testing.T) {
	mTmux := tmux.NewMockTmux(t)
	windowPanes := model.WindowPanes{Panes: []string{"logs", "server"}}
	panes := model.SeshPaneMap{
		lister.ConfigKey("logs"):   {Name: "logs", StartupCommand: "tail -f log"},
		lister.ConfigKey("server"): {Name: "server"},
	}
	mTmux.EXPECT().NewPaneInWindow(model.TmuxPaneOpts{TargetWindow: "myapp:dev"}).Return("%1", nil).Once()
	mTmux.EXPECT().SendKeys("%1", "tail -f log").Return("", errors.New("no pane")).Once()

	s := &RealStartup{tmux: mTmux}
	err := s.createPanes("myapp:dev", "", windowPanes, panes)
	assert.EqualError(t, err, "pane logs: run startup command: no pane")
}

func TestCreatePanesWithoutFocusKeepsMainPane(t *testing.T) {
	mTmux := tmux.NewMockTmux(t)
	windowPanes := model.WindowPanes{Panes: []string{"logs"}}
	panes := model.SeshPaneMap{lister.ConfigKey("logs"): {Name: "logs"}}
	mTmux.EXPECT().NewPaneInWindow(model.TmuxPaneOpts{TargetWindow: "myapp:dev"}).Return("%1", nil).Once()

	s := &RealStartup{tmux: mTmux}
	err := s.createPanes("myapp:dev", "", windowPanes, panes)
	assert.NoError(t, err)
	mTmux.AssertNotCalled(t, "SelectPaneByID", "%1")
}

func TestCreatePanesZoomsFocusedPane(t *testing.T) {
	mTmux := tmux.NewMockTmux(t)
	windowPanes := model.WindowPanes{Panes: []string{"logs"}, Zoom: true}
	panes := model.SeshPaneMap{lister.ConfigKey("logs"): {Name: "logs", Focus: true}}
	mTmux.EXPECT().NewPaneInWindow(model.TmuxPaneOpts{TargetWindow: "myapp:dev"}).Return("%1", nil).Once()
	mTmux.EXPECT().SelectPaneByID("%1").Return("", nil).Once()
	mTmux.EXPECT().ZoomPane("%1").Return("", nil).Once()

	s := &RealStartup{tmux: mTmux}
	err := s.createPanes("myapp:dev", "", windowPanes, panes)
	assert.NoError(t, err)
}

func TestCreatePanesZoomedPaneWinsOverFocus(t *testing.T) {
	mTmux := tmux.NewMockTmux(t)
	windowPanes := model.WindowPanes{Panes: []string{"logs", "server"}}
	panes := model.SeshPaneMap{
		lister.ConfigKey("logs"):   {Name: "logs", Focus: true},
		lister.ConfigKey("server"): {Name: "server", Zoom: true},
	}
	mTmux.EXPECT().NewPaneInWindow(model.TmuxPaneOpts{TargetWindow: "myapp:dev"}).Return("%1", nil).Once()
	mTmux.EXPECT().NewPaneInWindow(model.TmuxPaneOpts{TargetWindow: "myapp:dev", TargetPane: "%1"}).Return("%2", nil).Once()
	mTmux.EXPECT().SelectPaneByID("%1").Return("", nil).Once()
	mTmux.EXPECT().ZoomPane("%2").Return("", nil).Once()

	s := &RealStartup{tmux: mTmux}
	err := s.createPanes("myapp:dev", "", windowPanes, panes)
	assert.NoError(t, err)
}

func TestCreatePanesZoomsMainPaneWithoutFocus(t *testing.T) {
	mTmux := tmux.NewMockTmux(t)
	windowPanes := model.WindowPanes{Panes: []string{"logs"}, Zoom: true}
	panes := model.SeshPaneMap{lister.ConfigKey("logs"): {Name: "logs"}}
	mTmux.EXPECT().NewPaneInWindow(model.TmuxPaneOpts{TargetWindow: "myapp:dev"}).Return("%1", nil).Once()
	mTmux.EXPECT().ZoomPane("myapp:dev").Return("", nil).Once()

	s := &RealStartup{tmux: mTmux}
	err := s.createPanes("myapp:dev", "", windowPanes, panes)
	assert.NoError(t, err)
}

func TestExecCreatesFirstWindowPanesAfterStartupCommand(t *testing.T) {
	mLister := lister.NewMockLister(t)
	mTmux := tmux.NewMockTmux(t)
	mHome := home.NewMockHome(t)
	mReplacer := replacer.NewMockReplacer(t)

	const sessionPath = "/repo/myapp"
	var calls []string
	mTmux.EXPECT().NextWindowInSession("myapp").Return("", nil)
	mLister.EXPECT().FindConfigSession("myapp").Return(model.SeshSession{StartupCommand: "nvim"}, true)
	mReplacer.EXPECT().Replace("nvim", map[string]string{"{}": sessionPath}).Return("nvim")
	mTmux.EXPECT().SendKeys("myapp:", "nvim").
		Run(func(string, string) { calls = append(calls, "startup command") }).Return("", nil)
	mHome.EXPECT().ExpandPath(sessionPath).Return(sessionPath, nil)
	mTmux.EXPECT().NewPaneInWindow(model.TmuxPaneOpts{
		TargetWindow: "myapp:^", StartDir: sessionPath,
	}).Run(func(model.TmuxPaneOpts) { calls = append(calls, "pane") }).Return("%1", nil)
	mTmux.EXPECT().SendKeys("%1", "lazydocker").Return("", nil)
	mTmux.EXPECT().SelectLayout("myapp:^", "tiled").Return("", nil)

	s := &RealStartup{
		lister: mLister, tmux: mTmux, home: mHome, replacer: mReplacer,
		config: model.Config{
			SessionConfigs: []model.SessionConfig{{
				Name:                 "myapp",
				DefaultSessionConfig: model.DefaultSessionConfig{WindowPanes: model.WindowPanes{Panes: []string{"docker"}, Layout: "tiled"}},
			}},
			PaneConfigs: []model.PaneConfig{{Name: "docker", StartupCommand: "lazydocker"}},
		},
	}

	message, err := s.Exec(model.SeshSession{Name: "myapp", Path: sessionPath})

	assert.NoError(t, err)
	assert.Equal(t, "executing startup command: nvim", message)
	assert.Equal(t, []string{"startup command", "pane"}, calls)
}

func TestExecFirstWindowPanePrecedence(t *testing.T) {
	const sessionPath = "/repo/myapp"
	sessionCfg := model.SessionConfig{
		Name:                 "myapp",
		DefaultSessionConfig: model.DefaultSessionConfig{WindowPanes: model.WindowPanes{Panes: []string{"a"}, Layout: "tiled", Zoom: true}},
	}
	wildcardCfg := model.WildcardConfig{WindowPanes: model.WindowPanes{Panes: []string{"b"}, Layout: "even-vertical"}}
	defaultCfg := model.DefaultSessionConfig{WindowPanes: model.WindowPanes{Panes: []string{"c"}, Layout: "main-vertical"}}
	noPanes := model.SessionConfig{Name: "myapp"}

	tests := []struct {
		name     string
		sessions []model.SessionConfig
		wildcard model.WildcardConfig
		wildOK   bool
		pane     string
		layout   string
		zoom     bool
	}{
		{"session wins", []model.SessionConfig{sessionCfg}, wildcardCfg, true, "a", "tiled", true},
		{"wildcard beats default", []model.SessionConfig{noPanes}, wildcardCfg, true, "b", "even-vertical", false},
		{"default as fallback", nil, model.WildcardConfig{}, true, "c", "main-vertical", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mLister := lister.NewMockLister(t)
			mTmux := tmux.NewMockTmux(t)
			mHome := home.NewMockHome(t)
			mLister.EXPECT().FindConfigSession("myapp").Return(model.SeshSession{}, false).Maybe()
			mLister.EXPECT().FindConfigWildcard(sessionPath).Return(tt.wildcard, tt.wildOK).Maybe()
			mTmux.EXPECT().NextWindowInSession("myapp").Return("", nil)
			mHome.EXPECT().ExpandPath(sessionPath).Return(sessionPath, nil)
			mTmux.EXPECT().NewPaneInWindow(model.TmuxPaneOpts{
				TargetWindow: "myapp:^", StartDir: sessionPath,
			}).Return("%1", nil).Once()
			mTmux.EXPECT().SelectLayout("myapp:^", tt.layout).Return("", nil).Once()
			if tt.zoom {
				mTmux.EXPECT().ZoomPane("myapp:^").Return("", nil).Once()
			}

			s := &RealStartup{
				lister: mLister, tmux: mTmux, home: mHome,
				config: model.Config{
					SessionConfigs:       tt.sessions,
					DefaultSessionConfig: defaultCfg,
					PaneConfigs:          []model.PaneConfig{{Name: "a"}, {Name: "b"}, {Name: "c"}},
				},
			}

			_, err := s.Exec(model.SeshSession{Name: "myapp", Path: sessionPath})
			assert.NoError(t, err)
		})
	}
}

func TestExecFirstWindowUndefinedPane(t *testing.T) {
	mLister := lister.NewMockLister(t)
	mTmux := tmux.NewMockTmux(t)
	mHome := home.NewMockHome(t)
	mLister.EXPECT().FindConfigSession("myapp").Return(model.SeshSession{}, false)
	mLister.EXPECT().FindConfigWildcard("/repo/myapp").Return(model.WildcardConfig{}, false)
	mTmux.EXPECT().NextWindowInSession("myapp").Return("", nil)
	mHome.EXPECT().ExpandPath("/repo/myapp").Return("/repo/myapp", nil)

	s := &RealStartup{
		lister: mLister, tmux: mTmux, home: mHome,
		config: model.Config{DefaultSessionConfig: model.DefaultSessionConfig{WindowPanes: model.WindowPanes{Panes: []string{"x"}}}},
	}

	_, err := s.Exec(model.SeshSession{Name: "myapp", Path: "/repo/myapp"})

	assert.EqualError(t, err, "first window: pane x is not defined in config")
}

func TestExecWithoutFirstWindowPanesCreatesNoPanes(t *testing.T) {
	mLister := lister.NewMockLister(t)
	mTmux := tmux.NewMockTmux(t)
	mLister.EXPECT().FindConfigSession("myapp").Return(model.SeshSession{}, false)
	mLister.EXPECT().FindConfigWildcard("/repo/myapp").Return(model.WildcardConfig{}, false)
	mTmux.EXPECT().NextWindowInSession("myapp").Return("", nil)

	s := &RealStartup{lister: mLister, tmux: mTmux}

	message, err := s.Exec(model.SeshSession{Name: "myapp", Path: "/repo/myapp"})

	assert.NoError(t, err)
	assert.Empty(t, message)
	mTmux.AssertNotCalled(t, "NewPaneInWindow", mock.Anything)
}
