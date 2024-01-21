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

// BenchmarkFormat
// - Initial: BenchmarkFormat-10  649407  1830 ns/op  9936 B/op  41 allocs/op
func BenchmarkFormat(i *testing.B) {
	for n := 0; n < i.N; n++ {
		format()
	}
}
