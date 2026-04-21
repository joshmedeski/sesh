package namer

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConvertToValidName(t *testing.T) {
	t.Run("tmux backend with dot", func(t *testing.T) {
		input := "test.name"
		want := "test_name"
		assert.Equal(t, want, convertToValidName(input, ""))
	})

	t.Run("tmux backend with colon", func(t *testing.T) {
		input := "test:name"
		want := "test_name"
		assert.Equal(t, want, convertToValidName(input, ""))
	})

	t.Run("tmux backend with multiple special characters", func(t *testing.T) {
		input := "test.name:with.multiple"
		want := "test_name_with_multiple"
		assert.Equal(t, want, convertToValidName(input, ""))
	})

	t.Run("wezterm backend preserves dots", func(t *testing.T) {
		input := "test.name"
		want := "test.name"
		assert.Equal(t, want, convertToValidName(input, "wezterm"))
	})

	t.Run("wezterm backend preserves colons", func(t *testing.T) {
		input := "test:name"
		want := "test:name"
		assert.Equal(t, want, convertToValidName(input, "wezterm"))
	})

	t.Run("wezterm backend preserves multiple special characters", func(t *testing.T) {
		input := "test.name:with.multiple"
		want := "test.name:with.multiple"
		assert.Equal(t, want, convertToValidName(input, "wezterm"))
	})
}
