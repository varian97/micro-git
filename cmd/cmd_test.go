package cmd

import (
	"crypto/sha1"
	"encoding/hex"
	"os"
	"path/filepath"
	"testing"

	"micro-git/object"
	"micro-git/root"
	"micro-git/testutil"

	"github.com/stretchr/testify/suite"
)

type CmdTestSuite struct {
	suite.Suite
	currWd string
	tmpDir string
}

func (suite *CmdTestSuite) SetupSuite() {
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
}

func (suite *CmdTestSuite) TearDownSuite() {
	os.RemoveAll(suite.tmpDir)
	os.Chdir(suite.currWd)
}

func (suite *CmdTestSuite) SetupTest() {
	err := root.InitDB()
	if err != nil {
		suite.FailNow("Failed to execute Init command", "Error: %v", err)
	}
}

func (suite *CmdTestSuite) TearDownTest() {
	os.RemoveAll(".microgit")
}

func (suite *CmdTestSuite) TestHashBlobObjectNotWriteToDisk() {
	hexSum, err := HashObject("test.txt", "blob", false)
	if err != nil {
		suite.FailNow("HashObject failed", "Error: %v", err)
	}

	combined := append([]byte("blob"), []byte(" 5")...)
	combined = append(combined, '\x00')
	combined = append(combined, []byte("Hello")...)
	shaSum := sha1.Sum(combined)
	expected := hex.EncodeToString(shaSum[:])

	suite.Equal(expected, hexSum)
}

func (suite *CmdTestSuite) TestHashBlobObjectWriteToDisk() {
	hexSum, err := HashObject("test.txt", "blob", true)
	if err != nil {
		suite.FailNow("HashObject failed", "Error: %v", err)
	}

	combined := append([]byte("blob"), []byte(" 5")...)
	combined = append(combined, '\x00')
	combined = append(combined, []byte("Hello")...)
	shaSum := sha1.Sum(combined)
	expected := hex.EncodeToString(shaSum[:])

	suite.Equal(expected, hexSum)

	objectPath := filepath.Join(".microgit", "objects", hexSum[:2], hexSum[2:])
	fileContent, err := os.ReadFile(objectPath)
	if err != nil {
		suite.FailNow("Error when opening the object file", "Error: %v", err)
	}

	suite.Equal(string(combined), string(fileContent))
}

func (suite *CmdTestSuite) TestHashObjectInvalidObjectType() {
	hexSum, err := HashObject("test.txt", "invalid_object_type", false)

	suite.Empty(hexSum, "HashObject should return empty result because of invalid object type")
	suite.NotEmpty(err, "HashObject should return error because of invalid object type")
}

func (suite *CmdTestSuite) TestCatFileReturnCorrectResult() {
	hexSum, err := HashObject("test.txt", "blob", true)
	if err != nil {
		suite.FailNow("HashObject failed", "Error: %v", err)
	}

	objInfo, err := object.Read(hexSum)
	if err != nil {
		suite.FailNow("CatFile failed: ", "Error: %v", err)
	}

	suite.Equal("blob", objInfo.Type)
	suite.Equal(5, objInfo.Size)
	suite.Equal("Hello", string(objInfo.Content))
}

func TestCmdTestSuite(t *testing.T) {
	suite.Run(t, new(CmdTestSuite))
}
