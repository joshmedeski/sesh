package config_test

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

func prepareSeshConfig(t *testing.T) string {
	userConfigPath, err := os.MkdirTemp(os.TempDir(), "config")
	if err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(path.Join(userConfigPath, "sesh"), fs.ModePerm); err != nil {
		t.Fatal(err)
	}
	tempConfigPath := path.Join(userConfigPath, "sesh", "sesh.toml")
	secondTempConfigPath := path.Join(userConfigPath, "sesh", "sesh2.toml")

	err = os.WriteFile(tempConfigPath, []byte(fmt.Sprintf(`
		default_startup_script = "default"

		[[startup_scripts]]
		session_path = "~/dev/first_session"
		script_path = "~/.config/sesh/scripts/first_script"

		[[startup_scripts]]
		session_path = "~/dev/second_session"
		script_path = "~/.config/sesh/scripts/second_script"

		[[included_paths]]
		path = "%s"
		`, secondTempConfigPath),
	), fs.ModePerm)
	if err != nil {
		t.Fatal(err)
	}

	err = os.WriteFile(secondTempConfigPath, []byte(`
		[[startup_scripts]]
		session_path = "~/dev/third_session"
		script_path = "~/.config/sesh/scripts/third_script"
	`), fs.ModePerm)
	if err != nil {
		t.Fatal(err)
	}

	return userConfigPath
}

func TestParseConfigFile(t *testing.T) {
	t.Parallel()

	userConfigPath := prepareSeshConfig(t)
	defer os.Remove(userConfigPath)

	t.Run("ParseConfigFile", func(t *testing.T) {
		fetcher := &mockConfigDirectoryFetcher{dir: userConfigPath}
		config := config.ParseConfigFile(fetcher)

		if config.DefaultStartupScript != "default" {
			t.Errorf("Expected %s, got %s", "default", config.DefaultStartupScript)
		}

		if len(config.IncludedPaths) != 1 {
			t.Errorf("Expected %d, got %d", 1, len(config.IncludedPaths))
		}
		if config.IncludedPaths[0].Path != path.Join(userConfigPath, "sesh", "sesh2.toml") {
			t.Errorf("Expected %s, got %s", path.Join(userConfigPath, "sesh", "sesh2.toml"), config.IncludedPaths[0].Path)
		}

		if len(config.StartupScripts) != 3 {
			t.Errorf("Expected %d, got %d", 3, len(config.StartupScripts))
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
		if config.StartupScripts[2].SessionPath != "~/dev/third_session" {
			t.Errorf("Expected %s, got %s", "~/dev/third_session", config.StartupScripts[2].SessionPath)
		}
		if config.StartupScripts[2].ScriptPath != "~/.config/sesh/scripts/third_script" {
			t.Errorf("Expected %s, got %s", "~/.config/sesh/scripts/third_script", config.StartupScripts[2].ScriptPath)
		}
	})
}
