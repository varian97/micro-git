package cat

import (
	"fmt"
	"strconv"

	"micro-git/cmd"
	"micro-git/object"

	"github.com/akamensky/argparse"
)

type CatCommand struct {
	cmd.CommandMeta
	fileInput            *string
	shouldShowObjectType *bool
	shouldShowSize       *bool
}

func (c *CatCommand) Register(p *argparse.Parser) {
	c.Command = p.NewCommand("cat-file", "Provide content or type and size information for repository objects")
	c.fileInput = c.Command.StringPositional(&argparse.Options{
		Required: true,
		Help:     "The name of the object to show",
	})
	c.shouldShowObjectType = c.Command.Flag("t", "type", &argparse.Options{
		Default: false,
		Help:    "Instead of the content, show the object type",
	})
	c.shouldShowSize = c.Command.Flag("s", "size", &argparse.Options{
		Default: false,
		Help:    "Instead of the content, show the object size",
	})
}

func (c *CatCommand) Handle() bool {
	if c.Command.Happened() {
		info, err := catFile(*c.fileInput, *c.shouldShowObjectType, *c.shouldShowSize)
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println(info)
		}
		return true
	}
	return false
}

func catFile(oid string, shouldShowType, shouldShowSize bool) (string, error) {
	objectInfo, err := object.Read(oid)
	if err != nil {
		return "", err
	}

	if shouldShowType && shouldShowSize {
		err = fmt.Errorf("-t and -s cannot be used altogether")
		return "", err
	}

	if shouldShowType {
		return objectInfo.Type, nil
	}
	if shouldShowSize {
		return strconv.Itoa(objectInfo.Size), nil
	}
	return string(objectInfo.Content), nil
}
