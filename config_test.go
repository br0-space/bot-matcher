package matcher_test

import (
	"os"
	"path/filepath"
	"testing"

	matcher "github.com/br0-space/bot-matcher"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestLoadMatcherConfig_Success verifies that LoadMatcherConfig reads and unmarshals YAML into the target struct.
func TestLoadMatcherConfig_Success(t *testing.T) { //nolint:paralleltest
	// Prepare a temp working directory with config/tester.yml
	dir := t.TempDir()
	cfgDir := filepath.Join(dir, "config")
	require.NoError(t, os.MkdirAll(cfgDir, 0o755))

	content := []byte("name: Alice\nage: 42\n")
	require.NoError(t, os.WriteFile(filepath.Join(cfgDir, "tester.yml"), content, 0o600))

	// chdir into temp dir so LoadMatcherConfig finds config/tester.yml
	t.Chdir(dir)

	// target struct for unmarshal
	type target struct {
		Name string `mapstructure:"name"`
		Age  int    `mapstructure:"age"`
	}

	out, err := matcher.LoadMatcherConfig[target]("tester")
	require.NoError(t, err)
	assert.Equal(t, target{Name: "Alice", Age: 42}, out[0])
}

// TestLoadMatcherConfig_FileMissingReturnsError asserts that missing config files return an error.
func TestLoadMatcherConfig_FileMissingReturnsError(t *testing.T) { //nolint:paralleltest
	dir := t.TempDir()
	t.Chdir(dir)

	_, err := matcher.LoadMatcherConfig[struct{}]("does-not-exist")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to read fallback config")
}

// TestLoadMatcherConfig_InvalidYAML asserts that invalid YAML returns an error.
func TestLoadMatcherConfig_InvalidYAML(t *testing.T) { //nolint:paralleltest
	dir := t.TempDir()
	cfgDir := filepath.Join(dir, "config")
	require.NoError(t, os.MkdirAll(cfgDir, 0o755))

	// Write invalid YAML (mismatched indentation / bad syntax)
	invalidYAML := []byte("name: Alice\n  age: 42\n    invalid indentation")
	require.NoError(t, os.WriteFile(filepath.Join(cfgDir, "invalid.yml"), invalidYAML, 0o600))

	t.Chdir(dir)

	type target struct {
		Name string `mapstructure:"name"`
		Age  int    `mapstructure:"age"`
	}

	_, err := matcher.LoadMatcherConfig[target]("invalid")
	require.Error(t, err)
	// Invalid YAML can fail at parse or unmarshal stage
	assert.Contains(t, err.Error(), "failed to read fallback config")
}

// TestLoadMatcherConfig_PerChatConfigError verifies that errors in per-chat configs are returned.
func TestLoadMatcherConfig_PerChatConfigError(t *testing.T) { //nolint:paralleltest
	dir := t.TempDir()
	cfgDir := filepath.Join(dir, "config")
	require.NoError(t, os.MkdirAll(cfgDir, 0o755))

	// Write valid fallback config
	content := []byte("name: Default\nage: 0\n")
	require.NoError(t, os.WriteFile(filepath.Join(cfgDir, "perchat.yml"), content, 0o600))

	// Create per-chat directory and invalid per-chat config
	chatDir := filepath.Join(cfgDir, "123456")
	require.NoError(t, os.MkdirAll(chatDir, 0o755))

	invalidYAML := []byte("name: Chat\n  bad:\n    indentation")
	require.NoError(t, os.WriteFile(filepath.Join(chatDir, "perchat.yml"), invalidYAML, 0o600))

	t.Chdir(dir)

	type target struct {
		Name string `mapstructure:"name"`
		Age  int    `mapstructure:"age"`
	}

	_, err := matcher.LoadMatcherConfig[target]("perchat")
	require.Error(t, err)
	// Invalid YAML in per-chat config can fail at parse or unmarshal stage
	assert.Contains(t, err.Error(), "per-chat config")
}

// TestLoadMatcherConfig_WithPerChatConfigs verifies successful loading of both fallback and per-chat configs.
func TestLoadMatcherConfig_WithPerChatConfigs(t *testing.T) { //nolint:paralleltest
	dir := t.TempDir()
	cfgDir := filepath.Join(dir, "config")
	require.NoError(t, os.MkdirAll(cfgDir, 0o755))

	// Write fallback config
	fallback := []byte("name: Default\nage: 0\n")
	require.NoError(t, os.WriteFile(filepath.Join(cfgDir, "multichat.yml"), fallback, 0o600))

	// Create per-chat configs for multiple chat IDs
	for _, chatID := range []string{"111", "222", "333"} {
		chatDir := filepath.Join(cfgDir, chatID)
		require.NoError(t, os.MkdirAll(chatDir, 0o755))

		chatConfig := []byte("name: Chat" + chatID + "\nage: " + chatID + "\n")
		require.NoError(t, os.WriteFile(filepath.Join(chatDir, "multichat.yml"), chatConfig, 0o600))
	}

	t.Chdir(dir)

	type target struct {
		Name string `mapstructure:"name"`
		Age  int    `mapstructure:"age"`
	}

	out, err := matcher.LoadMatcherConfig[target]("multichat")
	require.NoError(t, err)

	// Check fallback config
	assert.Equal(t, target{Name: "Default", Age: 0}, out[0])

	// Check per-chat configs
	assert.Equal(t, target{Name: "Chat111", Age: 111}, out[111])
	assert.Equal(t, target{Name: "Chat222", Age: 222}, out[222])
	assert.Equal(t, target{Name: "Chat333", Age: 333}, out[333])
}
