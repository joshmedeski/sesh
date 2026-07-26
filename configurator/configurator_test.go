package configurator

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/joshmedeski/sesh/v2/model"
	"github.com/joshmedeski/sesh/v2/pathwrap"
	"github.com/joshmedeski/sesh/v2/runtimewrap"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// testOs implements oswrap.Os for testing
type testOs struct {
	homeDir     string
	homeDirErr  error
	configDir   string
	configErr   error
	files       map[string][]byte
	readFileErr map[string]error
	envVars     map[string]string
}

func (o *testOs) UserHomeDir() (string, error) {
	return o.homeDir, o.homeDirErr
}

func (o *testOs) UserConfigDir() (string, error) {
	return o.configDir, o.configErr
}

func (o *testOs) ReadFile(name string) ([]byte, error) {
	if o.readFileErr != nil {
		if err, ok := o.readFileErr[name]; ok {
			return nil, err
		}
	}
	if o.files != nil {
		if data, ok := o.files[name]; ok {
			return data, nil
		}
	}
	return nil, &os.PathError{Op: "open", Path: name, Err: os.ErrNotExist}
}

func (o *testOs) Getenv(key string) string {
	if o.envVars != nil {
		return o.envVars[key]
	}
	return ""
}

func (o *testOs) ExpandEnv(s string) string {
	return os.Expand(s, func(key string) string {
		if o.envVars != nil {
			return o.envVars[key]
		}
		return ""
	})
}

func (o *testOs) Stat(name string) (os.FileInfo, error) {
	return nil, nil
}

func (o *testOs) MkdirAll(path string, perm os.FileMode) error {
	return nil
}

func testdataPath(name string) string {
	abs, _ := filepath.Abs(filepath.Join("testdata", name))
	return abs
}

func TestGetConfig_DefaultPath(t *testing.T) {
	mockOs := &testOs{
		homeDir: "/home/testuser",
		files:   map[string][]byte{},
	}
	mockPath := pathwrap.NewPath()
	mockRuntime := &runtimewrap.MockRunTime{}

	c := NewConfigurator(mockOs, mockPath, mockRuntime)
	config, err := c.GetConfig()

	assert.NoError(t, err)
	assert.Equal(t, 1, config.DirLength) // default
}

func TestGetConfig_CustomPathValid(t *testing.T) {
	configFile := testdataPath("sesh.toml")
	data, err := os.ReadFile(configFile)
	require.NoError(t, err)

	mockOs := &testOs{
		homeDir: "/home/testuser",
		files: map[string][]byte{
			configFile: data,
		},
	}
	mockPath := pathwrap.NewPath()
	mockRuntime := &runtimewrap.MockRunTime{}

	c := NewConfiguratorWithPath(mockOs, mockPath, mockRuntime, configFile)
	config, err := c.GetConfig()

	assert.NoError(t, err)
	assert.Equal(t, "echo test", config.DefaultSessionConfig.StartupCommand)
	assert.Len(t, config.SessionConfigs, 1)
	assert.Equal(t, "test-session", config.SessionConfigs[0].Name)
	assert.Equal(t, "/tmp/test", config.SessionConfigs[0].Path)
}

func TestGetConfig_CustomPathNotFound(t *testing.T) {
	mockOs := &testOs{
		homeDir: "/home/testuser",
		files:   map[string][]byte{},
	}
	mockPath := pathwrap.NewPath()
	mockRuntime := &runtimewrap.MockRunTime{}

	c := NewConfiguratorWithPath(mockOs, mockPath, mockRuntime, "/nonexistent/sesh.toml")
	_, err := c.GetConfig()

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "couldn't read config file")
	assert.Contains(t, err.Error(), "/nonexistent/sesh.toml")
}

func TestGetConfig_CustomPathInvalidTOML(t *testing.T) {
	invalidFile := testdataPath("invalid.toml")
	data, err := os.ReadFile(invalidFile)
	require.NoError(t, err)

	mockOs := &testOs{
		homeDir: "/home/testuser",
		files: map[string][]byte{
			invalidFile: data,
		},
	}
	mockPath := pathwrap.NewPath()
	mockRuntime := &runtimewrap.MockRunTime{}

	c := NewConfiguratorWithPath(mockOs, mockPath, mockRuntime, invalidFile)
	_, err = c.GetConfig()

	assert.Error(t, err)
	var configErr *ConfigError
	assert.True(t, errors.As(err, &configErr))
}

