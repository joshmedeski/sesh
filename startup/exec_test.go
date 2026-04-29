package startup

import (
	"testing"

	"github.com/joshmedeski/sesh/v2/home"
	"github.com/joshmedeski/sesh/v2/lister"
	"github.com/joshmedeski/sesh/v2/model"
	"github.com/joshmedeski/sesh/v2/oswrap"
	"github.com/joshmedeski/sesh/v2/replacer"
	"github.com/joshmedeski/sesh/v2/tmux"
	"github.com/stretchr/testify/assert"
)

func TestExecCreatesWindowsInTargetSession(t *testing.T) {
	mockOs := new(oswrap.MockOs)
	mockLister := new(lister.MockLister)
	mockTmux := new(tmux.MockTmux)
	mockHome := new(home.MockHome)
	mockReplacer := new(replacer.MockReplacer)

	s := &RealStartup{
		os:       mockOs,
		lister:   mockLister,
		tmux:     mockTmux,
		config:   model.Config{WindowConfigs: []model.WindowConfig{{Name: "editor", StartupScript: "echo hi"}}},
		home:     mockHome,
		replacer: mockReplacer,
	}
	session := model.SeshSession{Name: "demo", Path: "/tmp", WindowNames: []string{"editor"}}

	mockHome.On("ExpandPath", "/tmp").Return("/tmp", nil)
	mockOs.On("Getenv", "SHELL").Return("/bin/zsh")
	mockTmux.On("NewWindowInSession", "editor", "/tmp", "demo", `'/bin/zsh' -i -c 'echo hi; exec /bin/zsh -i -f'`).Return("", nil)
	mockTmux.On("SelectWindow", "demo:^").Return("", nil)
	mockLister.On("FindConfigSession", "demo").Return(model.SeshSession{}, false)
	mockLister.On("FindConfigWildcard", "/tmp").Return(model.WildcardConfig{}, false)

	msg, err := s.Exec(session)
	assert.Nil(t, err)
	assert.Equal(t, "", msg)
	mockTmux.AssertNotCalled(t, "NewWindow", "/tmp", "editor", `'/bin/zsh' -i -c 'echo hi; exec /bin/zsh -i -f'`)
}
