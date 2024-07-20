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

	hashObjectCommand := parser.NewCommand("hash-object", "Compute object ID and optionally creates a blob from a file")
	hashObjectFileType := hashObjectCommand.String("t", "type", &argparse.Options{
		Help:    "Specify the type (default: \"blob\")",
		Default: "blob",
	})
	hashObjectWriteFlag := hashObjectCommand.Flag("w", "write", &argparse.Options{
		Required: false,
		Help:     "Actually write the object into the object database",
	})
	hashObjectFileInput := hashObjectCommand.StringPositional(&argparse.Options{
		Required: true,
		Help:     "The file to be hashed",
	})

	err := parser.Parse(os.Args)
	if err != nil {
		fmt.Print(parser.Usage(err))
		return
	}

	if initCommand.Happened() {
		err := filesystem.Init()
		if err != nil {
			panic(err)
		}
	}

	if hashObjectCommand.Happened() {
		hexSum, err := filesystem.WriteObject(*hashObjectFileInput, *hashObjectFileType, *hashObjectWriteFlag)
		if err != nil {
			panic(err)
		}
		fmt.Println(hexSum)
	}
}
