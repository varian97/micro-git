package init

import (
	"fmt"

	"micro-git/cmd"
	"micro-git/root"

	"github.com/akamensky/argparse"
)

type InitCommand struct {
	cmd.CommandMeta
}

func (c *InitCommand) Register(p *argparse.Parser) {
	c.Command = p.NewCommand("init", "Create an empty micro-git repository or reinitialize an existing one")
}

func (c *InitCommand) Handle() bool {
	if c.Command.Happened() {
		err := root.InitDB()
		if err != nil {
			fmt.Println(err)
		}
		return true
	}
	return false
}
