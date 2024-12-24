package object

import (
	"crypto/sha1"
	"encoding/hex"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"micro-git/testutil"
)

func TestGenInfo(t *testing.T) {
	objectType := "blob"
	fileContent := []byte("Hello World")
	rawContent := append([]byte("blob "), []byte("11")...)
	rawContent = append(rawContent, '\x00')
	rawContent = append(rawContent, fileContent...)

	sha := sha1.Sum(rawContent)
	oid := hex.EncodeToString(sha[:])

	objectInfo := GenInfo(objectType, fileContent)

	if objectInfo.Type != objectType {
		t.Fatalf("objectInfo.Type is incorrect. Expected: %v, got: %v", objectType, objectInfo.Type)
	}
	if objectInfo.Size != len(fileContent) {
		t.Fatalf("objectInfo.Size is incorrect. Expected: %v, got: %v", len(fileContent), objectInfo.Size)
	}
	if string(objectInfo.Content) != string(fileContent) {
		t.Fatalf("objectInfo.Content is incorrect. Expected: %v, got: %v", fileContent, objectInfo.Content)
	}
	if string(objectInfo.RawContent) != string(rawContent) {
		t.Fatalf("objectInfo.RawContent is incorrect. Expected: %v, got: %v", rawContent, objectInfo.RawContent)
	}
	if objectInfo.Oid != oid {
		t.Fatalf("objectInfo.Oid is incorrect. Expected: %v, got: %v", oid, objectInfo.Oid)
	}
}

func TestWriteTree(t *testing.T) {
	currWd, _ := os.Getwd()
	defer os.Chdir(currWd)

	tmpDir := testutil.CreateTestDir(t)
	defer os.RemoveAll(tmpDir)

	err := InitDB()
	if err != nil {
		t.Fatalf("Failed to execute Init command, error: %v", err)
	}

	/*
		folder structure:
		.microgit
		test.txt
		src
		  test2.txt
	*/
	os.WriteFile("test.txt", []byte("Hello"), 0o664)

	srcDir := filepath.Join(tmpDir, "src")
	err = os.Mkdir(srcDir, 0o774)
	if err != nil {
		t.Fatalf("Failed to create src directory, error: %v", err)
	}
	defer os.RemoveAll(srcDir)

	os.WriteFile(filepath.Join(srcDir, "test2.txt"), []byte("Hello World"), 0o664)

	oid, err := WriteTree(".")
	if err != nil {
		t.Fatalf("Failed to execute WriteTree, error: %v", err)
	}

	objectInfo, err := Read(oid)
	if err != nil {
		t.Fatalf("The created tree cannot read, error: %v", err)
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
		t.Fatalf("test.txt not found in the tree")
	}

	if info.objectType != "blob" {
		t.Fatalf("test.txt object type is wrong, expected: blob, got: %v", info.objectType)
	}

	objInfo, err := Read(info.oid)
	if err != nil {
		t.Fatalf("test.txt cannot be read, error: %v", err)
	}
	if string(objInfo.Content) != "Hello" {
		t.Fatalf("test.txt content is wrong, expected: Hello, got: %v", string(objInfo.Content))
	}

	// assert the folder src
	info, ok = infoByFilename["src"]
	if !ok {
		t.Fatalf("src folder not found in the tree")
	}

	if info.objectType != "tree" {
		t.Fatalf("test.txt object type is wrong, expected: tree, got: %v", info.objectType)
	}

	objectInfo, err = Read(info.oid)
	if err != nil {
		t.Fatalf("src's tree object file cannot be read, error: %v", err)
	}

	// assert the src/test2.txt
	srcFolderEntry := strings.Trim(string(objectInfo.Content), "\n")
	parts := strings.Split(srcFolderEntry, "\t")
	headerParts := strings.Split(parts[0], " ")

	filename := parts[1]
	objectType := headerParts[0]
	entryOid := headerParts[1]

	if filename != "test2.txt" {
		t.Fatalf("The filename of file inside src folder is not correct, expected: test2.txt, got: %v", filename)
	}
	if objectType != "blob" {
		t.Fatalf("The object type of file inside src folder is not correct, expected: blob, got: %v", objectType)
	}

	objInfo, err = Read(entryOid)
	if err != nil {
		t.Fatalf("test2.txt cannot be read, error: %v", err)
	}
	if string(objInfo.Content) != "Hello World" {
		t.Fatalf("test2.txt content is wrong, expected: Hello World, got: %v", string(objInfo.Content))
	}
}
