package lister

import (
	"fmt"
	"testing"

	"github.com/joshmedeski/sesh/v2/model"
	"github.com/stretchr/testify/assert"
)

func mkSess(src, name, path string) model.SeshSession {
	return model.SeshSession{Src: src, Name: name, Path: path}
}

func mkSessions(ss ...model.SeshSession) model.SeshSessions {
	out := model.SeshSessions{
		Directory:    make(model.SeshSessionMap, len(ss)),
		OrderedIndex: make([]string, 0, len(ss)),
	}
	for i, s := range ss {
		key := fmt.Sprintf("%s:%d", s.Src, i)
		out.Directory[key] = s
		out.OrderedIndex = append(out.OrderedIndex, key)
	}
	return out
}

func TestApplyDedup(t *testing.T) {
	cases := []struct {
		name     string
		input    model.SeshSessions
		expected []string
		// expectedSrcs is set for cases where two candidates share the
		// rendered Name/Path; asserting on Src ensures the right source
		// survives.
		expectedSrcs []string
	}{
		{
			name:     "two_tmux_same_path_different_names",
			input:    mkSessions(mkSess("tmux", "a", "/p"), mkSess("tmux", "b", "/p")),
			expected: []string{"a", "b"},
		},
		{
			name:         "tmux_and_config_same_name",
			input:        mkSessions(mkSess("tmux", "myproj", "/x"), mkSess("config", "myproj", "/y")),
			expected:     []string{"myproj"},
			expectedSrcs: []string{"tmux"},
		},
		{
			name:         "tmux_and_zoxide_same_path",
			input:        mkSessions(mkSess("tmux", "t", "/a"), mkSess("zoxide", "", "/a")),
			expected:     []string{"t"},
			expectedSrcs: []string{"tmux"},
		},
		{
			name:     "tmux_and_config_same_path_diff_names",
			input:    mkSessions(mkSess("tmux", "x", "/a"), mkSess("config", "y", "/a")),
			expected: []string{"x", "y"},
		},
		{
			name:     "config_and_tmuxinator_same_name",
			input:    mkSessions(mkSess("config", "myproj", "/a"), mkSess("tmuxinator", "myproj", "")),
			expected: []string{"myproj", "myproj"},
		},
		{
			name:         "config_and_zoxide_same_path",
			input:        mkSessions(mkSess("config", "x", "/x.path"), mkSess("zoxide", "", "/x.path")),
			expected:     []string{"x"},
			expectedSrcs: []string{"config"},
		},
		{
			name:     "two_config_same_path_diff_names",
			input:    mkSessions(mkSess("config", "a", "/p"), mkSess("config", "b", "/p")),
			expected: []string{"a", "b"},
		},
		{
			name:     "empty_input",
			input:    mkSessions(),
			expected: []string{},
		},
		{
			name:     "single_source_only_tmux",
			input:    mkSessions(mkSess("tmux", "a", "/1"), mkSess("tmux", "b", "/2"), mkSess("tmux", "c", "/3")),
			expected: []string{"a", "b", "c"},
		},
		{
			name:     "bug_repro_two_tmux_same_path",
			input:    mkSessions(mkSess("tmux", "ddd-sandbox", "/s"), mkSess("tmux", "claude c02df2feed", "/s")),
			expected: []string{"ddd-sandbox", "claude c02df2feed"},
		},
		{
			name:         "order_independence_zoxide_before_tmux",
			input:        mkSessions(mkSess("zoxide", "", "/a"), mkSess("tmux", "t", "/a")),
			expected:     []string{"t"},
			expectedSrcs: []string{"tmux"},
		},
		{
			// tmuxinator currently has no Path, so a path-keyed dedup
			// rule cannot match. Both entries must survive.
			name:     "zoxide_and_tmuxinator_kept_no_path_on_tmuxinator",
			input:    mkSessions(mkSess("tmuxinator", "myproj", ""), mkSess("zoxide", "~/code/myproj", "/home/u/code/myproj")),
			expected: []string{"myproj", "~/code/myproj"},
		},
		{
			// Cascading drop guard: config loses to tmux on name, but
			// zoxide must still be kept because it would only have lost
			// to the now-dropped config's path, not to any kept session.
			name: "zoxide_survives_when_only_dropped_config_owned_path",
			input: mkSessions(
				mkSess("tmux", "foo", "/q"),
				mkSess("config", "foo", "/p"),
				mkSess("zoxide", "", "/p"),
			),
			expected: []string{"foo", "/p"},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			result := applyDedup(tc.input)
			got := make([]string, 0, len(result))
			gotSrcs := make([]string, 0, len(result))
			for _, idx := range result {
				s := tc.input.Directory[idx]
				if s.Name != "" {
					got = append(got, s.Name)
				} else {
					got = append(got, s.Path)
				}
				gotSrcs = append(gotSrcs, s.Src)
			}
			assert.Equal(t, tc.expected, got)
			if tc.expectedSrcs != nil {
				assert.Equal(t, tc.expectedSrcs, gotSrcs)
			}
		})
	}
}
