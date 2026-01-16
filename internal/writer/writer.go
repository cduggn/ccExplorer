package writer

import (
	"os"
)

// OutputDir is the default output directory for generated files
var OutputDir = "./writer"

func init() {
	if _, err := os.Stat(OutputDir); os.IsNotExist(err) {
		if err := os.Mkdir(OutputDir, 0755); err != nil {
			panic("unable to create writer directory: " + err.Error())
		}
	}
}
