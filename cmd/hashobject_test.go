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

type HashObjectTestSuite struct {
	suite.Suite
	currWd string
	tmpDir string
}

func (suite *HashObjectTestSuite) SetupSuite() {
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

func (suite *HashObjectTestSuite) TearDownSuite() {
	os.RemoveAll(suite.tmpDir)
	os.Chdir(suite.currWd)
}

func (suite *HashObjectTestSuite) SetupTest() {
	err := root.InitDB()
	if err != nil {
		suite.FailNow("Failed to execute Init command", "Error: %v", err)
	}
}

func (suite *HashObjectTestSuite) TearDownTest() {
	os.RemoveAll(".microgit")
}

func (suite *HashObjectTestSuite) TestHashBlobObjectNotWriteToDisk() {
	hexSum, err := hashObject("test.txt", "blob", false)
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

func (suite *HashObjectTestSuite) TestHashBlobObjectWriteToDisk() {
	hexSum, err := hashObject("test.txt", "blob", true)
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

func (suite *HashObjectTestSuite) TestHashObjectInvalidObjectType() {
	hexSum, err := hashObject("test.txt", "invalid_object_type", false)

	suite.Empty(hexSum, "HashObject should return empty result because of invalid object type")
	suite.NotEmpty(err, "HashObject should return error because of invalid object type")
}

func TestHashObjectTestSuite(t *testing.T) {
	suite.Run(t, new(HashObjectTestSuite))
}
