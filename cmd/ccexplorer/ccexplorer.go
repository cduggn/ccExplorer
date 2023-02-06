package main

import (
	"fmt"
	"github.com/cduggn/ccexplorer/internal/commands"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	commands.Execute()

	// testing github secrete scanner
	aws_access_key_id := "YIYIOUY£465465465hjkhjk"
	aws_secret_access_key := "YIYIOUY£465465465098908908"
	fmt.Print(aws_access_key_id, aws_secret_access_key)
}
