package session

import (
	"fmt"
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

type MockSeshConfig struct {
	Name     string
	Contents string
	Imports  string
}

func prepareSeshConfig(t *testing.T, mockSeshConfigs []MockSeshConfig) string {
	userConfigPath, err := os.MkdirTemp(os.TempDir(), "config")
	if err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(path.Join(userConfigPath, "sesh"), fs.ModePerm); err != nil {
		t.Fatal(err)
	}

	// create a temp config file for each supplied config
	for _, mockSesh := range mockSeshConfigs {
		tempConfigPath := path.Join(userConfigPath, "sesh", mockSesh.Name)

		var contents string
		if mockSesh.Imports != "" {
			contents = fmt.Sprintf("import = [\"%s\"]\n%s", path.Join(userConfigPath, "sesh", mockSesh.Imports), mockSesh.Contents)
		} else {
			contents = mockSesh.Contents
		}
		err = os.WriteFile(tempConfigPath, []byte(contents), fs.ModePerm)

		if err != nil {
			t.Fatal(err)
		}
	}

	return userConfigPath
}

// Verify that custom sessions defined in config with a path
// equal to a tmux session is still added to the session list.
//
// Only tmux sessions created from custom config sessions directly
// should overwrite config sessions
func TestListConfigSessionsShouldAddConfigSessionWithDuplicatePath(t *testing.T) {
	mockCfg := []MockSeshConfig{
		{
			Name: "sesh.toml",
			Contents: `
			[[session]]
			name = "test-session"
			path = "~/dev/duplicate_session"
		
			[[session]] 
			name = "test-session-2"
			path = "~/dev/duplicate_session"
		`},
	}
	cfgDir := prepareSeshConfig(t, mockCfg)
	defer os.Remove(cfgDir)

	fetcher := &mockConfigDirectoryFetcher{dir: cfgDir}
	cfg := config.ParseConfigFile(fetcher)

	home, err := os.UserHomeDir()
	if err != nil {
		t.Fatal("unable to get user home directory")
	}

	sessions := []Session{
		// session has been created from config session
		// include it, but not the one from the config
		{Src: "tmux", Name: "test-session", Path: home + "/dev/duplicate_session"},
	}

	configSessions, err := listConfigSessions(&cfg, sessions)
	if err != nil {
		t.Fatalf("error: %v", err.Error())
	}

	if len(configSessions) != 1 {
		t.Fatalf("Expected created configSessions to be %d, got %d", 1, len(configSessions))
	}

	if configSessions[0].Name != "test-session-2" {
		t.Fatalf("Expected created configSession to be %s, got %s", "test-session-2", configSessions[0].Name)
	}
}