func TestGetConfig_EmptyConfigPath(t *testing.T) {
	// Empty configPath should fall back to default behavior
	mockOs := &testOs{
		homeDir: "/home/testuser",
		files:   map[string][]byte{},
	}
	mockPath := pathwrap.NewPath()
	mockRuntime := &runtimewrap.MockRunTime{}

	c := NewConfiguratorWithPath(mockOs, mockPath, mockRuntime, "")
	config, err := c.GetConfig()

	assert.NoError(t, err)
	assert.Equal(t, 1, config.DirLength)
}

func TestGetConfig_XDGConfigHome(t *testing.T) {
	configFile := testdataPath("sesh.toml")
	data, err := os.ReadFile(configFile)
	require.NoError(t, err)

	mockOs := &testOs{
		homeDir: "/home/testuser",
		envVars: map[string]string{
			"XDG_CONFIG_HOME": "/custom/config",
		},
		files: map[string][]byte{
			"/custom/config/sesh/sesh.toml": data,
		},
	}
	mockPath := pathwrap.NewPath()
	mockRuntime := &runtimewrap.MockRunTime{}

	c := NewConfigurator(mockOs, mockPath, mockRuntime)
	config, err := c.GetConfig()

	assert.NoError(t, err)
	assert.Equal(t, "echo test", config.DefaultSessionConfig.StartupCommand)
	assert.Len(t, config.SessionConfigs, 1)
	assert.Equal(t, "test-session", config.SessionConfigs[0].Name)
}

func TestGetConfig_ImportPathWithEnvVar(t *testing.T) {
	// Import path uses $VAR syntax — should be env-expanded via oswrap
	importFile := testdataPath("sesh.toml")
	importData, err := os.ReadFile(importFile)
	require.NoError(t, err)

	mainTOML := []byte(`import = ["$CONFIGS/imported.toml"]` + "\n")

	mockOs := &testOs{
		homeDir: "/home/testuser",
		envVars: map[string]string{
			"CONFIGS": "/custom/dir",
		},
		files: map[string][]byte{
			"/main/sesh.toml":           mainTOML,
			"/custom/dir/imported.toml": importData,
		},
	}
	mockPath := pathwrap.NewPath()
	mockRuntime := &runtimewrap.MockRunTime{}

	c := NewConfiguratorWithPath(mockOs, mockPath, mockRuntime, "/main/sesh.toml")
	config, err := c.GetConfig()

	assert.NoError(t, err)
	// The imported config contributes "test-session"; its presence proves the env var
	// expanded to resolve /custom/dir/imported.toml
	assert.Len(t, config.SessionConfigs, 1)
	assert.Equal(t, "test-session", config.SessionConfigs[0].Name)
}

func TestGetConfig_ImportPathWithTilde(t *testing.T) {
	// Import path uses ~ syntax — resolved against UserHomeDir via c.path.Join
	importFile := testdataPath("sesh.toml")
	importData, err := os.ReadFile(importFile)
	require.NoError(t, err)

	mainTOML := []byte(`import = ["~/imports/imported.toml"]` + "\n")

	mockOs := &testOs{
		homeDir: "/home/testuser",
		files: map[string][]byte{
			"/main/sesh.toml":                      mainTOML,
			"/home/testuser/imports/imported.toml": importData,
		},
	}
	mockPath := pathwrap.NewPath()
	mockRuntime := &runtimewrap.MockRunTime{}

	c := NewConfiguratorWithPath(mockOs, mockPath, mockRuntime, "/main/sesh.toml")
	config, err := c.GetConfig()

	assert.NoError(t, err)
	assert.Len(t, config.SessionConfigs, 1)
	assert.Equal(t, "test-session", config.SessionConfigs[0].Name)
}

// configFromTOML loads a config straight from TOML content, going through the
// same defaults and validation as a real config file.
func configFromTOML(t *testing.T, contents string) (model.Config, error) {
	t.Helper()
	mockOs := &testOs{
		homeDir: "/home/testuser",
		files:   map[string][]byte{"/main/sesh.toml": []byte(contents)},
	}
	c := NewConfiguratorWithPath(mockOs, pathwrap.NewPath(), &runtimewrap.MockRunTime{}, "/main/sesh.toml")
	return c.GetConfig()
}

func TestGetConfig_AliasAutoConnectDelayDefault(t *testing.T) {
	config, err := configFromTOML(t, "")
	assert.NoError(t, err)
	assert.Equal(t, "150ms", config.TUI.AliasAutoConnectDelay)
}

func TestGetConfig_AliasAutoConnectDelayOverride(t *testing.T) {
	config, err := configFromTOML(t, "[tui]\nalias_auto_connect_delay = \"300ms\"\n")
	assert.NoError(t, err)
	assert.Equal(t, "300ms", config.TUI.AliasAutoConnectDelay)
}

