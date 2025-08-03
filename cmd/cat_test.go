package cmd

import (
	"crypto/sha1"
	"encoding/hex"
	"os"
	"path/filepath"
	"testing"

	"micro-git/root"
	"micro-git/testutil"

	"github.com/stretchr/testify/suite"
)

type CatTestSuite struct {
	suite.Suite
	currWd    string
	tmpDir    string
	sampleOid string
}

func (suite *CatTestSuite) SetupSuite() {
	currWd, err := os.Getwd()
	if err != nil {
		suite.FailNow("Failed to get current working directory", "Error: %v", err)
	}

	suite.currWd = currWd

	tmpDir := testutil.CreateTestDir(suite.T())
	suite.tmpDir = tmpDir

	err = os.WriteFile("test.txt", []byte("Hello"), 0o664)
	if err != nil {
		suite.FailNow("Failed to create test.txt file for testing", "Error: %v", err)
	}

	// create the .microgit folder
	err = root.InitDB()
	if err != nil {
		suite.FailNow("Failed to execute Init command", "Error: %v", err)
	}

	// hash the test.txt file
	combined := append([]byte("blob"), []byte(" 5")...)
	combined = append(combined, '\x00')
	combined = append(combined, []byte("Hello")...)
	shaSum := sha1.Sum(combined)
	hexSum := hex.EncodeToString(shaSum[:])

	suite.sampleOid = hexSum

	folderPath := filepath.Join(".microgit", "objects", hexSum[:2])
	objectPath := filepath.Join(folderPath, hexSum[2:])

	err = os.MkdirAll(folderPath, 0o774)
	if err != nil {
		suite.FailNow("Failed to create folder for hashed file of test.txt", "Error: %v", err)
	}

	err = os.WriteFile(objectPath, combined, 0o664)
	if err != nil {
		suite.FailNow("Failed to create hashed file of test.txt", "Error: %v", err)
	}
}

func (suite *CatTestSuite) TearDownSuite() {
	os.RemoveAll(suite.tmpDir)
	os.Chdir(suite.currWd)
}

func (suite *CatTestSuite) TestCatFileContentReturnCorrectResult() {
	content, err := catFile(suite.sampleOid, false, false)
	if err != nil {
		suite.FailNow("catFile failed", "Error: %v", err)
	}
	suite.Equal("Hello", string(content))
}

func (suite *CatTestSuite) TestCatFileTypeReturnCorrectResult() {
	fileType, err := catFile(suite.sampleOid, true, false)
	if err != nil {
		suite.FailNow("catFile failed", "Error: %v", err)
	}
	suite.Equal("blob", fileType)
}

func (suite *CatTestSuite) TestCatFileSizeReturnCorrectResult() {
	fileSize, err := catFile(suite.sampleOid, false, true)
	if err != nil {
		suite.FailNow("catFile failed", "Error: %v", err)
	}
	suite.Equal("5", fileSize)
}

func TestCatTestSuite(t *testing.T) {
	suite.Run(t, new(CatTestSuite))
}
