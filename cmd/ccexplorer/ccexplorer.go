package main

import (
	"fmt"
	"github.com/cduggn/ccexplorer/internal/core/handlers/commandline"
	"os"
)

func main() {
	root := commandline.RootCommand()

	if err := root.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(126)
	}

}
