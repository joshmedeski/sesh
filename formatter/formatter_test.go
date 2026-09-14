package formatter

import (
	"testing"

	"github.com/joshmedeski/sesh/v2/icon"
	"github.com/joshmedeski/sesh/v2/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestColorCode(t *testing.T) {
	tests := []struct {
		name  string
		color string
		code  int
		ok    bool
	}{
		{"standard color", "blue", 34, true},
		{"bright color", "bright-blue", 94, true},
		{"case and whitespace", "  Bright-Magenta  ", 95, true},
		{"gray alias", "gray", 90, true},
		{"grey alias", "grey", 90, true},
		{"empty color", "", 0, true},
		{"unsupported color", "orange", 0, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			code, ok := colorCode(tt.color)
			assert.Equal(t, tt.code, code)
			assert.Equal(t, tt.ok, ok)
		})
	}
}

func TestRenderColors(t *testing.T) {
	t.Run("renders foreground color and reset", func(t *testing.T) {
		got, err := renderColors("{fg:blue}hello{/fg}", false)
		require.NoError(t, err)
		assert.Equal(t, "\x1b[34mhello\x1b[39m", got)
	})

	t.Run("renders multiple scopes", func(t *testing.T) {
		got, err := renderColors(
			"{fg:red}red{/fg} {fg:bright-blue}blue{/fg}",
			false,
		)
		require.NoError(t, err)
		assert.Equal(t, "\x1b[31mred\x1b[39m \x1b[94mblue\x1b[39m", got)
	})

	t.Run("removes color tokens with no color", func(t *testing.T) {
		got, err := renderColors("{fg:blue}hello{/fg}", true)
		require.NoError(t, err)
		assert.Equal(t, "hello", got)
	})

	t.Run("preserves placeholders", func(t *testing.T) {
		got, err := renderColors(
			"{fg:blue}{active_window_name_prefix}{/fg}{name}",
			false,
		)
		require.NoError(t, err)
		assert.Equal(
			t,
			"\x1b[34m{active_window_name_prefix}\x1b[39m{name}",
			got,
		)
	})

	t.Run("rejects unsupported colors", func(t *testing.T) {
		_, err := renderColors("{fg:orange}hello{/fg}", false)
		require.Error(t, err)
		assert.Contains(t, err.Error(), `unsupported format color "orange"`)
	})

	t.Run("rejects empty colors", func(t *testing.T) {
		_, err := renderColors("{fg:}hello{/fg}", false)
		require.Error(t, err)
		assert.Contains(t, err.Error(), `unsupported format color ""`)
	})

	t.Run("rejects unterminated tokens", func(t *testing.T) {
		_, err := renderColors("{fg:blue", false)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "unterminated foreground color token")
	})
}

func TestUsesActiveWindowName(t *testing.T) {
	f := NewFormatter(icon.NewIcon(model.Config{}))

	assert.True(t, f.UsesActiveWindowName("{active_window_name} {name}"))
	assert.True(t, f.UsesActiveWindowName("{active_window_name_prefix}{name}"))
	assert.False(t, f.UsesActiveWindowName("{source} {name}"))
}

func TestFormat(t *testing.T) {
	f := NewFormatter(icon.NewIcon(model.Config{}))
	session := model.SeshSession{
		Src:  "tmux",
		Name: "work",
		Path: "/home/user/work",
	}
	original := model.SeshSessions{
		OrderedIndex: []string{"tmux:work"},
		Directory: model.SeshSessionMap{
			"tmux:work": session,
		},
	}
	formatOne := func(t *testing.T, opts Options) string {
		t.Helper()
		got, err := f.Format(original, opts)
		require.NoError(t, err)
		return got.Directory["tmux:work"].Name
	}

	t.Run("replaces all supported fields", func(t *testing.T) {
		template := "{source}|{session}|{path}|{active_window_name}|{active_window_name_prefix}{name}"
		got := formatOne(t, Options{
			Template:          &template,
			Icons:             true,
			NoColor:           true,
			ActiveWindowNames: map[string]string{"work": ""},
		})
		assert.Equal(t, "tmux|work|/home/user/work||  work", got)
	})

	t.Run("prefix has one trailing space", func(t *testing.T) {
		template := "{active_window_name_prefix}{name}"
		got := formatOne(t, Options{
			Template:          &template,
			ActiveWindowNames: map[string]string{"work": ""},
		})
		assert.Equal(t, " work", got)
	})

	t.Run("empty prefix disappears completely", func(t *testing.T) {
		template := "{active_window_name_prefix}{name}"
		got := formatOne(t, Options{Template: &template})
		assert.Equal(t, "work", got)
	})

	t.Run("active window name has no implicit spacing", func(t *testing.T) {
		template := "{active_window_name}:{name}"
		got := formatOne(t, Options{
			Template:          &template,
			ActiveWindowNames: map[string]string{"work": ""},
		})
		assert.Equal(t, ":work", got)
	})

	t.Run("works with inline colors", func(t *testing.T) {
		template := "{fg:blue}{active_window_name_prefix}{/fg}{name}"
		got := formatOne(t, Options{
			Template:          &template,
			ActiveWindowNames: map[string]string{"work": ""},
		})
		assert.Equal(t, "\x1b[34m \x1b[39mwork", got)
	})

	t.Run("adds icons without a template", func(t *testing.T) {
		got := formatOne(t, Options{Icons: true, NoColor: true})
		assert.Equal(t, " work", got)
	})

	t.Run("can exclude a source icon", func(t *testing.T) {
		got := formatOne(t, Options{
			Icons:        true,
			NoColor:      true,
			IconExcludes: []string{" TMUX "},
		})
		assert.Equal(t, "work", got)
	})

	t.Run("does not mutate the input directory", func(t *testing.T) {
		got := formatOne(t, Options{Icons: true, NoColor: true})
		assert.Equal(t, " work", got)
		assert.Equal(t, "work", original.Directory["tmux:work"].Name)
	})

	t.Run("validates a template even for an empty list", func(t *testing.T) {
		template := "{fg:orange}{name}{/fg}"
		_, err := f.Format(model.SeshSessions{}, Options{Template: &template})
		require.Error(t, err)
		assert.Contains(t, err.Error(), `unsupported format color "orange"`)
	})
}
