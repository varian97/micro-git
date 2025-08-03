package cmd

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"micro-git/object"
	"micro-git/root"
	"micro-git/testutil"

	"github.com/stretchr/testify/suite"
)

type CommitTestSuite struct {
	suite.Suite
	currWd string
	tmpDir string
}

func (suite *CommitTestSuite) SetupSuite() {
	currWd, err := os.Getwd()
	if err != nil {
		suite.FailNow("Failed to get current working directory", "Error: %v", err)
	}

	suite.currWd = currWd

	tmpDir := testutil.CreateTestDir(suite.T())
	suite.tmpDir = tmpDir

	// create .microgit root folder
	err = root.InitDB()
	if err != nil {
		suite.FailNow("Failed to execute Init command", "Error: %v", err)
	}

	// Create the directory contents
	/*
		folder structure:
		.microgit
		test.txt
		src
		  test2.txt
	*/
	err = os.WriteFile("test.txt", []byte("Hello"), 0o664)
	if err != nil {
		suite.FailNow("Failed to create test.txt file", "Error: %v", err)
	}

	srcDir := filepath.Join(suite.tmpDir, "src")
	err = os.Mkdir(srcDir, 0o774)
	if err != nil {
		suite.FailNow("Failed to create src directory", "Error: %v", err)
	}

	err = os.WriteFile(filepath.Join(srcDir, "test2.txt"), []byte("Hello World"), 0o664)
	if err != nil {
		suite.FailNow("Failed to create test2.txt file", "Error: %v", err)
	}

	emptyDir := filepath.Join(suite.tmpDir, "empty")
	err = os.Mkdir(emptyDir, 0o774)
	if err != nil {
		suite.FailNow("Failed to create empty directory", "Error: %v", err)
	}
}

func (suite *CommitTestSuite) TearDownSuite() {
	os.RemoveAll(suite.tmpDir)
	os.Chdir(suite.currWd)
}

func (suite *CommitTestSuite) TestCommit() {
	oidFromWriteTree, err := object.WriteTree(".")
	if err != nil {
		suite.FailNow("Failed to execute WriteTree", "Error: %v", err)
	}

	oid, err := commit("commit message")
	if err != nil {
		suite.FailNow("Failed to execute WriteTree", "Error: %v", err)
	}

	objectInfo, err := object.Read(oid)
	if err != nil {
		suite.FailNow("Cannot read the commit file", "Error: %v", err)
	}

	commitFileContent := string(objectInfo.Content)
	commitFileLines := strings.Split(commitFileContent, "\n")

	suite.Equal("commit message", commitFileLines[4])
	suite.Equal(oidFromWriteTree, strings.Split(commitFileLines[0], " ")[1])
}

func TestCommitTestSuite(t *testing.T) {
	suite.Run(t, new(CommitTestSuite))
}
