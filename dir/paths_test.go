package dir

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAlternatePath(t *testing.T) {
	t.Run("absolute path", func(t *testing.T) {
		require.Equal(t, "/foo/bar", AlternatePath("/foo/bar"))
	})
	t.Run("home directory", func(t *testing.T) {
		homeDir, err := os.UserHomeDir()
		require.NoError(t, err)
		require.Equal(t, homeDir+"/foo/bar", AlternatePath("~/foo/bar"))
	})
	t.Run("relative path", func(t *testing.T) {
		wd, err := os.Getwd()
		require.NoError(t, err)
		require.Equal(t, wd+"/foo/bar", AlternatePath("./foo/bar"))
	})
}
