package lister

import (
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/joshmedeski/sesh/v2/home"
	"github.com/joshmedeski/sesh/v2/model"
	"github.com/joshmedeski/sesh/v2/tmux"
	"github.com/joshmedeski/sesh/v2/tmuxinator"
	"github.com/joshmedeski/sesh/v2/zoxide"
	"github.com/stretchr/testify/assert"
)

func TestListTmuxSessions(t *testing.T) {
	mockTmux := new(tmux.MockTmux)
	t.Run("should list tmux sessions", func(t *testing.T) {
		const timeFormat = "2006-01-02 15:04:05 -0700 MST"
		createdFA, _ := time.Parse(timeFormat, "2024-04-25 19:02:45 -0500 CDT")
		lastAttachedFA, _ := time.Parse(timeFormat, "2024-04-25 19:30:06 -0500 CDT")
		activityFA, _ := time.Parse(timeFormat, "2024-04-25 19:44:06 -0500 CDT")
		firstAttached := model.TmuxSession{
			Created:           &createdFA,
			LastAttached:      &lastAttachedFA,
			Activity:          &activityFA,
			Group:             "",
			Path:              "/Users/joshmedeski/c/sesh/main",
			Name:              "sesh/main",
			ID:                "$1",
			AttachedList:      []string{""},
			GroupList:         []string{""},
			GroupAttachedList: []string{""},
			Stack:             []int{2, 1},
			Alerts:            []int{},
			GroupSize:         0,
			GroupAttached:     0,
			Attached:          0,
			Windows:           2,
			Format:            true,
			GroupManyAttached: false,
			Grouped:           false,
			ManyAttached:      false,
			Marked:            false,
		}

		createdLA, _ := time.Parse(timeFormat, "2024-04-25 19:02:45 -0500 CDT")
		lastAttachedLA, _ := time.Parse(timeFormat, "2024-04-25 19:44:06 -0500 CDT")
		activityLA, _ := time.Parse(timeFormat, "2024-04-25 19:44:06 -0500 CDT")
		lastAttached := model.TmuxSession{
			Created:           &createdLA,
			LastAttached:      &lastAttachedLA,
			Activity:          &activityLA,
			Group:             "",
			Path:              "/Users/joshmedeski/c/sesh/v2",
			Name:              "sesh/v2",
			ID:                "$1",
			AttachedList:      []string{""},
			GroupList:         []string{""},
			GroupAttachedList: []string{""},
			Stack:             []int{2, 1},
			Alerts:            []int{},
			GroupSize:         0,
			GroupAttached:     0,
			Attached:          0,
			Windows:           2,
			Format:            true,
			GroupManyAttached: false,
			Grouped:           false,
			ManyAttached:      false,
			Marked:            false,
		}
		mockTmux.On("ListSessions").Return([]*model.TmuxSession{&firstAttached, &lastAttached}, nil)

		mockConfig := model.Config{}
		mockHome := new(home.MockHome)
		mockZoxide := new(zoxide.MockZoxide)
		mockTmuxinator := new(tmuxinator.MockTmuxinator)
		lister := NewLister(mockConfig, mockHome, mockTmux, mockZoxide, mockTmuxinator)

		realLister, ok := lister.(*RealLister)
		if !ok {
			log.Fatal("Cannot convert lister to *RealLister")
		}

		sessions, err := listTmux(realLister)
		assert.Equal(t, "tmux:sesh/main", sessions.OrderedIndex[0])
		assert.Equal(t, "sesh/main", sessions.Directory["tmux:sesh/main"].Name)
		assert.Equal(t, "tmux:sesh/v2", sessions.OrderedIndex[1])
		assert.Equal(t, nil, err)
	})
}

func makeTmuxSession(name, path string) *model.TmuxSession {
	return &model.TmuxSession{
		Name:    name,
		Path:    path,
		ID:      "$1",
		Windows: 1,
	}
}

func TestGetLastTmuxSession(t *testing.T) {
	tests := []struct {
		name      string
		sessions  []*model.TmuxSession
		blacklist []string
		wantName  string
		wantOk    bool
	}{
		{
			name: "no blacklist returns second session",
			sessions: []*model.TmuxSession{
				makeTmuxSession("current", "/tmp/current"),
				makeTmuxSession("previous", "/tmp/previous"),
				makeTmuxSession("oldest", "/tmp/oldest"),
			},
			blacklist: nil,
			wantName:  "previous",
			wantOk:    true,
		},
		{
			name: "second session blacklisted skips to next",
			sessions: []*model.TmuxSession{
				makeTmuxSession("current", "/tmp/current"),
				makeTmuxSession("blocked", "/tmp/blocked"),
				makeTmuxSession("fallback", "/tmp/fallback"),
			},
			blacklist: []string{"blocked"},
			wantName:  "fallback",
			wantOk:    true,
		},
		{
			name: "all sessions blacklisted returns false",
			sessions: []*model.TmuxSession{
				makeTmuxSession("a", "/tmp/a"),
				makeTmuxSession("b", "/tmp/b"),
			},
			blacklist: []string{"a", "b"},
			wantOk:    false,
		},
		{
			name: "only one non-blacklisted session returns false",
			sessions: []*model.TmuxSession{
				makeTmuxSession("keep", "/tmp/keep"),
				makeTmuxSession("blocked", "/tmp/blocked"),
			},
			blacklist: []string{"blocked"},
			wantOk:    false,
		},
		{
			name: "regex blacklist pattern matches",
			sessions: []*model.TmuxSession{
				makeTmuxSession("current", "/tmp/current"),
				makeTmuxSession("scratch-123", "/tmp/scratch-123"),
				makeTmuxSession("work", "/tmp/work"),
			},
			blacklist: []string{"^scratch-.*"},
			wantName:  "work",
			wantOk:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockTmux := new(tmux.MockTmux)
			mockTmux.On("ListSessions").Return(tt.sessions, nil)

			config := model.Config{Blacklist: tt.blacklist}
			lister := NewLister(config, new(home.MockHome), mockTmux, new(zoxide.MockZoxide), new(tmuxinator.MockTmuxinator))

			session, ok := lister.GetLastTmuxSession()
			assert.Equal(t, tt.wantOk, ok)
			if tt.wantOk {
				assert.Equal(t, tt.wantName, session.Name)
			}
		})
	}
}

func TestListTmuxSessionsError(t *testing.T) {
	mockTmux := new(tmux.MockTmux)
	t.Run("should return error when unable to list tmux sessions", func(t *testing.T) {
		mockTmux.On("ListSessions").Return(nil, fmt.Errorf("some error"))

		mockConfig := model.Config{}
		mockHome := new(home.MockHome)
		mockZoxide := new(zoxide.MockZoxide)
		mockTmuxinator := new(tmuxinator.MockTmuxinator)
		lister := NewLister(mockConfig, mockHome, mockTmux, mockZoxide, mockTmuxinator)

		realLister, ok := lister.(*RealLister)
		if !ok {
			log.Fatal("Cannot convert lister to *RealLister")
		}

		sessions, err := listTmux(realLister)
		assert.Equal(t, model.SeshSessions{}, sessions)
		assert.NotNil(t, err)
		assert.Contains(t, err.Error(), "couldn't list tmux sessions")
	})
}
