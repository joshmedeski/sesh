package connector

import (
	"testing"

	"github.com/joshmedeski/sesh/v2/dir"
	"github.com/joshmedeski/sesh/v2/focuser"
	"github.com/joshmedeski/sesh/v2/home"
	"github.com/joshmedeski/sesh/v2/lister"
	"github.com/joshmedeski/sesh/v2/model"
	"github.com/joshmedeski/sesh/v2/tmux"
	"github.com/joshmedeski/sesh/v2/zoxide"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// windowConnector wires the mocks every window connect test needs, with the
// target session already attached so the tests only have to describe the
// window behaviour they care about.
func windowConnector(t *testing.T) (*RealConnector, *tmux.MockTmux, *lister.MockLister, *dir.MockDir, *home.MockHome) {
	mTmux := tmux.NewMockTmux(t)
	mLister := lister.NewMockLister(t)
	mDir := dir.NewMockDir(t)
	mHome := home.NewMockHome(t)

	c := NewConnector(
		model.Config{}, mDir, mHome, mLister, nil, nil, mTmux, nil, nil, nil,
	).(*RealConnector)

	return c, mTmux, mLister, mDir, mHome
}

func attachedTo(mLister *lister.MockLister, name string, path string) {
	mLister.EXPECT().GetAttachedTmuxSession().Return(model.SeshSession{Name: name, Path: path}, true)
}

// notADirectory keeps a plain window name from being read as a path.
func notADirectory(mHome *home.MockHome, mDir *dir.MockDir, name string) {
	mHome.EXPECT().ExpandPath(name).Return(name, nil)
	mDir.EXPECT().Dir(name).Return(false, "")
}

func TestConnectWindow(t *testing.T) {
	t.Run("sends the command into a window that already exists", func(t *testing.T) {
		c, mTmux, mLister, mDir, mHome := windowConnector(t)
		attachedTo(mLister, "second-brain", "/home/dev/brain")
		notADirectory(mHome, mDir, "notes")
		mTmux.EXPECT().ListWindows("second-brain").Return([]*model.TmuxWindow{
			{Name: "shell", Index: 1},
			{Name: "notes", Index: 4},
		}, nil)
		mTmux.EXPECT().SendKeys("second-brain:4", "claude").Return("", nil)
		mTmux.EXPECT().SelectWindow("second-brain:4").Return("", nil)
		mTmux.EXPECT().SwitchOrAttach("second-brain", model.ConnectOpts{}).Return("", nil)

		target, err := c.ConnectWindow(model.WindowConnectOpts{Name: "notes", Command: "claude"})

		require.NoError(t, err)
		assert.Equal(t, "second-brain:4", target)
	})

	t.Run("creates the window when the name doesn't match", func(t *testing.T) {
		c, mTmux, mLister, mDir, mHome := windowConnector(t)
		attachedTo(mLister, "second-brain", "/home/dev/brain")
		notADirectory(mHome, mDir, "notes")
		mTmux.EXPECT().ListWindows("second-brain").Return([]*model.TmuxWindow{{Name: "shell", Index: 1}}, nil)
		mTmux.EXPECT().NewWindowInSession(model.TmuxWindowOpts{
			Name:          "notes",
			StartDir:      "/home/dev/brain",
			TargetSession: "second-brain",
			Command:       "claude",
		}).Return("second-brain:2", nil)
		mTmux.EXPECT().SwitchOrAttach("second-brain", model.ConnectOpts{}).Return("", nil)

		target, err := c.ConnectWindow(model.WindowConnectOpts{Name: "notes", Command: "claude"})

		require.NoError(t, err)
		assert.Equal(t, "second-brain:2", target)
	})

	t.Run("always creates when no name is given", func(t *testing.T) {
		c, mTmux, mLister, _, _ := windowConnector(t)
		attachedTo(mLister, "second-brain", "/home/dev/brain")
		mTmux.EXPECT().NewWindowInSession(model.TmuxWindowOpts{
			StartDir:      "/home/dev/brain",
			TargetSession: "second-brain",
		}).Return("second-brain:2", nil)
		mTmux.EXPECT().SwitchOrAttach("second-brain", model.ConnectOpts{}).Return("", nil)

		target, err := c.ConnectWindow(model.WindowConnectOpts{})

		require.NoError(t, err)
		assert.Equal(t, "second-brain:2", target)
	})

	t.Run("--new creates even when the name matches, without listing windows", func(t *testing.T) {
		c, mTmux, mLister, mDir, mHome := windowConnector(t)
		attachedTo(mLister, "second-brain", "/home/dev/brain")
		notADirectory(mHome, mDir, "claude")
		mTmux.EXPECT().NewWindowInSession(model.TmuxWindowOpts{
			Name:          "claude",
			StartDir:      "/home/dev/brain",
			TargetSession: "second-brain",
			Command:       "claude",
			Background:    true,
		}).Return("second-brain:5", nil)

		target, err := c.ConnectWindow(model.WindowConnectOpts{
			Name: "claude", Command: "claude", New: true, Background: true,
		})

		require.NoError(t, err)
		assert.Equal(t, "second-brain:5", target)
	})

	t.Run("reuse picks the lowest-indexed window when names are shared", func(t *testing.T) {
		c, mTmux, mLister, mDir, mHome := windowConnector(t)
		attachedTo(mLister, "second-brain", "/home/dev/brain")
		notADirectory(mHome, mDir, "claude")
		mTmux.EXPECT().ListWindows("second-brain").Return([]*model.TmuxWindow{
			{Name: "claude", Index: 2},
			{Name: "claude", Index: 5},
		}, nil)
		mTmux.EXPECT().SendKeys("second-brain:2", "another prompt").Return("", nil)

		target, err := c.ConnectWindow(model.WindowConnectOpts{
			Name: "claude", Command: "another prompt", Background: true,
		})

		require.NoError(t, err)
		assert.Equal(t, "second-brain:2", target)
	})

	t.Run("background leaves the client where it is", func(t *testing.T) {
		c, mTmux, mLister, mDir, mHome := windowConnector(t)
		attachedTo(mLister, "second-brain", "/home/dev/brain")
		notADirectory(mHome, mDir, "notes")
		mTmux.EXPECT().ListWindows("second-brain").Return(nil, nil)
		mTmux.EXPECT().NewWindowInSession(model.TmuxWindowOpts{
			Name:          "notes",
			StartDir:      "/home/dev/brain",
			TargetSession: "second-brain",
			Background:    true,
		}).Return("second-brain:2", nil)

		target, err := c.ConnectWindow(model.WindowConnectOpts{Name: "notes", Background: true})

		require.NoError(t, err)
		assert.Equal(t, "second-brain:2", target)
	})

	t.Run("names a directory argument after its basename and roots it there", func(t *testing.T) {
		c, mTmux, mLister, mDir, mHome := windowConnector(t)
		attachedTo(mLister, "second-brain", "/home/dev/brain")
		mHome.EXPECT().ExpandPath("~/c/sesh").Return("/home/dev/c/sesh", nil)
		mDir.EXPECT().Dir("/home/dev/c/sesh").Return(true, "/home/dev/c/sesh")
		mTmux.EXPECT().ListWindows("second-brain").Return(nil, nil)
		mTmux.EXPECT().NewWindowInSession(model.TmuxWindowOpts{
			Name:          "sesh",
			StartDir:      "/home/dev/c/sesh",
			TargetSession: "second-brain",
		}).Return("second-brain:2", nil)
		mTmux.EXPECT().SwitchOrAttach("second-brain", model.ConnectOpts{}).Return("", nil)

		target, err := c.ConnectWindow(model.WindowConnectOpts{Name: "~/c/sesh"})

		require.NoError(t, err)
		assert.Equal(t, "second-brain:2", target)
	})

	t.Run("reuses the window a directory argument created on an earlier run", func(t *testing.T) {
		c, mTmux, mLister, mDir, mHome := windowConnector(t)
		attachedTo(mLister, "second-brain", "/home/dev/brain")
		mHome.EXPECT().ExpandPath("~/c/sesh").Return("/home/dev/c/sesh", nil)
		mDir.EXPECT().Dir("/home/dev/c/sesh").Return(true, "/home/dev/c/sesh")
		mTmux.EXPECT().ListWindows("second-brain").Return([]*model.TmuxWindow{{Name: "sesh", Index: 3}}, nil)
		mTmux.EXPECT().SelectWindow("second-brain:3").Return("", nil)
		mTmux.EXPECT().SwitchOrAttach("second-brain", model.ConnectOpts{}).Return("", nil)

		target, err := c.ConnectWindow(model.WindowConnectOpts{Name: "~/c/sesh"})

		require.NoError(t, err)
		assert.Equal(t, "second-brain:3", target)
	})

	t.Run("switches the client explicitly when triggered from outside tmux", func(t *testing.T) {
		c, mTmux, mLister, _, _ := windowConnector(t)
		mLister.EXPECT().FindTmuxSession("second-brain").Return(model.SeshSession{
			Src: "tmux", Name: "second-brain", Path: "/home/dev/brain",
		}, true)
		mTmux.EXPECT().IsAttached().Return(false)
		mTmux.EXPECT().NewWindowInSession(model.TmuxWindowOpts{
			StartDir:      "/home/dev/brain",
			TargetSession: "second-brain",
			Command:       "claude",
		}).Return("second-brain:2", nil)
		mTmux.EXPECT().ListClients().Return([]string{"/dev/ttys002"}, nil)
		mTmux.EXPECT().SwitchClientTarget("/dev/ttys002", "second-brain").Return("", nil)

		mZoxide := zoxide.NewMockZoxide(t)
		mZoxide.EXPECT().Add("/home/dev/brain").Return(nil)
		c.zoxide = mZoxide
		mFocuser := focuser.NewMockFocuser(t)
		mFocuser.EXPECT().Activate("").Return(false, nil)
		c.focuser = mFocuser

		target, err := c.ConnectWindow(model.WindowConnectOpts{
			Session: "second-brain", Command: "claude", Switch: true,
		})

		require.NoError(t, err)
		assert.Equal(t, "second-brain:2", target)
	})

	t.Run("errors when no session is given and tmux isn't attached", func(t *testing.T) {
		c, _, mLister, _, _ := windowConnector(t)
		mLister.EXPECT().GetAttachedTmuxSession().Return(model.SeshSession{}, false)

		_, err := c.ConnectWindow(model.WindowConnectOpts{Name: "notes"})

		assert.ErrorContains(t, err, "use --target")
	})
}
