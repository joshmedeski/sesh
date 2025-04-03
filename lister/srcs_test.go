package lister

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSortSources(t *testing.T) {
	tests := map[string]struct {
		sources   []string
		sortOrder []string
		expected  []string
	}{
		"a normal configuration": {
			sources:   []string{"tmux", "config", "tmuxinator", "zoxide"},
			sortOrder: []string{"tmuxinator", "zoxide", "config", "tmux"},
			expected:  []string{"tmuxinator", "zoxide", "config", "tmux"},
		},
		"empty configuration": {
			sources:   []string{"tmux", "config", "tmuxinator", "zoxide"},
			sortOrder: []string{},
			expected:  []string{"tmux", "config", "tmuxinator", "zoxide"},
		},
		"partial configuration": {
			sources:   []string{"tmux", "config", "tmuxinator", "zoxide"},
			sortOrder: []string{"tmuxinator"},
			expected:  []string{"tmuxinator", "tmux", "config", "zoxide"},
		},
		"superfluous elements": {
			sources:   []string{"tmux", "config", "tmuxinator", "zoxide"},
			sortOrder: []string{"tmuxinator", "apple", "zoxide", "banana", "config", "chocolate", "tmux"},
			expected:  []string{"tmuxinator", "zoxide", "config", "tmux"},
		},
		"configuration with capitalization": {
			sources:   []string{"tmux", "config", "tmuxinator", "zoxide"},
			sortOrder: []string{"tMuxiNator", "Zoxide", "conFIg", "tmux"},
			expected:  []string{"tmuxinator", "zoxide", "config", "tmux"},
		},
		"configuration with duplicate elements": {
			sources:   []string{"tmux", "config", "tmuxinator", "zoxide"},
			sortOrder: []string{"tmuxinator", "zoxide", "tmuxinator", "config", "tmuxinator", "tmux", "tmuxinator", "tmuxinator"},
			expected:  []string{"zoxide", "config", "tmux", "tmuxinator"},
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			sortSources(tt.sources, tt.sortOrder)
			assert.Equal(t, tt.expected, tt.sources)
		})
	}
}

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