func TestGetConfig_AliasAutoConnectDelayInvalid(t *testing.T) {
	_, err := configFromTOML(t, "[tui]\nalias_auto_connect_delay = \"soon\"\n")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid alias_auto_connect_delay")
}

func TestGetConfig_AliasFilterPrefixDefault(t *testing.T) {
	config, err := configFromTOML(t, "")
	assert.NoError(t, err)
	require.NotNil(t, config.TUI.AliasFilterPrefix, "an absent key is filled in with the default")
	assert.Equal(t, "/", *config.TUI.AliasFilterPrefix)
}

func TestGetConfig_AliasFilterPrefixOverride(t *testing.T) {
	config, err := configFromTOML(t, "[tui]\nalias_filter_prefix = \"@\"\n")
	assert.NoError(t, err)
	require.NotNil(t, config.TUI.AliasFilterPrefix)
	assert.Equal(t, "@", *config.TUI.AliasFilterPrefix)
}

func TestGetConfig_AliasFilterPrefixDisabled(t *testing.T) {
	config, err := configFromTOML(t, "[tui]\nalias_filter_prefix = \"\"\n")
	assert.NoError(t, err)
	require.NotNil(t, config.TUI.AliasFilterPrefix,
		"an explicit empty string disables the mode and must survive the defaults")
	assert.Equal(t, "", *config.TUI.AliasFilterPrefix)
}

func TestGetConfig_AliasFilterPrefixInvalid(t *testing.T) {
	_, err := configFromTOML(t, "[tui]\nalias_filter_prefix = \"//\"\n")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "must be a single character")

	_, err = configFromTOML(t, "[tui]\nalias_filter_prefix = \" \"\n")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "must not be whitespace")
}

func TestGetConfig_SessionAliases(t *testing.T) {
	config, err := configFromTOML(t, `
[[session]]
name = "wallpaper"
path = "~/c/wallpaper"
alias = "wp"
alias_auto_connect = true

[[session]]
name = "dotfiles"
path = "~/.config"
alias = "dot"
`)
	assert.NoError(t, err)
	require.Len(t, config.SessionConfigs, 2)
	assert.Equal(t, "wp", config.SessionConfigs[0].Alias)
	assert.True(t, config.SessionConfigs[0].AliasAutoConnect)
	assert.Equal(t, "dot", config.SessionConfigs[1].Alias)
	assert.False(t, config.SessionConfigs[1].AliasAutoConnect,
		"alias_auto_connect is opt-in per session")
}

func TestGetConfig_SessionIcons(t *testing.T) {
	config, err := configFromTOML(t, `
[[session]]
name = "notes"
path = "~/second-brain"
icon = "📓"

[[session]]
name = "dotfiles"
path = "~/.config"
`)
	assert.NoError(t, err)
	require.Len(t, config.SessionConfigs, 2)
	assert.Equal(t, "📓", config.SessionConfigs[0].Icon)
	assert.Equal(t, "", config.SessionConfigs[1].Icon,
		"a session without an icon keeps its source glyph")
}

func TestGetConfig_DuplicateAliases(t *testing.T) {
	_, err := configFromTOML(t, `
[[session]]
name = "wallpaper"
alias = "wp"

[[session]]
name = "wordpress"
alias = "WP"
`)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "duplicate alias")
	assert.Contains(t, err.Error(), "wordpress")
}

func TestGetConfig_OverlappingAliasPrefixesAllowed(t *testing.T) {
	// `w` and `wp` share a prefix on purpose: the auto-connect delay is what
	// makes them usable together.
	config, err := configFromTOML(t, `
[[session]]
name = "wallpaper"
alias = "wp"

[[session]]
name = "work"
alias = "w"
`)
	assert.NoError(t, err)
	assert.Len(t, config.SessionConfigs, 2)
}

func TestGetConfig_XDGConfigHomeNotSet(t *testing.T) {
	// When XDG_CONFIG_HOME is not set, should fall back to $HOME/.config
	configFile := testdataPath("sesh.toml")
	data, err := os.ReadFile(configFile)
	require.NoError(t, err)

	mockOs := &testOs{
		homeDir: "/home/testuser",
		// envVars not set, so XDG_CONFIG_HOME will return ""
		files: map[string][]byte{
			"/home/testuser/.config/sesh/sesh.toml": data,
		},
	}
	mockPath := pathwrap.NewPath()
	mockRuntime := &runtimewrap.MockRunTime{}

	c := NewConfigurator(mockOs, mockPath, mockRuntime)
	config, err := c.GetConfig()

	assert.NoError(t, err)
	assert.Equal(t, "echo test", config.DefaultSessionConfig.StartupCommand)
	assert.Len(t, config.SessionConfigs, 1)
	assert.Equal(t, "test-session", config.SessionConfigs[0].Name)
}
