package writer

import (
	"os"
	"path/filepath"
)

// OutputDir is the directory generated CSV and chart files are written to,
// relative to the working directory.
var OutputDir = "./writer"

// ensureOutputDir creates the output directory on demand.
//
// This used to run in a package init(), which meant every invocation of the
// binary created ./writer in the caller's working directory - including
// `--help`, `--version`, and `mcp serve`, none of which write a file - and
// panicked outright if the directory could not be created. Creating it at
// render time keeps the filesystem untouched unless a file is actually
// produced, and turns a permission problem into a returned error.
func ensureOutputDir() error {
	return os.MkdirAll(OutputDir, 0o755)
}

// outputPath returns the full path for a generated file.
func outputPath(filename string) string {
	return filepath.Join(OutputDir, filename)
}
