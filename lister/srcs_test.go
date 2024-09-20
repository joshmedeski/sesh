package lister

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSrcs(t *testing.T) {
	tests := []struct {
		name     string
		expected []string
		opts     ListOptions
	}{
		{
			name:     "All options are false",
			opts:     ListOptions{},
			expected: []string{"tmux", "config", "tmuxinator", "zoxide"},
		},
		{
			name:     "Only Tmux is true",
			opts:     ListOptions{Tmux: true},
			expected: []string{"tmux"},
		},
		{
			name:     "Only Config is true",
			opts:     ListOptions{Config: true},
			expected: []string{"config"},
		},
		{
			name:     "Only Zoxide is true",
			opts:     ListOptions{Zoxide: true},
			expected: []string{"zoxide"},
		},
		{
			name:     "Tmux and Config are true",
			opts:     ListOptions{Tmux: true, Config: true},
			expected: []string{"tmux", "config"},
		},
		{
			name:     "Tmux and Zoxide are true",
			opts:     ListOptions{Tmux: true, Zoxide: true},
			expected: []string{"tmux", "zoxide"},
		},
		{
			name:     "Config and Zoxide are true",
			opts:     ListOptions{Config: true, Zoxide: true},
			expected: []string{"config", "zoxide"},
		},
		{
			name:     "All options are true",
			opts:     ListOptions{Tmux: true, Config: true, Zoxide: true},
			expected: []string{"tmux", "config", "zoxide"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := srcs(tt.opts)
			assert.Equal(t, tt.expected, result)
		})
	}
}
