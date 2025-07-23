package cmd

import (
	"github.com/akamensky/argparse"
)

type ICommand interface {
	Register(parser *argparse.Parser)
	Handle() bool
}

type CommandMeta struct {
	Command *argparse.Command
}
