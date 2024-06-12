package session

import (
	"io/fs"
	"os"
	"path"
	"testing"

	"github.com/joshmedeski/sesh/config"
)

type mockConfigDirectoryFetcher struct {
	dir string
}

func (m *mockConfigDirectoryFetcher) GetUserConfigDir() (string, error) {
	return m.dir, nil
}

func prepareSeshConfig(t *testing.T) string {
	userConfigPath, err := os.MkdirTemp(os.TempDir(), "config")
	if err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(path.Join(userConfigPath, "sesh"), fs.ModePerm); err != nil {
		t.Fatal(err)
	}
	tempConfigPath := path.Join(userConfigPath, "sesh", "sesh.toml")

	err = os.WriteFile(tempConfigPath, []byte(`
		[[session]]
		name = "test-session"
		path = "~/dev/duplicate_session"
		`,
	), fs.ModePerm)
	if err != nil {
		t.Fatal(err)
	}

	return userConfigPath
}

func mockCfg(t *testing.T) config.Config {
	cfgDir := prepareSeshConfig(t)
	fetcher := &mockConfigDirectoryFetcher{dir: cfgDir}
	defer os.Remove(cfgDir)
	cfg := config.ParseConfigFile(fetcher)
	return cfg
}

// Verify that custom sessions defined in config with a path
// equal to a tmux session is still added to the session list
func TestListConfigSessionsShouldAddConfigSessionWithDuplicatePath(t *testing.T) {
	cfg := mockCfg(t)
	sessions := []Session{
		{Src: "tmux", Name: "some-other-name", Path: "~/dev/duplicate_session"},
	}
	configSessions, err := listConfigSessions(&cfg, sessions)
	if err != nil {
		t.Fatalf("error: %v", err.Error())
	}

	if len(configSessions) != 2 {
		t.Fatal("error: Expected config session with duplicate path to be added to tmux session")
	}
}
