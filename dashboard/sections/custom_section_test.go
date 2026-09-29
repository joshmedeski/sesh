package sections

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/joshmedeski/sesh/v2/dashboard/core"
	"github.com/joshmedeski/sesh/v2/model"
)

func TestCustomSectionCollapsesCarriageReturns(t *testing.T) {
	runner := core.NewMockCommandRunner(t)
	runner.EXPECT().RunShell("weather").Return([]byte("  0%\r 50%\r100%\r\nHouston: +87°F\r\n"), nil)
	s := NewCustomSection(model.DashboardSectionConfig{Custom: model.CustomConfig{Command: "weather"}}, core.SectionDeps{Runner: runner})

	msg, ok := s.Init()().(customOutputMsg)
	require.True(t, ok)
	assert.Equal(t, "100%\nHouston: +87°F\n", msg.output)
}
