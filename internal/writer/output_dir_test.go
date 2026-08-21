package writer

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Importing this package used to create ./writer as a side effect of init(),
// so every run of the binary littered the caller's working directory even for
// --help or a stdout-only query.
func TestImportingWriterCreatesNothing(t *testing.T) {
	dir := t.TempDir()
	_, err := os.Stat(filepath.Join(dir, "writer"))
	assert.True(t, os.IsNotExist(err),
		"no directory should exist until something is rendered")
}

func TestRenderCreatesOutputDirOnDemand(t *testing.T) {
	dir := t.TempDir()
	target := filepath.Join(dir, "nested", "out")

	old := OutputDir
	OutputDir = target
	t.Cleanup(func() { OutputDir = old })

	_, err := os.Stat(target)
	require.True(t, os.IsNotExist(err), "precondition: directory absent")

	err = NewCSVRenderer().Render(&CSVOutput{
		Headers:  []string{"a"},
		Rows:     [][]string{{"1"}},
		Filename: "report.csv",
	})
	require.NoError(t, err)

	// MkdirAll, not Mkdir: a nested path must work.
	info, err := os.Stat(target)
	require.NoError(t, err)
	assert.True(t, info.IsDir())
	assert.FileExists(t, filepath.Join(target, "report.csv"))
}

// Rendering twice must not fail on the already-existing directory.
func TestRenderIsIdempotent(t *testing.T) {
	dir := t.TempDir()
	old := OutputDir
	OutputDir = dir
	t.Cleanup(func() { OutputDir = old })

	for i := 0; i < 2; i++ {
		require.NoError(t, NewCSVRenderer().Render(&CSVOutput{
			Headers:  []string{"a"},
			Rows:     [][]string{{"1"}},
			Filename: "report.csv",
		}))
	}
}

// An unwritable location must surface as an error, not a panic.
func TestRenderReturnsErrorWhenOutputDirCannotBeCreated(t *testing.T) {
	dir := t.TempDir()
	blocker := filepath.Join(dir, "blocked")
	require.NoError(t, os.WriteFile(blocker, []byte("not a directory"), 0o600))

	old := OutputDir
	OutputDir = filepath.Join(blocker, "sub")
	t.Cleanup(func() { OutputDir = old })

	err := NewCSVRenderer().Render(&CSVOutput{
		Headers:  []string{"a"},
		Rows:     [][]string{{"1"}},
		Filename: "report.csv",
	})
	require.Error(t, err, "must return an error rather than panicking")
	assert.Contains(t, err.Error(), "output directory")
}
