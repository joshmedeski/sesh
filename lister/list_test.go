package lister

import (
	"testing"

	"github.com/joshmedeski/sesh/v2/home"
	"github.com/joshmedeski/sesh/v2/model"
	"github.com/joshmedeski/sesh/v2/projects"
	"github.com/joshmedeski/sesh/v2/tmux"
	"github.com/joshmedeski/sesh/v2/tmuxinator"
	"github.com/joshmedeski/sesh/v2/zoxide"
	"github.com/stretchr/testify/assert"
)

func TestHideDuplicates(t *testing.T) {
	tests := []struct {
		name              string
		tmuxSessions      []*model.TmuxSession
		zoxideResults     []*model.ZoxideResult
		configSessions    []model.SessionConfig
		tmuxinatorConfigs []*model.TmuxinatorConfig
		homeShortenHome   map[string]string
		expectedNames     []string
	}{
		{
			name: "no duplicates",
			tmuxSessions: []*model.TmuxSession{
				{Name: "session1", Path: "/path/to/session1"},
			},
			zoxideResults: []*model.ZoxideResult{
				{Path: "/path/to/session2", Score: 1.0},
			},
			homeShortenHome: map[string]string{
				"/path/to/session2": "session2",
			},
			expectedNames: []string{"session1", "session2"},
		},
		{
			name: "name duplicates only",
			tmuxSessions: []*model.TmuxSession{
				{Name: "dev", Path: "/path/to/dev1"},
			},
			zoxideResults: []*model.ZoxideResult{
				{Path: "/path/to/dev2", Score: 1.0},
			},
			homeShortenHome: map[string]string{
				"/path/to/dev2": "dev",
			},
			expectedNames: []string{"dev"},
		},
		{
			name: "path duplicates only",
			tmuxSessions: []*model.TmuxSession{
				{Name: "dev1", Path: "/path/to/dev"},
			},
			zoxideResults: []*model.ZoxideResult{
				{Path: "/path/to/dev", Score: 1.0},
			},
			homeShortenHome: map[string]string{
				"/path/to/dev": "dev2",
			},
			expectedNames: []string{"dev1"},
		},
		{
			name: "empty path tmuxinator sessions filtered by name only",
			tmuxinatorConfigs: []*model.TmuxinatorConfig{
				{Name: "tmux1"},
				{Name: "tmux1"},
				{Name: "tmux2"},
			},
			expectedNames: []string{"tmux1", "tmux2"},
		},
		{
			name: "mixed empty path and regular sessions with name conflicts",
			tmuxinatorConfigs: []*model.TmuxinatorConfig{
				{Name: "dev"},
			},
			zoxideResults: []*model.ZoxideResult{
				{Path: "/path/to/dev", Score: 1.0},
			},
			homeShortenHome: map[string]string{
				"/path/to/dev": "dev",
			},
			expectedNames: []string{"dev"},
		},
		{
			name: "order preservation",
			tmuxSessions: []*model.TmuxSession{
				{Name: "a", Path: "/path/to/a"},
			},
			zoxideResults: []*model.ZoxideResult{
				{Path: "/path/to/b", Score: 1.0},
			},
			configSessions: []model.SessionConfig{
				{Name: "a", Path: "/path/to/c"},
			},
			homeShortenHome: map[string]string{
				"/path/to/b": "b",
			},
			expectedNames: []string{"a", "b"},
		},
		{
			name:              "empty input",
			tmuxSessions:      []*model.TmuxSession{},
			zoxideResults:     []*model.ZoxideResult{},
			configSessions:    []model.SessionConfig{},
			tmuxinatorConfigs: []*model.TmuxinatorConfig{},
			expectedNames:     []string{},
		},
		{
			name: "single session",
			tmuxSessions: []*model.TmuxSession{
				{Name: "single", Path: "/path/to/single"},
			},
			expectedNames: []string{"single"},
		},
		{
			name: "all duplicates",
			tmuxSessions: []*model.TmuxSession{
				{Name: "dup", Path: "/path/to/dup"},
			},
			zoxideResults: []*model.ZoxideResult{
				{Path: "/path/to/dup", Score: 1.0},
			},
			configSessions: []model.SessionConfig{
				{Name: "dup", Path: "/path/to/dup"},
			},
			homeShortenHome: map[string]string{
				"/path/to/dup": "dup",
			},
			expectedNames: []string{"dup"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mocks
			mockTmux := new(tmux.MockTmux)
			mockZoxide := new(zoxide.MockZoxide)
			mockHome := new(home.MockHome)
			mockTmuxinator := new(tmuxinator.MockTmuxinator)

			// Set up mocks for all sources that have data (including empty data)
			if tt.tmuxSessions != nil {
				mockTmux.On("ListSessions").Return(tt.tmuxSessions, nil)
			}
			if tt.zoxideResults != nil {
				mockZoxide.On("ListResults").Return(tt.zoxideResults, nil)
				for path, shortened := range tt.homeShortenHome {
					mockHome.On("ShortenHome", path).Return(shortened, nil)
				}
			}
			if tt.configSessions != nil {
				for _, session := range tt.configSessions {
					mockHome.On("ExpandHome", session.Path).Return(session.Path, nil)
				}
			}
			if tt.tmuxinatorConfigs != nil {
				mockTmuxinator.On("List").Return(tt.tmuxinatorConfigs, nil)
			}

			config := model.Config{
				SessionConfigs: tt.configSessions,
			}
			mockProjects := new(projects.MockProjects)

			lister := NewLister(config, mockHome, mockTmux, mockZoxide, mockTmuxinator, mockProjects)

			// Call the actual List function with HideDuplicates
			result, err := lister.List(ListOptions{
				Tmux:           tt.tmuxSessions != nil,
				Zoxide:         tt.zoxideResults != nil,
				Config:         tt.configSessions != nil,
				Tmuxinator:     tt.tmuxinatorConfigs != nil,
				HideDuplicates: true,
			})

			assert.NoError(t, err)
			assert.Equal(t, len(tt.expectedNames), len(result.OrderedIndex))

			for i, expectedName := range tt.expectedNames {
				if i < len(result.OrderedIndex) {
					session := result.Directory[result.OrderedIndex[i]]
					assert.Equal(t, expectedName, session.Name)
				}
			}

			// Verify all mocks were called
			mockTmux.AssertExpectations(t)
			mockZoxide.AssertExpectations(t)
			mockTmuxinator.AssertExpectations(t)
		})
	}
}
