package cmd

import (
	"fmt"
	"os"

	"micro-git/object"

	"github.com/akamensky/argparse"
)

type HashObjectCommand struct {
	CommandMeta
	fileInput   *string
	objectType  *string
	shouldWrite *bool
}

func (c *HashObjectCommand) Register(p *argparse.Parser) {
	c.Command = p.NewCommand("hash-object", "Compute object ID and optionally creates a blob from a file")
	c.objectType = c.Command.String("t", "type", &argparse.Options{
		Help:    "Specify the type (default: \"blob\")",
		Default: "blob",
	})
	c.shouldWrite = c.Command.Flag("w", "write", &argparse.Options{
		Required: false,
		Help:     "Actually write the object into the object database",
	})
	c.fileInput = c.Command.StringPositional(&argparse.Options{
		Required: true,
		Help:     "The file to be hashed",
	})
}

func (c *HashObjectCommand) Handle() bool {
	if c.Command.Happened() {
		hexSum, err := hashObject(*c.fileInput, *c.objectType, *c.shouldWrite)
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println(hexSum)
		}
		return true
	}
	return false
}

func hashObject(path, objectType string, shouldWrite bool) (string, error) {
	fileContent, err := os.ReadFile(path)
	if err != nil {
		err := fmt.Errorf("failed to read file content, %v", err)
		return "", err
	}

	if shouldWrite {
		return object.Write(objectType, fileContent)
	}

	objectInfo, err := object.GenInfo(objectType, fileContent)
	if err != nil {
		return "", err
	}

	return objectInfo.Oid, nil
}
