package seshcli

import (
	"testing"

	"github.com/joshmedeski/sesh/v2/model"
	"github.com/joshmedeski/sesh/v2/refresher"
)

func TestMaybeWarmStatus(t *testing.T) {
	t.Run("spawns when caching enabled", func(t *testing.T) {
		rf := new(refresher.MockRefresher)
		rf.On("Spawn", "").Return(nil)
		deps := &Deps{}
		deps.Refresher = rf
		deps.Config = model.Config{} // IssueTTL nil → EffectiveTTL 60

		maybeWarmStatus(deps)

		rf.AssertExpectations(t)
	})

	t.Run("does not spawn when issue_ttl is zero", func(t *testing.T) {
		rf := new(refresher.MockRefresher)
		deps := &Deps{}
		deps.Refresher = rf
		zero := 0
		deps.Config = model.Config{Github: model.GithubConfig{IssueTTL: &zero}}

		maybeWarmStatus(deps)

		rf.AssertNotCalled(t, "Spawn", "")
	})
}
