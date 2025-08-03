package cmd

import (
	"fmt"
	"os/user"
	"time"

	"micro-git/object"
	"micro-git/refs"

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
		commitOid, err := commit(*c.commitMessage)
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println(commitOid)
		}

		return true
	}
	return false
}

func commit(msg string) (string, error) {
	/*
		Create commit file
		format:

		tree <oid>
		parent <oid parent>
		author: Author timestamp

		commit message
	*/
	treeOid, err := object.WriteTree(".")
	if err != nil {
		return "", err
	}

	fileContent := fmt.Appendf([]byte{}, "%v %v\n", object.TREE_OBJECT_TYPE, treeOid)

	// parent commit
	// @todo: commit the same working directory will cause parent to point to the same commit
	// How to check no changes since last commit?
	headRef, err := refs.GetCurrentHead()
	if err != nil {
		return "", fmt.Errorf("failed to read HEAD file, %v", err)
	}

	prevCommitOid, err := refs.GetRefContent(headRef)
	if err != nil {
		return "", fmt.Errorf("failed to read ref file, %v", err)
	}
	fileContent = fmt.Appendf(fileContent, "parent %v\n", prevCommitOid)

	// author
	usr, err := user.Current()
	var username string
	if err == nil {
		username = usr.Username
	} else {
		username = "<anonymous>"
	}
	fileContent = fmt.Appendf(fileContent, "author %v %v\n\n", username, time.Now().Unix())

	// commit message
	fileContent = append(fileContent, []byte(msg)...)

	commitOid, err := object.Write(object.COMMIT_OBJECT_TYPE, fileContent)
	if err != nil {
		return "", err
	}

	// @todo: How to handle if error happened here?
	refs.SetRefContent(headRef, commitOid)

	return commitOid, nil
}
