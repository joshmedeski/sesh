package tmux

import (
	"testing"
	"time"

	"github.com/joshmedeski/sesh/v2/model"
	"github.com/joshmedeski/sesh/v2/shell"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestListSessions(t *testing.T) {
	t.Run("List tmux session", func(t *testing.T) {
		mockShell := &shell.MockShell{}
		tmux := &RealTmux{shell: mockShell}
		mockShell.EXPECT().ListCmd("tmux", "list-sessions", "-F", mock.Anything).Return([]string{"1714092246::::0::::1714089765::1::::::::::::::0::$1::1714092246::0::0::sesh/main::/Users/joshmedeski/c/sesh/main::2,1::2"},
			nil,
		)
		sessions, err := tmux.ListSessions()
		assert.Nil(t, err)
		for _, session := range sessions {
			assert.Equal(t, "sesh/main", session.Name)
			assert.Equal(t, "/Users/joshmedeski/c/sesh/main", session.Path)
		}
	})

	t.Run("parseTmuxSessionsOutput", func(t *testing.T) {
		rawSessions := []string{
			"1714092246::::0::::1714089765::1::::::::::::::0::$1::1714092246::0::0::sesh/main::/Users/joshmedeski/c/sesh/main::2,1::2",
		}
		sessions, err := parseTmuxSessionsOutput(rawSessions)
		assert.Nil(t, err)

		expectedName := "sesh/main"
		expectedPath := "/Users/joshmedeski/c/sesh/main"

		for _, session := range sessions {
			assert.Equal(t, expectedName, session.Name)
			assert.Equal(t, expectedPath, session.Path)
		}
	})

	t.Run("sortByLastAttached", func(t *testing.T) {
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

		expectedSortedSessions := []*model.TmuxSession{&lastAttached, &firstAttached}
		actualSortedSessionsOutOfOrder := sortByLastAttached([]*model.TmuxSession{&firstAttached, &lastAttached})
		assert.Equal(t, expectedSortedSessions, actualSortedSessionsOutOfOrder)
		actualSortedSessionsInOrder := sortByLastAttached([]*model.TmuxSession{&lastAttached, &firstAttached})
		assert.Equal(t, expectedSortedSessions, actualSortedSessionsInOrder)
	})
}
