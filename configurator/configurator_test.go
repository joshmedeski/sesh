package configurator

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

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

func (o *testOs) Stat(name string) (os.FileInfo, error) {
	return nil, nil
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
