package main

import (
	"github.com/cduggn/ccexplorer/internal/commands"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	commands.Execute()
}
