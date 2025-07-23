package commit

import (
	"fmt"

	"micro-git/cmd"
	"micro-git/object"

	"github.com/akamensky/argparse"
)

type LogCommand struct {
	cmd.CommandMeta
}

func (c *LogCommand) Register(p *argparse.Parser) {
	c.Command = p.NewCommand("log", "Shows the commit logs.")
}

func (c *LogCommand) Handle() bool {
	if c.Command.Happened() {
		err := object.PrintCommitLogs()
		if err != nil {
			fmt.Println(err)
		}

		return true
	}
	return false
}
