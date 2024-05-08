package lister

import (
	"errors"
	"testing"

	"github.com/joshmedeski/sesh/config"
	"github.com/joshmedeski/sesh/home"
	"github.com/joshmedeski/sesh/model"
	"github.com/joshmedeski/sesh/tmux"
	"github.com/joshmedeski/sesh/zoxide"
	"github.com/stretchr/testify/assert"
)

func TestList_withTmux(t *testing.T) {
	mockTmux := new(tmux.MockTmux)
	mockConfig := new(config.MockConfig)
	mockZoxide := new(zoxide.MockZoxide)
	lister := NewLister(mockConfig, nil, mockTmux, mockZoxide)

	mockTmuxSessions := []*model.TmuxSession{{Name: "test", Path: "/test", Attached: 1}}
	mockTmux.On("ListSessions").Return(mockTmuxSessions, nil)

	list, err := lister.List(ListOptions{Tmux: true})
	assert.NoError(t, err)
	assert.Len(t, list, 1)
	assert.Equal(t, "test", list[0].Name)
	mockTmux.AssertExpectations(t)
}

func TestList_withConfig(t *testing.T) {
	mockTmux := new(tmux.MockTmux)
	mockConfig := new(config.MockConfig)
	mockZoxide := new(zoxide.MockZoxide)
	lister := NewLister(mockConfig, nil, mockTmux, mockZoxide)

	mockConfigSessions := model.Config{SessionConfigs: []model.SessionConfig{{Name: "configSession", Path: "/config"}}}
	mockConfig.On("GetConfig").Return(mockConfigSessions, nil)

	list, err := lister.List(ListOptions{Config: true})
	assert.NoError(t, err)
	assert.Len(t, list, 1)
	assert.Equal(t, "configSession", list[0].Name)
	mockConfig.AssertExpectations(t)
}

func TestList_withZoxide(t *testing.T) {
	mockTmux := new(tmux.MockTmux)
	mockConfig := new(config.MockConfig)
	mockZoxide := new(zoxide.MockZoxide)
	mockHome := new(home.MockHome)
	lister := NewLister(mockConfig, mockHome, mockTmux, mockZoxide)

	mockZoxideResults := []*model.ZoxideResult{{Path: "/zoxidePath", Score: 0.5}}
	mockZoxide.On("ListResults").Return(mockZoxideResults, nil)

	list, err := lister.List(ListOptions{Zoxide: true})
	assert.NoError(t, err)
	assert.Len(t, list, 1)
	assert.Equal(t, "/zoxidePath", list[0].Path)
	mockZoxide.AssertExpectations(t)
}

func TestList_Errors(t *testing.T) {
	mockTmux := new(tmux.MockTmux)
	mockConfig := new(config.MockConfig)
	mockHome := new(home.MockHome)
	mockZoxide := new(zoxide.MockZoxide)
	lister := NewLister(mockConfig, mockHome, mockTmux, mockZoxide)

	mockTmux.On("ListSessions").Return(nil, errors.New("tmux error"))
	mockConfig.On("GetConfig").Return(model.Config{}, errors.New("config error"))
	mockZoxide.On("ListResults").Return(nil, errors.New("zoxide error"))

	_, err := lister.List(ListOptions{Tmux: true})
	assert.Error(t, err, "tmux error")

	_, err = lister.List(ListOptions{Config: true})
	assert.Error(t, err, "config error")

	_, err = lister.List(ListOptions{Zoxide: true})
	assert.Error(t, err, "zoxide error")
}
