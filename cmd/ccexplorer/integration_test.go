package main

import (
	"bytes"
	"flag"
	"os"
	"path/filepath"
	"testing"

	"github.com/cduggn/ccexplorer/cmd/cli"
	"github.com/stretchr/testify/require"
)

var updateGolden = flag.Bool("update", false, "update golden files")

func TestCLIHelpGolden(t *testing.T) {
	t.Parallel()

	command := cli.RootCommand()
	var output bytes.Buffer
	command.SetOut(&output)
	command.SetErr(&output)
	command.SetArgs([]string{"--help"})

	err := command.Execute()
	require.NoError(t, err)

	goldenPath := filepath.Join("testdata", "help.golden")
	if *updateGolden {
		require.NoError(t, os.MkdirAll(filepath.Dir(goldenPath), 0o755))
		require.NoError(t, os.WriteFile(goldenPath, output.Bytes(), 0o644))
		return
	}

	expected, err := os.ReadFile(goldenPath)
	require.NoError(t, err)
	require.Equal(t, string(expected), output.String())
}
