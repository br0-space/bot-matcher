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

	out := matcher.LoadMatcherConfig[target]("tester")
	assert.Equal(t, target{Name: "Alice", Age: 42}, out[0])
}

// TestLoadMatcherConfig_FileMissingPanics asserts that missing config files cause a panic.
func TestLoadMatcherConfig_FileMissingPanics(t *testing.T) { //nolint:paralleltest
	dir := t.TempDir()
	t.Chdir(dir)

	require.Panics(t, func() { _ = matcher.LoadMatcherConfig[struct{}]("does-not-exist") })
}
