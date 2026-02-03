package configurator

import (
	"testing"

	"github.com/joshmedeski/sesh/v2/model"
	"github.com/pelletier/go-toml/v2"
	"github.com/stretchr/testify/assert"
)

func TestWildcardConfigParsing(t *testing.T) {
	t.Run("should parse wildcard configs from TOML", func(t *testing.T) {
		input := []byte(`
[[wildcard]]
pattern = "~/projects/*"
startup_command = "nvim"

[[wildcard]]
pattern = "~/work/*"
startup_command = "make dev"
preview_command = "ls -la"
disable_startup_command = true
`)
		config := model.Config{}
		err := toml.Unmarshal(input, &config)
		assert.Nil(t, err)
		assert.Len(t, config.WildcardConfigs, 2)
		assert.Equal(t, "~/projects/*", config.WildcardConfigs[0].Pattern)
		assert.Equal(t, "nvim", config.WildcardConfigs[0].StartupCommand)
		assert.Equal(t, "~/work/*", config.WildcardConfigs[1].Pattern)
		assert.Equal(t, "make dev", config.WildcardConfigs[1].StartupCommand)
		assert.Equal(t, "ls -la", config.WildcardConfigs[1].PreviewCommand)
		assert.True(t, config.WildcardConfigs[1].DisableStartCommand)
	})
}

func TestWildcardImportMerge(t *testing.T) {
	t.Run("should merge wildcard configs from imports", func(t *testing.T) {
		mainConfig := model.Config{
			WildcardConfigs: []model.WildcardConfig{
				{Pattern: "~/projects/*", StartupCommand: "nvim"},
			},
		}
		importConfig := model.Config{
			WildcardConfigs: []model.WildcardConfig{
				{Pattern: "~/imported/*", StartupCommand: "echo imported"},
			},
		}
		mainConfig.WildcardConfigs = append(mainConfig.WildcardConfigs, importConfig.WildcardConfigs...)
		assert.Len(t, mainConfig.WildcardConfigs, 2)
		assert.Equal(t, "~/imported/*", mainConfig.WildcardConfigs[1].Pattern)
	})
}
