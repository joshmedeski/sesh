package tmux

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFormat(t *testing.T) {
	want := "#{session_activity} #{session_alerts} #{session_attached}" +
		" #{session_attached_list} #{session_created} #{session_format} " +
		"#{session_group} #{session_group_attached} " +
		"#{session_group_attached_list} #{session_group_list} " +
		"#{session_group_many_attached} #{session_group_size} " +
		"#{session_grouped} #{session_id} #{session_last_attached} " +
		"#{session_many_attached} #{session_marked} #{session_name} " +
		"#{session_path} #{session_stack} #{session_windows}"
	got := format()
	require.Equal(t, want, got)
}

func BenchmarkFormat(i *testing.B) {
	for n := 0; n < i.N; n++ {
		format()
	}
}

func TestProcessSessions(t *testing.T) {
	testCases := map[string]struct {
		Input    []string
		Expected []*TmuxSession
	}{
		"Single active session": {
			Input: []string{
				"1705879337  1 /dev/ttys000 1705878987 1       0 $2 1705879328 0 0 session-1 /some/test/path 1 1",
			},
			Expected: []*TmuxSession{},
		},
		"Single inactive session": {
			Input: []string{
				"1705879002  0  1705878987 1       0 $2 1705878987 0 0 session-1 /some/test/path 1 1",
			},
			Expected: make([]*TmuxSession, 1),
		},
		"Two inactive session": {
			Input: []string{
				"1705879002  0  1705878987 1       0 $2 1705878987 0 0 session-1 /some/test/path 1 1",
				"1705879063  0  1705879002 1       0 $3 1705879002 0 0 session-2 /some/other/test/path 1 1",
			},
			Expected: make([]*TmuxSession, 2),
		},
		"Two active session": {
			Input: []string{
				"1705879337  1 /dev/ttys000 1705878987 1       0 $2 1705879328 0 0 session-1 /some/test/path 1 1",
				"1705879337  1 /dev/ttys000 1705878987 1       0 $2 1705879328 0 0 session-1 /some/test/path 1 1",
			},
			Expected: []*TmuxSession{},
		},
		"No sessions": {
			Expected: []*TmuxSession{},
		},
		"Invalid LastAttached (Issue 34)": {
			Input: []string{
				"1705879002  0  1705878987 1       0 $2 1705878987 0 0 session-1 /some/test/path 1 1",
				"1705879063  0  1705879002 1       0 $3  0 0 session-2 /some/other/test/path 1 1",
			},
			Expected: make([]*TmuxSession, 2),
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			got := processSessions(tc.Input)
			require.Equal(t, len(tc.Expected), len(got))
		})
	}
}

func BenchmarkProcessSessions(b *testing.B) {
	for n := 0; n < b.N; n++ {
		processSessions([]string{
			"1705879337  1 /dev/ttys000 1705878987 1       0 $2 1705879328 0 0 session-1 /some/test/path 1 1",
			"1705879337  1 /dev/ttys000 1705878987 1       0 $2 1705879328 0 0 session-1 /some/test/path 1 1",
		})
	}
}
