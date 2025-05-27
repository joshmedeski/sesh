package namer

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConvertToValidName(t *testing.T) {
	t.Run("Test with dot", func(t *testing.T) {
		input := "test.name"
		want := "test_name"
		assert.Equal(t, want, convertToValidName(input))
	})

	t.Run("Test with colon", func(t *testing.T) {
		input := "test:name"
		want := "test_name"
		assert.Equal(t, want, convertToValidName(input))
	})

	t.Run("Test with space", func(t *testing.T) {
		input := "test name"
		want := "test_name"
		assert.Equal(t, want, convertToValidName(input))
	})

	t.Run("Test with multiple special characters", func(t *testing.T) {
		input := "test.name:with.multiple"
		want := "test_name_with_multiple"
		assert.Equal(t, want, convertToValidName(input))
	})

	t.Run("Test with spaces and other characters", func(t *testing.T) {
		input := "my project.folder:test name"
		want := "my_project_folder_test_name"
		assert.Equal(t, want, convertToValidName(input))
	})
}
