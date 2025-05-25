package object

import (
	"crypto/sha1"
	"encoding/hex"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"micro-git/root"
	"micro-git/testutil"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

func TestGenInfo(t *testing.T) {
	objectType := "blob"
	fileContent := []byte("Hello World")
	rawContent := append([]byte("blob "), []byte("11")...)
	rawContent = append(rawContent, '\x00')
	rawContent = append(rawContent, fileContent...)

	sha := sha1.Sum(rawContent)
	oid := hex.EncodeToString(sha[:])

	objectInfo, err := GenInfo(objectType, fileContent)
	if err != nil {
		t.Fatalf("GenInfo returns error %v", err)
	}

	assert.Equal(t, objectType, objectInfo.Type)
	assert.Equal(t, len(fileContent), objectInfo.Size)
	assert.Equal(t, string(fileContent), string(objectInfo.Content))
	assert.Equal(t, string(rawContent), string(objectInfo.RawContent))
	assert.Equal(t, oid, objectInfo.Oid)
}

type TreeAndCommitTestSuite struct {
	suite.Suite
	currWd string
	tmpDir string
}

func (suite *TreeAndCommitTestSuite) SetupSuite() {
	currWd, err := os.Getwd()
	if err != nil {
		suite.FailNow("Failed to get current working directory", "Error: %v", err)
	}

	suite.currWd = currWd

	tmpDir := testutil.CreateTestDir(suite.T())
	suite.tmpDir = tmpDir
}

func (suite *TreeAndCommitTestSuite) TearDownSuite() {
	os.RemoveAll(suite.tmpDir)
	os.Chdir(suite.currWd)
}

func (suite *TreeAndCommitTestSuite) SetupTest() {
	err := root.InitDB()
	if err != nil {
		suite.FailNow("Failed to execute Init command", "Error: %v", err)
	}

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

func (suite *TreeAndCommitTestSuite) TearDownTest() {
	os.RemoveAll(".microgit")
	os.RemoveAll("src")
	os.RemoveAll("empty")
	os.Remove("test.txt")
}

func (suite *TreeAndCommitTestSuite) TestWriteTree() {
	oid, err := WriteTree(".")
	if err != nil {
		suite.FailNow("Failed to execute WriteTree", "Error: %v", err)
	}

	objectInfo, err := Read(oid)
	if err != nil {
		suite.FailNow("Cannot read the tree file", "Error: %v", err)
	}

	infoByFilename := make(map[string]struct {
		oid        string
		objectType string
	})
	entries := strings.Split(strings.Trim(string(objectInfo.Content), "\n"), "\n")
	for _, entry := range entries {
		parts := strings.Split(entry, "\t")
		headerParts := strings.Split(parts[0], " ")

		filename := parts[1]
		objectType := headerParts[0]
		entryOid := headerParts[1]

		infoByFilename[filename] = struct {
			oid        string
			objectType string
		}{entryOid, objectType}
	}

	// assert the single file test.txt
	info, ok := infoByFilename["test.txt"]
	if !ok {
		suite.FailNow("test.txt does not found inside the tree file")
	}
	suite.Equal("blob", info.objectType)

	objInfo, err := Read(info.oid)
	if err != nil {
		suite.FailNow("Failed to read the stored version of test.txt", "Error: %v", err)
	}
	suite.Equal("Hello", string(objInfo.Content))

	// assert the src folder
	info, ok = infoByFilename["src"]
	if !ok {
		suite.FailNow("src folder does not found inside the tree file")
	}
	suite.Equal("tree", info.objectType)

	// assert the src/test2.txt
	objectInfo, err = Read(info.oid)
	if err != nil {
		suite.FailNow("Failed to read src folder's tree file", "Error: %v", err)
	}

	srcFolderEntry := strings.Trim(string(objectInfo.Content), "\n")
	parts := strings.Split(srcFolderEntry, "\t")
	headerParts := strings.Split(parts[0], " ")

	filename := parts[1]
	objectType := headerParts[0]
	entryOid := headerParts[1]

	suite.Equal("test2.txt", filename)
	suite.Equal("blob", objectType)

	objInfo, err = Read(entryOid)
	if err != nil {
		suite.FailNow("Failed to read the stored version of test2.txt", "Error: %v", err)
	}
	suite.Equal("Hello World", string(objInfo.Content))

	// assert the empty folder
	info, ok = infoByFilename["empty"]
	if !ok {
		suite.FailNow("empty folder does not found inside the tree file")
	}
	suite.Equal("tree", info.objectType)
}

func (suite *TreeAndCommitTestSuite) TestReadTree() {
	oid, err := WriteTree(".")
	if err != nil {
		suite.FailNow("Failed to execute WriteTree", "Error: %v", err)
	}

	os.RemoveAll("src")
	os.RemoveAll("empty")
	os.Remove("test.txt")

	err = ReadTree(oid)
	if err != nil {
		suite.FailNow("ReadTree failed", "Error: %v", err)
	}

	_, err = os.Stat("test.txt")
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			suite.FailNow("test.txt not exists after read-tree")
		} else {
			suite.FailNow("test.txt not exists after read-tree", "Error: %v", err)
		}
	}

	_, err = os.Stat(filepath.Join("src", "test2.txt"))
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			suite.FailNow("test2.txt not exists after read-tree")
		} else {
			suite.FailNow("test2.txt not exists after read-tree", "Error: %v", err)
		}
	}

	_, err = os.Stat("src")
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			suite.FailNow("src folder not exists after read-tree")
		} else {
			suite.FailNow("src folder not exists after read-tree", "Error: %v", err)
		}
	}

	_, err = os.Stat("empty")
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			suite.FailNow("empty folder not exists after read-tree")
		} else {
			suite.FailNow("empty folder not exists after read-tree", "Error: %v", err)
		}
	}

	fileContent, err := os.ReadFile("test.txt")
	if err != nil {
		suite.FailNow("Failed to read test.txt file", "Error: %v", err)
	}
	suite.Equal("Hello", string(fileContent))

	fileContent, err = os.ReadFile(filepath.Join("src", "test2.txt"))
	if err != nil {
		suite.FailNow("Failed to read test2.txt file", "Error: %v", err)
	}
	suite.Equal("Hello World", string(fileContent))
}

func (suite *TreeAndCommitTestSuite) TestCommit() {
	oidFromWriteTree, err := WriteTree(".")
	if err != nil {
		suite.FailNow("Failed to execute WriteTree", "Error: %v", err)
	}

	oid, err := Commit("commit message")
	if err != nil {
		suite.FailNow("Failed to execute WriteTree", "Error: %v", err)
	}

	objectInfo, err := Read(oid)
	if err != nil {
		suite.FailNow("Cannot read the commit file", "Error: %v", err)
	}

	commitFileContent := string(objectInfo.Content)
	commitFileLines := strings.Split(commitFileContent, "\n")

	suite.Equal("commit message", commitFileLines[4])
	suite.Equal(oidFromWriteTree, strings.Split(commitFileLines[0], " ")[1])
}

func TestTreeAndCommitTestSuite(t *testing.T) {
	suite.Run(t, new(TreeAndCommitTestSuite))
}
