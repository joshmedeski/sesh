package config_test

import (
	"fmt"
	"io/fs"
	"os"
	"path"
	"strings"
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

		[[extended_configs]]
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

		if len(config.ExtendedConfigs) != 1 {
			t.Errorf("Expected %d, got %d", 1, len(config.ExtendedConfigs))
		}
		if config.ExtendedConfigs[0].Path != path.Join(userConfigPath, "sesh", "sesh2.toml") {
			t.Errorf("Expected %s, got %s", path.Join(userConfigPath, "sesh", "sesh2.toml"), config.ExtendedConfigs[0].Path)
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

func prepareSeshConfigForBench(b *testing.B, extended_configs_count int) string {
	userConfigPath, err := os.MkdirTemp(os.TempDir(), "config")
	if err != nil {
		b.Fatal(err)
	}
	if err := os.MkdirAll(path.Join(userConfigPath, "sesh"), fs.ModePerm); err != nil {
		b.Fatal(err)
	}
	tempConfigPath := path.Join(userConfigPath, "sesh", "sesh.toml")

	extendedConfigsStringBuilder := strings.Builder{}
	extendedConfigs := make([]string, extended_configs_count)
	for i := 0; i < extended_configs_count; i++ {
		configPath := path.Join(userConfigPath, "sesh", fmt.Sprintf("sesh%d.toml", i))
		extendedConfigs[i] = configPath
		extendedConfigsStringBuilder.WriteString(fmt.Sprintf(`
		[[extended_configs]]
		path = "%s"
		`, configPath))
	}

	err = os.WriteFile(tempConfigPath, []byte(fmt.Sprintf(`
		default_startup_script = "default"

		[[startup_scripts]]
		session_path = "~/dev/first_session"
		script_path = "~/.config/sesh/scripts/first_script"

		[[startup_scripts]]
		session_path = "~/dev/second_session"
		script_path = "~/.config/sesh/scripts/second_script"

		%s
		`, extendedConfigsStringBuilder.String()),
	), fs.ModePerm)
	if err != nil {
		b.Fatal(err)
	}

	for i, configPath := range extendedConfigs {
		err = os.WriteFile(configPath, []byte(fmt.Sprintf(`
		[[startup_scripts]]
		session_path = "~/dev/session_%d"
		script_path = "~/.config/sesh/scripts/script"
		`, i),
		), fs.ModePerm)
		if err != nil {
			b.Fatal(err)
		}
	}

	return userConfigPath
}

func BenchmarkParseConfigFile(b *testing.B) {
	var table = []struct {
		input int
	}{
		{input: 1},
		{input: 10},
		{input: 100},
		{input: 1000},
		{input: 10000},
	}

	for _, test := range table {

		b.Run(fmt.Sprintf("ParseConfigFile_%d", test.input), func(b *testing.B) {
			userConfigPath := prepareSeshConfigForBench(b, test.input)
			defer os.Remove(userConfigPath)

			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				fetcher := &mockConfigDirectoryFetcher{dir: userConfigPath}
				config.ParseConfigFile(fetcher)
			}
		})
	}
}
