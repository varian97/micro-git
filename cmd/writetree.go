package cmd

import (
	"fmt"

	"micro-git/object"

	"github.com/akamensky/argparse"
)

type WriteTreeCommand struct {
	CommandMeta
	subdir *string
}

func (c *WriteTreeCommand) Register(p *argparse.Parser) {
	c.Command = p.NewCommand("write-tree", "Create tree object from the current index")
	c.subdir = c.Command.String("", "prefix", &argparse.Options{
		Required: false,
		Help:     "Write a tree object from subdirectory <prefix>",
		Default:  ".",
	})
}

func (c *WriteTreeCommand) Handle() bool {
	if c.Command.Happened() {
		treeId, err := object.WriteTree(*c.subdir)
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println(treeId)
		}

		return true
	}
	return false
}
