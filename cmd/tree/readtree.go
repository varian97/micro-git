package tree

import (
	"fmt"

	"micro-git/cmd"
	"micro-git/object"

	"github.com/akamensky/argparse"
)

type ReadTreeCommand struct {
	cmd.CommandMeta
	treeOid *string
}

func (c *ReadTreeCommand) Register(p *argparse.Parser) {
	c.Command = p.NewCommand("read-tree", "Reads tree information into the index")
	c.treeOid = c.Command.StringPositional(&argparse.Options{
		Required: true,
		Help:     "The id of the tree object to be read/merged",
	})
}

func (c *ReadTreeCommand) Handle() bool {
	if c.Command.Happened() {
		err := object.ReadTree(*c.treeOid)
		if err != nil {
			fmt.Println(err)
		}
		return true
	}
	return false
}
