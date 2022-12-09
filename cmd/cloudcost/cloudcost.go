package main

import (
	"github.com/cduggn/cloudcost/internal/pkg/cmd"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	cmd.Execute()
}
