package main

import (
	"github.com/cduggn/ccexplorer/cmd/cli"
	"log/slog"
	"os"
)

// version is the build version. It is overwritten at release time by
// goreleaser via -X main.version and stays "dev" for local builds.
var version = "dev"

func main() {
	root := cli.RootCommand(version)

	if err := root.Execute(); err != nil {
		slog.Error("error", ErrAttr(err))
		os.Exit(126)
	}
}

func ErrAttr(err error) slog.Attr {
	return slog.Any("error", err)
}
