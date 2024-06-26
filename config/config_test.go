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

func TestParseConfigFile(t *testing.T) {
	t.Parallel()
	home, err := os.UserHomeDir()
	if err != nil {
		t.Fatal("unable to get user's home directory")
	}

	mockSessions := []MockSeshConfig{
		{
			Name: "sesh.toml",
			Contents: fmt.Sprintf(`
			[default_session]
			startup_script = "default"

			[[session]]
			path = "~/dev/first_session"
			startup_script = "~/.config/sesh/scripts/first_script"

			[[session]]
			path = "%s/dev/second_session"
			startup_script = "%s/.config/sesh/scripts/second_script"
			`, home, home),
			Imports: "sesh2.toml",
		},
		{
			Name: "sesh2.toml",
			Contents: `
			[[session]]
			path = "~/dev/third_session"
			startup_script = "~/.config/sesh/scripts/third_script"
			`,
		},
	}

	userConfigPath := prepareSeshConfig(t, mockSessions)
	defer os.Remove(userConfigPath)

	t.Run("ParseConfigFile", func(t *testing.T) {
		fetcher := &mockConfigDirectoryFetcher{dir: userConfigPath}
		config := config.ParseConfigFile(fetcher)

		if config.DefaultSessionConfig.StartupScript != "default" {
			t.Errorf("Expected %s, got %s", "default", config.DefaultSessionConfig.StartupScript)
		}

		if len(config.ImportPaths) != 1 {
			t.Errorf("Expected %d, got %d", 1, len(config.ImportPaths))
		}
		if config.ImportPaths[0] != path.Join(userConfigPath, "sesh", "sesh2.toml") {
			t.Errorf("Expected %s, got %s", path.Join(userConfigPath, "sesh", "sesh2.toml"), config.ImportPaths[0])
		}

		if len(config.SessionConfigs) != 3 {
			t.Errorf("Expected %d, got %d", 3, len(config.SessionConfigs))
		}
		if config.SessionConfigs[0].Path != (home + "/dev/first_session") {
			t.Errorf("Expected %s, got %s", home+"/dev/first_session", config.SessionConfigs[0].Path)
		}
		if config.SessionConfigs[0].StartupScript != (home + "/.config/sesh/scripts/first_script") {
			t.Errorf("Expected %s, got %s", home+"/.config/sesh/scripts/first_script", config.SessionConfigs[0].StartupScript)
		}
		if config.SessionConfigs[1].Path != (home + "/dev/second_session") {
			t.Errorf("Expected %s, got %s", home+"/dev/second_session", config.SessionConfigs[1].Path)
		}
		if config.SessionConfigs[1].StartupScript != (home + "/.config/sesh/scripts/second_script") {
			t.Errorf("Expected %s, got %s", home+"/.config/sesh/scripts/second_script", config.SessionConfigs[1].StartupScript)
		}
		if config.SessionConfigs[2].Path != (home + "/dev/third_session") {
			t.Errorf("Expected %s, got %s", home+"/dev/third_session", config.SessionConfigs[2].Path)
		}
		if config.SessionConfigs[2].StartupScript != (home + "/.config/sesh/scripts/third_script") {
			t.Errorf("Expected %s, got %s", home+"/.config/sesh/scripts/third_script", config.SessionConfigs[2].StartupScript)
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

	importPathsStringBuilder := strings.Builder{}
	importPathsStringBuilder.WriteString("import = [")
	importPaths := make([]string, extended_configs_count)
	for i := 0; i < extended_configs_count; i++ {
		configPath := path.Join(userConfigPath, "sesh", fmt.Sprintf("sesh%d.toml", i))
		importPaths[i] = configPath
		importPathsStringBuilder.WriteString(fmt.Sprintf(`"%s",`, configPath))
	}
	importPathsStringBuilder.WriteString("]")

	err = os.WriteFile(tempConfigPath, []byte(fmt.Sprintf(`
		%s
		default_startup_script = "default"

		[[startup_scripts]]
		session_path = "~/dev/first_session"
		script_path = "~/.config/sesh/scripts/first_script"

		[[startup_scripts]]
		session_path = "~/dev/second_session"
		script_path = "~/.config/sesh/scripts/second_script"	
		`, importPathsStringBuilder.String()),
	), fs.ModePerm)
	if err != nil {
		b.Fatal(err)
	}

	for i, configPath := range importPaths {
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
	b.Skip("Skipping benchmark because it will be failing on CI")
	table := []struct {
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

func TestConfigPathsAreCleaned(t *testing.T) {
	t.Parallel()

	home, err := os.UserHomeDir()
	if err != nil {
		t.Fatalf("failed to get home dir: %s", err.Error())
	}

	mockSessions := []MockSeshConfig{
		{
			Name: "sesh.toml",
			Contents: fmt.Sprintf(`
			[default_session]
			startup_script = "default"

			[[session]]
			path = "~/tilde-prefixed"
			name = "my-tilde-prefixed-session"
			startup_script = "~/my_script"

			[[session]]
			path = "%s/tilde-prefixed"
			name = "my-tilde-prefixed-session"
			startup_script = "%s/my_script"
			`, home, home),
		},
	}

	userConfigPath := prepareSeshConfig(t, mockSessions)
	defer os.Remove(userConfigPath)

	fetcher := &mockConfigDirectoryFetcher{dir: userConfigPath}
	cfg := config.ParseConfigFile(fetcher)

	err = cfg.CleanPaths()
	if err != nil {
		t.Fatalf("failed to clean paths: %s", err.Error())
	}

	// test all config sessions
	for _, sessionConfig := range cfg.SessionConfigs {
		if !strings.HasPrefix(sessionConfig.Path, home) {
			t.Fatalf("config session path is not absolute")
		}

		if !strings.HasPrefix(sessionConfig.StartupScript, home) {
			t.Fatalf("config session startup script is not absolute")
		}
	}
}
