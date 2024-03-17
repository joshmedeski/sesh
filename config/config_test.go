package config_test

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
		default_startup_script = "default"

		[[startup_scripts]]
		session_path = "~/dev/first_session"
		script_path = "~/.config/sesh/scripts/first_script"

		[[startup_scripts]]
		session_path = "~/dev/second_session"
		script_path = "~/.config/sesh/scripts/second_script"
		`), fs.ModePerm)
	if err != nil {
		t.Fatal(err)
	}

	return userConfigPath
}

func TestParseConfigFile(t *testing.T) {
	userConfigPath := prepareSeshConfig(t)
	defer os.Remove(userConfigPath)

	t.Run("ParseConfigFile", func(t *testing.T) {
		fetcher := &mockConfigDirectoryFetcher{dir: userConfigPath}
		config := config.ParseConfigFile(fetcher)

		if config.DefaultStartupScript != "default" {
			t.Errorf("Expected %s, got %s", "default", config.DefaultStartupScript)
		}
		if len(config.StartupScripts) != 2 {
			t.Errorf("Expected %d, got %d", 2, len(config.StartupScripts))
		}
		if config.StartupScripts[0].SessionPath != "~/dev/first_session" {
			t.Errorf("Expected %s, got %s", "~/dev/first_session", config.StartupScripts[0].SessionPath)
		}
		if config.StartupScripts[0].ScriptPath != "~/.config/sesh/scripts/first_script" {
			t.Errorf("Expected %s, got %s", "~/.config/sesh/scripts/first_script", config.StartupScripts[0].ScriptPath)
		}
		if config.StartupScripts[1].SessionPath != "~/dev/second_session" {
			t.Errorf("Expected %s, got %s", "~/dev/second_session", config.StartupScripts[1].SessionPath)
		}
		if config.StartupScripts[1].ScriptPath != "~/.config/sesh/scripts/second_script" {
			t.Errorf("Expected %s, got %s", "~/.config/sesh/scripts/second_script", config.StartupScripts[1].ScriptPath)
		}
	})
}
