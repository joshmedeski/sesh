package seshcli

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/joshmedeski/sesh/v2/model"
	"github.com/joshmedeski/sesh/v2/picker"
	"github.com/joshmedeski/sesh/v2/worktree"
)

func TestWorktreeLabel(t *testing.T) {
	assert.Equal(t, "409  Worktree support",
		worktreeLabel(model.WorktreeEntry{Number: 409, Title: "Worktree support"}))
	assert.Equal(t, "409", worktreeLabel(model.WorktreeEntry{Number: 409}),
		"a worktree whose title never resolved shows its bare number")
}

func TestPickWorktree(t *testing.T) {
	listOpts := model.WorktreeListOpts{Repo: "joshmedeski/sesh"}

	t.Run("connects to the picked worktree", func(t *testing.T) {
		mockPicker := new(picker.MockPicker)
		mockWorktree := new(worktree.MockWorktree)
		mockPicker.On("PickWorktree", mock.Anything, mock.Anything).
			Return(model.WorktreeEntry{Number: 409, Path: "/repo/.wk/409"}, true, nil)
		mockWorktree.On("Connect", model.WorktreeConnectOpts{
			Number: 409,
			Repo:   "joshmedeski/sesh",
			Switch: true,
		}).Return("409", nil)

		deps := &Deps{Picker: mockPicker, Worktree: mockWorktree}
		assert.NoError(t, pickWorktree(deps, listOpts, true, picker.WorktreePickerOptions{}))

		mockWorktree.AssertExpectations(t)
	})

	t.Run("connects to nothing when the picker is quit", func(t *testing.T) {
		mockPicker := new(picker.MockPicker)
		mockWorktree := new(worktree.MockWorktree)
		mockPicker.On("PickWorktree", mock.Anything, mock.Anything).
			Return(model.WorktreeEntry{}, false, nil)

		deps := &Deps{Picker: mockPicker, Worktree: mockWorktree}
		assert.NoError(t, pickWorktree(deps, listOpts, false, picker.WorktreePickerOptions{}))

		mockWorktree.AssertNotCalled(t, "Connect", mock.Anything)
	})

	t.Run("returns the picker's error without connecting", func(t *testing.T) {
		mockPicker := new(picker.MockPicker)
		mockWorktree := new(worktree.MockWorktree)
		mockPicker.On("PickWorktree", mock.Anything, mock.Anything).
			Return(model.WorktreeEntry{}, false, errors.New("no [[worktree]] block"))

		deps := &Deps{Picker: mockPicker, Worktree: mockWorktree}
		assert.Error(t, pickWorktree(deps, listOpts, false, picker.WorktreePickerOptions{}))

		mockWorktree.AssertNotCalled(t, "Connect", mock.Anything)
	})

	t.Run("lists through the list opts it was given", func(t *testing.T) {
		mockPicker := new(picker.MockPicker)
		mockWorktree := new(worktree.MockWorktree)
		mockWorktree.On("List", listOpts).Return([]model.WorktreeEntry{{Number: 409}}, nil)
		// The picker is handed a fetch function rather than a list, so calling it
		// is what proves the opts made it through.
		mockPicker.On("PickWorktree", mock.Anything, mock.Anything).
			Run(func(args mock.Arguments) {
				fetchFunc := args.Get(0).(picker.WorktreeFetchFunc)
				entries, err := fetchFunc(false)
				assert.NoError(t, err)
				assert.Len(t, entries, 1)
			}).
			Return(model.WorktreeEntry{}, false, nil)

		deps := &Deps{Picker: mockPicker, Worktree: mockWorktree}
		assert.NoError(t, pickWorktree(deps, listOpts, false, picker.WorktreePickerOptions{}))

		mockWorktree.AssertExpectations(t)
	})

	t.Run("a refreshing fetch lists with Refresh set", func(t *testing.T) {
		mockPicker := new(picker.MockPicker)
		mockWorktree := new(worktree.MockWorktree)
		refreshOpts := listOpts
		refreshOpts.Refresh = true
		mockWorktree.On("List", refreshOpts).Return([]model.WorktreeEntry{{Number: 409}}, nil)
		mockPicker.On("PickWorktree", mock.Anything, mock.Anything).
			Run(func(args mock.Arguments) {
				fetchFunc := args.Get(0).(picker.WorktreeFetchFunc)
				_, err := fetchFunc(true)
				assert.NoError(t, err)
			}).
			Return(model.WorktreeEntry{}, false, nil)

		deps := &Deps{Picker: mockPicker, Worktree: mockWorktree}
		assert.NoError(t, pickWorktree(deps, listOpts, false, picker.WorktreePickerOptions{}))

		mockWorktree.AssertExpectations(t)
	})
}
