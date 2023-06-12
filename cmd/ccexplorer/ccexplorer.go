package main

import (
	"github.com/cduggn/ccexplorer/internal/core/handlers/commandline"
	"github.com/cduggn/ccexplorer/internal/core/logger"
	"os"
)

func main() {
	root := commandline.RootCommand()

	if err := root.Execute(); err != nil {
		logger.ErrorOut(err)
		os.Exit(126)
	}

}
