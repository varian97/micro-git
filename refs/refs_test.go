package refs

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/suite"

	"micro-git/root"
	"micro-git/testutil"
)

type RefsTestSuite struct {
	suite.Suite
	currWd string
	tmpDir string
}

func (suite *RefsTestSuite) SetupSuite() {
	currWd, err := os.Getwd()
	if err != nil {
		suite.FailNow("Failed to get current working directory", "Error: ", err)
	}

	suite.currWd = currWd

	tmpDir := testutil.CreateTestDir(suite.T())
	suite.tmpDir = tmpDir
}

func (suite *RefsTestSuite) TearDownSuite() {
	os.RemoveAll(suite.tmpDir)
	os.Chdir(suite.currWd)
}

func (suite *RefsTestSuite) SetupTest() {
	err := root.InitDB()
	if err != nil {
		suite.FailNow("Failed to execute Init command", "Error: %v", err)
	}
}

func (suite *RefsTestSuite) TearDownTest() {
	os.RemoveAll(".microgit")
}

func (suite *RefsTestSuite) TestGetCurrentHead() {
	refsPointedByHead, err := GetCurrentHead()
	if err != nil {
		suite.FailNow("Failed to read the HEAD file", "Error: %v", err)
	}

	if refsPointedByHead != "refs/heads/master" {
		suite.FailNow(
			"The content of HEAD file is not match",
			"Expected: refs/heads/master, got: %v", refsPointedByHead,
		)
	}
}

func (suite *RefsTestSuite) TestGetRefContentEmptyCommitOid() {
	commitOid, err := GetRefContent("refs/heads/master")
	if err != nil {
		suite.FailNow("Failed to read the ref file", "Error: %v", err)
	}

	if commitOid != "" {
		suite.FailNow(
			"The content of commit oid is not match",
			"Expected: '', got: %v", commitOid,
		)
	}
}

func (suite *RefsTestSuite) TestGetRefContentNonEmptyCommitOid() {
	expectedCommitOid := "asfdasldkfs1212321"
	err := os.WriteFile(filepath.Join(".microgit", "refs/heads/master"), []byte(expectedCommitOid), 0o664)
	if err != nil {
		suite.FailNow("Failed initialize the ref file", "Error: %v", err)
	}

	commitOid, err := GetRefContent("refs/heads/master")
	if err != nil {
		suite.FailNow("Failed to read the ref file", "Error: %v", err)
	}

	if commitOid != expectedCommitOid {
		suite.FailNow(
			"The content of commit oid is not match",
			"Expected: %v, got: %v", expectedCommitOid, commitOid,
		)
	}
}

func (suite *RefsTestSuite) TestSetRefContent() {
	expectedCommitOid := "asfdasldkfs1212321"

	_, err := SetRefContent("refs/heads/master", expectedCommitOid)
	if err != nil {
		suite.FailNow("Failed set the ref file content", "Error: %v", err)
	}

	commitOid, err := GetRefContent("refs/heads/master")
	if err != nil {
		suite.FailNow("Failed to read the ref file", "Error: %v", err)
	}

	if commitOid != expectedCommitOid {
		suite.FailNow(
			"The content of commit oid is not match",
			"Expected: %v, got: %v", expectedCommitOid, commitOid,
		)
	}
}

func TestRefsTestSuite(t *testing.T) {
	suite.Run(t, new(RefsTestSuite))
}
