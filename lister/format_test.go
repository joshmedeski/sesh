package lister

import (
	"testing"

	"github.com/joshmedeski/sesh/v2/formatter"
	"github.com/joshmedeski/sesh/v2/home"
	"github.com/joshmedeski/sesh/v2/icon"
	"github.com/joshmedeski/sesh/v2/model"
	"github.com/joshmedeski/sesh/v2/tmux"
	"github.com/joshmedeski/sesh/v2/tmuxinator"
	"github.com/joshmedeski/sesh/v2/zoxide"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestActiveWindowNameFormat(t *testing.T) {
	assert.Equal(t, "#{?window_active,#{window_name},}", activeWindowNameFormat)
}

func TestHasTmuxSessions(t *testing.T) {
	t.Run("true when tmux session is present", func(t *testing.T) {
		sessions := model.SeshSessions{
			OrderedIndex: []string{"config:notes", "tmux:work"},
			Directory: model.SeshSessionMap{
				"config:notes": {Src: "config", Name: "notes"},
				"tmux:work":    {Src: "tmux", Name: "work"},
			},
		}
		assert.True(t, hasTmuxSessions(sessions))
	})

	t.Run("false when no tmux session is present", func(t *testing.T) {
		sessions := model.SeshSessions{
			OrderedIndex: []string{"config:notes", "zoxide:work"},
			Directory: model.SeshSessionMap{
				"config:notes": {Src: "config", Name: "notes"},
				"zoxide:work":  {Src: "zoxide", Name: "work"},
			},
		}
		assert.False(t, hasTmuxSessions(sessions))
	})
}

func TestFirstActiveWindowNameBySession(t *testing.T) {
	t.Run("uses first returned name", func(t *testing.T) {
		got := firstActiveWindowNameBySession(map[string][]string{
			"work":     {"nvim", "shell"},
			"dotfiles": {"zsh"},
			"empty":    {},
		})
		assert.Equal(t, map[string]string{
			"work":     "nvim",
			"dotfiles": "zsh",
		}, got)
	})

	t.Run("nil for no window names", func(t *testing.T) {
		assert.Nil(t, firstActiveWindowNameBySession(nil))
		assert.Nil(t, firstActiveWindowNameBySession(map[string][]string{}))
	})
}

func TestFormatSessions(t *testing.T) {
	newLister := func(mockTmux *tmux.MockTmux) Lister {
		config := model.Config{}
		return NewLister(
			config,
			new(home.MockHome),
			mockTmux,
			new(zoxide.MockZoxide),
			new(tmuxinator.MockTmuxinator),
			formatter.NewFormatter(icon.NewIcon(config)),
			nil,
		)
	}

	t.Run("formats a cloned directory", func(t *testing.T) {
		mockTmux := new(tmux.MockTmux)
		l := newLister(mockTmux)
		sessions := model.SeshSessions{
			OrderedIndex: []string{"tmux:work"},
			Directory: model.SeshSessionMap{
				"tmux:work": {Src: "tmux", Name: "work", Path: "/work"},
			},
		}
		template := "{source}:{name}"

		got, err := l.Format(sessions, ListOptions{
			Format:    template,
			FormatSet: true,
			Icons:     true,
			NoColor:   true,
		})
		require.NoError(t, err)
		assert.Equal(t, "tmux: work", got.Directory["tmux:work"].Name)
		assert.Equal(t, "work", sessions.Directory["tmux:work"].Name)
	})

	t.Run("fetches active window names only when requested", func(t *testing.T) {
		mockTmux := new(tmux.MockTmux)
		mockTmux.On("ListAllWindowNames", activeWindowNameFormat).Return(map[string][]string{
			"work": {"nvim"},
		}, nil).Once()
		l := newLister(mockTmux)
		sessions := model.SeshSessions{
			OrderedIndex: []string{"tmux:work"},
			Directory: model.SeshSessionMap{
				"tmux:work": {Src: "tmux", Name: "work"},
			},
		}

		got, err := l.Format(sessions, ListOptions{
			Format:    "{active_window_name_prefix}{name}",
			FormatSet: true,
		})
		require.NoError(t, err)
		assert.Equal(t, "nvim work", got.Directory["tmux:work"].Name)
		mockTmux.AssertExpectations(t)
	})

	t.Run("does not fetch active window names for unrelated formats", func(t *testing.T) {
		mockTmux := new(tmux.MockTmux)
		l := newLister(mockTmux)
		sessions := model.SeshSessions{
			OrderedIndex: []string{"tmux:work"},
			Directory: model.SeshSessionMap{
				"tmux:work": {Src: "tmux", Name: "work"},
			},
		}

		_, err := l.Format(sessions, ListOptions{Format: "{name}", FormatSet: true})
		require.NoError(t, err)
		mockTmux.AssertNotCalled(t, "ListAllWindowNames")
	})
}

func TestFormatJsonResolvesSessions(t *testing.T) {
	config := model.Config{
		DefaultSessionConfig: model.DefaultSessionConfig{
			StartupCommand: "default-start",
			PreviewCommand: "default-preview {}",
			Windows:        []string{"shell"},
		},
		SessionConfigs: []model.SessionConfig{
			{Name: "notes", Path: "~/notes", Icon: "📓", Alias: "n", AliasAutoConnect: true,
				DefaultSessionConfig: model.DefaultSessionConfig{Tmuxp: "notes", Windows: []string{"editor"}}},
			{Name: "work", Path: "~/elsewhere", Alias: "w", AliasAutoConnect: true},
		},
		WindowConfigs: []model.WindowConfig{
			{Name: "editor", StartupScript: "nvim"},
			{Name: "shell", Path: "~/scratch"},
		},
		WildcardConfigs: []model.WildcardConfig{
			{Pattern: "~/c/quiet*", DisableStartCommand: true},
			{Pattern: "~/c/*", StartupCommand: "nvim", PreviewCommand: "glow {}", Windows: []string{"editor"}, Icon: "🏠"},
		},
	}
	mockTmux := new(tmux.MockTmux)
	mockTmux.On("ListAllWindows").Return(map[string][]model.TmuxWindow{
		"work": {{Name: "nvim", Index: 1, Active: true}, {Name: "shell", Index: 2}},
	}, nil).Once()
	l := NewLister(config, iconTestHome(t), mockTmux, nil, nil, nil, nil)
	notes, _ := l.FindConfigSession("notes")
	sessions := model.SeshSessions{
		OrderedIndex: []string{"config:notes", "zoxide:app", "zoxide:quiet", "zoxide:other", "tmux:work"},
		Directory: model.SeshSessionMap{
			"config:notes": notes,
			"zoxide:app":   {Src: "zoxide", Name: "~/c/app", Path: "/home/user/c/app"},
			"zoxide:quiet": {Src: "zoxide", Name: "~/c/quiet", Path: "/home/user/c/quiet"},
			"zoxide:other": {Src: "zoxide", Name: "~/other", Path: "/home/user/other"},
			"tmux:work":    {Src: "tmux", Name: "work", Path: "/home/user/c/work"},
		},
	}

	got, err := l.Format(sessions, ListOptions{Json: true})
	require.NoError(t, err)

	assert.Equal(t, model.SeshSession{
		Src: "config", Name: "notes", Path: "/home/user/notes",
		StartupCommand: "default-start", PreviewCommand: "default-preview {}",
		WindowNames:   []string{"editor"},
		WindowConfigs: []model.WindowConfig{{Name: "editor", StartupScript: "nvim", Path: "/home/user/notes"}},
		Icon:          "📓", Alias: "n", AliasAutoConnect: true, Tmuxp: "notes",
	}, got.Directory["config:notes"])

	assert.Equal(t, model.SeshSession{
		Src: "zoxide", Name: "~/c/app", Path: "/home/user/c/app",
		StartupCommand: "nvim", PreviewCommand: "glow {}",
		WindowNames:   []string{"editor"},
		WindowConfigs: []model.WindowConfig{{Name: "editor", StartupScript: "nvim", Path: "/home/user/c/app"}},
		Icon:          "🏠", Wildcard: "~/c/*",
	}, got.Directory["zoxide:app"])

	quiet := got.Directory["zoxide:quiet"]
	assert.True(t, quiet.DisableStartupCommand)
	assert.Empty(t, quiet.StartupCommand, "a wildcard that disables the startup command skips the default too")
	assert.Equal(t, "~/c/quiet*", quiet.Wildcard)

	other := got.Directory["zoxide:other"]
	assert.Equal(t, "default-start", other.StartupCommand)
	assert.Equal(t, []model.WindowConfig{{Name: "shell", Path: "/home/user/scratch"}}, other.WindowConfigs)
	assert.Empty(t, other.Wildcard)

	work := got.Directory["tmux:work"]
	assert.Equal(t, []model.TmuxWindow{{Name: "nvim", Index: 1, Active: true}, {Name: "shell", Index: 2}}, work.TmuxWindows)
	assert.Equal(t, "🏠", work.Icon)
	assert.Empty(t, work.StartupCommand, "a running tmux session starts nothing on connect")
	assert.Nil(t, work.WindowNames)
	assert.Equal(t, "w", work.Alias, "a tmux session named after a [[session]] keeps its alias, even if --hide-duplicates drops the config entry")
	assert.True(t, work.AliasAutoConnect)

	assert.Empty(t, sessions.Directory["zoxide:app"].StartupCommand, "the input list is left untouched")
	mockTmux.AssertExpectations(t)
}
