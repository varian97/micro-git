package main

import (
	"fmt"
	"os"

	"micro-git/cmd"

	"github.com/akamensky/argparse"
)

type Pipeline struct {
	commands []cmd.ICommand
}

func (p *Pipeline) setNext(command cmd.ICommand) *Pipeline {
	p.commands = append(p.commands, command)
	return p
}

func (p Pipeline) execute() {
	if len(p.commands) == 0 {
		return
	}

	for _, command := range p.commands {
		if handled := command.Handle(); handled {
			break
		}
	}
}

func main() {
	pipeline := Pipeline{}
	parser := argparse.NewParser("mgit", "These are common microgit commands used in various situations.")

	// porcelain commands
	initCommand := &cmd.InitCommand{}
	commitCommand := &cmd.CommitCommand{}
	logCommand := &cmd.LogCommand{}

	// plumbing commands
	hashObjectCommand := &cmd.HashObjectCommand{}
	catCommand := &cmd.CatCommand{}
	writeTreeCommand := &cmd.WriteTreeCommand{}
	readTreeCommand := &cmd.ReadTreeCommand{}

	// Registration
	initCommand.Register(parser)
	commitCommand.Register(parser)
	logCommand.Register(parser)
	hashObjectCommand.Register(parser)
	catCommand.Register(parser)
	writeTreeCommand.Register(parser)
	readTreeCommand.Register(parser)

	pipeline.
		setNext(initCommand).
		setNext(commitCommand).
		setNext(logCommand).
		setNext(hashObjectCommand).
		setNext(catCommand).
		setNext(writeTreeCommand).
		setNext(readTreeCommand)

	err := parser.Parse(os.Args)
	if err != nil {
		fmt.Print(parser.Usage(err))
		return
	}

	pipeline.execute()
}
