package lister

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/joshmedeski/sesh/v2/model"
)

func TestSortSources(t *testing.T) {
	defaultSources := []string{"tmux", "config", "tmuxinator", "zoxide"}
	tests := map[string]struct {
		sortOrder model.SortOrder
		expected  []string
	}{
		"a normal configuration": {
			sortOrder: model.SortOrder{"tmuxinator", "zoxide", "config", "tmux"},
			expected:  []string{"tmuxinator", "zoxide", "config", "tmux"},
		},
		"empty configuration": {
			sortOrder: model.SortOrder{},
			expected:  []string{"tmux", "config", "tmuxinator", "zoxide"},
		},
		"partial configuration": {
			sortOrder: model.SortOrder{"tmuxinator"},
			expected:  []string{"tmuxinator", "tmux", "config", "zoxide"},
		},
		"superfluous elements": {
			sortOrder: model.SortOrder{"tmuxinator", "apple", "zoxide", "banana", "config", "chocolate", "tmux"},
			expected:  []string{"tmuxinator", "zoxide", "config", "tmux"},
		},
		"configuration with capitalization": {
			sortOrder: model.SortOrder{"tMuxiNator", "Zoxide", "conFIg", "tmux"},
			expected:  []string{"tmuxinator", "zoxide", "config", "tmux"},
		},
		"configuration with duplicate elements": {
			sortOrder: model.SortOrder{"tmuxinator", "zoxide", "tmuxinator", "config", "tmuxinator", "tmux", "tmuxinator", "tmuxinator"},
			expected:  []string{"zoxide", "config", "tmux", "tmuxinator"},
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			actual := sortSources(defaultSources, tt.sortOrder)
			assert.Equal(t, tt.expected, actual)
		})
	}
}

func TestGroupSources(t *testing.T) {
	defaultSources := []string{"tmux", "config", "tmuxinator", "zoxide"}
	tests := map[string]struct {
		sortOrder model.SortOrder
		expected  [][]string
	}{
		"no sort order leaves every source on its own": {
			sortOrder: nil,
			expected:  [][]string{{"tmux"}, {"config"}, {"tmuxinator"}, {"zoxide"}},
		},
		"a flat sort order leaves every source on its own": {
			sortOrder: model.SortOrder{"zoxide", "tmux", "config", "tmuxinator"},
			expected:  [][]string{{"zoxide"}, {"tmux"}, {"config"}, {"tmuxinator"}},
		},
		"a nested entry merges its sources into one group": {
			sortOrder: model.SortOrder{"tmux", []string{"config", "zoxide"}, "tmuxinator"},
			expected:  [][]string{{"tmux"}, {"config", "zoxide"}, {"tmuxinator"}},
		},
		"a nested entry decoded from toml merges too": {
			sortOrder: model.SortOrder{"tmux", []any{"config", "zoxide"}},
			expected:  [][]string{{"tmux"}, {"config", "zoxide"}, {"tmuxinator"}},
		},
		"unlisted sources trail as groups of one": {
			sortOrder: model.SortOrder{[]string{"zoxide", "tmux"}},
			expected:  [][]string{{"zoxide", "tmux"}, {"config"}, {"tmuxinator"}},
		},
		"a source named twice joins the last group that claims it": {
			sortOrder: model.SortOrder{[]string{"tmux", "zoxide"}, []string{"config", "zoxide"}},
			expected:  [][]string{{"tmux"}, {"config", "zoxide"}, {"tmuxinator"}},
		},
		"a group naming only inactive sources is skipped": {
			sortOrder: model.SortOrder{[]string{"apple", "banana"}, "zoxide"},
			expected:  [][]string{{"zoxide"}, {"tmux"}, {"config"}, {"tmuxinator"}},
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, tt.expected, groupSources(defaultSources, tt.sortOrder))
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
		{
			name:     "Only Panes is true",
			opts:     ListOptions{Panes: true},
			expected: []string{"tmux-pane"},
		},
		{
			name:     "Panes with Tmux",
			opts:     ListOptions{Panes: true, Tmux: true},
			expected: []string{"tmux", "tmux-pane"},
		},
		{
			name:     "Panes with Zoxide",
			opts:     ListOptions{Panes: true, Zoxide: true},
			expected: []string{"zoxide", "tmux-pane"},
		},
		{
			name:     "Panes with all sources",
			opts:     ListOptions{Panes: true, Tmux: true, Config: true, Zoxide: true, Tmuxinator: true},
			expected: []string{"tmux", "config", "tmuxinator", "zoxide", "tmux-pane"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := srcs(tt.opts)
			assert.Equal(t, tt.expected, result)
		})
	}
}
