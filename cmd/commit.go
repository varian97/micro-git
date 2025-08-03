package cmd

import (
	"fmt"

	"micro-git/object"

	"github.com/akamensky/argparse"
)

type CommitCommand struct {
	CommandMeta
	commitMessage *string
}

func (c *CommitCommand) Register(p *argparse.Parser) {
	c.Command = p.NewCommand("commit", "Record changes to the repository")
	c.commitMessage = c.Command.String("m", "message", &argparse.Options{
		Help:     "The commit message",
		Required: true,
	})
}

func (c *CommitCommand) Handle() bool {
	if c.Command.Happened() {
		commitOid, err := object.Commit(*c.commitMessage)
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println(commitOid)
		}

		return true
	}
	return false
}
