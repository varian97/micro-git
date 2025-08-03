package cmd

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"micro-git/object"
	"micro-git/refs"

	"github.com/akamensky/argparse"
)

type commitInfo struct {
	treeOid         string
	parentCommitOid string
	author          string
	time            int64
	message         string
}

type LogCommand struct {
	CommandMeta
}

func (c *LogCommand) Register(p *argparse.Parser) {
	c.Command = p.NewCommand("log", "Shows the commit logs.")
}

func (c *LogCommand) Handle() bool {
	if c.Command.Happened() {
		err := printCommitLogs()
		if err != nil {
			fmt.Println(err)
		}

		return true
	}
	return false
}

func printCommitLogs() error {
	headRef, err := refs.GetCurrentHead()
	if err != nil {
		return fmt.Errorf("failed to read HEAD file, %v", err)
	}

	commitOid, err := refs.GetRefContent(headRef)
	if err != nil {
		return fmt.Errorf("failed to read ref file, %v", err)
	}
	if commitOid == "" {
		return fmt.Errorf("fatal: Your branch does not have any commits yet")
	}

	commitInfo, err := parseCommitContent(commitOid)
	if err != nil {
		return err
	}

	/*
		commit: <oid>
		Author: author
		Date: <date>

		    commit message
	*/

	t := time.Unix(commitInfo.time, 0)
	fmt.Printf("commit: %v\nAuthor: %v\nDate: %v\n\n\t%v\n\n", commitOid, commitInfo.author, t.String(), commitInfo.message)

	for curr := commitInfo; curr.parentCommitOid != ""; {
		parentCommitInfo, err := parseCommitContent(curr.parentCommitOid)
		if err != nil {
			return err
		}
		t := time.Unix(parentCommitInfo.time, 0)
		fmt.Printf(
			"commit: %v\nAuthor: %v\nDate: %v\n\n\t%v\n\n",
			curr.parentCommitOid,
			parentCommitInfo.author,
			t.String(),
			parentCommitInfo.message,
		)
		curr = parentCommitInfo
	}

	return nil
}

func parseCommitContent(commitOid string) (*commitInfo, error) {
	objInfo, err := object.Read(commitOid)
	if err != nil {
		return nil, fmt.Errorf("failed to read commit file, %v", err)
	}
	commitContents := strings.Split(string(objInfo.Content), "\n\n")
	message := commitContents[1]

	infos := strings.Split(commitContents[0], "\n")
	treeOid := strings.Split(infos[0], " ")[1]
	parent := strings.Split(infos[1], " ")[1]

	authorAndTime := strings.Split(infos[2], " ")
	author := authorAndTime[1]
	timestampInt64, err := strconv.ParseInt(authorAndTime[2], 10, 64)
	if err != nil {
		timestampInt64 = 0
	}

	return &commitInfo{
		treeOid:         treeOid,
		parentCommitOid: parent,
		author:          author,
		time:            timestampInt64,
		message:         message,
	}, nil
}
