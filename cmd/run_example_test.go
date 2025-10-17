package main

import (
	"os"
	"testing"

	"github.com/spf13/pflag"
	"github.com/stretchr/testify/require"
)

// TestMain_ExecutesWithoutPanic ensures the example main() function runs end-to-end
// without panicking. It isolates pflag.CommandLine and os.Args to avoid interference
// with other tests in the suite.
func TestMain_ExecutesWithoutPanic(t *testing.T) {
	t.Parallel()

	// Save and replace pflag.CommandLine to avoid global flag pollution.
	savedFlagSet := pflag.CommandLine

	defer func() { pflag.CommandLine = savedFlagSet }()

	pflag.CommandLine = pflag.NewFlagSet("run_example_test", pflag.ContinueOnError)

	// Save and restore os.Args so pflag.Parse() sees a clean argument list.
	savedArgs := os.Args

	defer func() { os.Args = savedArgs }()

	os.Args = []string{"run_example"}

	// Ensure working directory is project root so config files are resolvable.
	cwd, err := os.Getwd()
	require.NoError(t, err)

	defer func() { _ = os.Chdir(cwd) }() //nolint:usetesting

	require.NoError(t, os.Chdir("..")) //nolint:usetesting

	require.NotPanics(t, func() { main() }, "main() should not panic")
}
