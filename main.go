package main

import (
	"fmt"
	"os"

	"micro-git/core"

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

	catFileCommand := parser.NewCommand("cat-file", "Provide content or type and size information for repository objects")
	catFileInput := catFileCommand.StringPositional(&argparse.Options{
		Required: true,
		Help:     "The name of the object to show",
	})
	catFileShouldShowObjectType := catFileCommand.Flag("t", "type", &argparse.Options{
		Default: false,
		Help:    "Instead of the content, show the object type",
	})
	catFileShouldShowSize := catFileCommand.Flag("s", "size", &argparse.Options{
		Default: false,
		Help:    "Instead of the content, show the object size",
	})

	err := parser.Parse(os.Args)
	if err != nil {
		fmt.Print(parser.Usage(err))
		return
	}

	if initCommand.Happened() {
		err := core.Init()
		if err != nil {
			fmt.Println(err)
			return
		}
	}

	if hashObjectCommand.Happened() {
		hexSum, err := core.HashObject(*hashObjectFileInput, *hashObjectFileType, *hashObjectWriteFlag)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println(hexSum)
	}

	if catFileCommand.Happened() {
		objectInfo, err := core.CatFile(*catFileInput)
		if err != nil {
			fmt.Println(err)
			return
		}

		if *catFileShouldShowObjectType && *catFileShouldShowSize {
			fmt.Println("-t and -s cannot be used altogether")
			return
		}

		if *catFileShouldShowObjectType {
			fmt.Println(objectInfo.Type)
		} else if *catFileShouldShowSize {
			fmt.Println(objectInfo.Size)
		} else {
			fmt.Println(string(objectInfo.Content))
		}
	}
}
