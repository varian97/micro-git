package main

import (
	"fmt"
	"os"

	"micro-git/filesystem"

	"github.com/akamensky/argparse"
)

func main() {
	parser := argparse.NewParser("microgit", "These are common microgit commands used in various situations.")

	initCommand := parser.NewCommand("init", "Create an empty micro-git repository or reinitialize an existing one")

	err := parser.Parse(os.Args)
	if err != nil {
		fmt.Print(parser.Usage(err))
		return
	}

	if initCommand.Happened() {
		filesystem.Init()
	}
}
