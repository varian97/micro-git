package main

import (
	"crypto/sha1"
	"encoding/hex"
	"os"
	"path/filepath"
	"testing"

	"micro-git/object"
	"micro-git/testutil"
)

func TestInit(t *testing.T) {
	currWd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current working directory, %v", err)
	}

	tmpDir := testutil.CreateTestDir(t)
	defer os.RemoveAll(tmpDir)
	defer os.Chdir(currWd)

	err = object.InitDB()
	if err != nil {
		t.Fatalf("Failed to execute Init command, error: %v", err)
	}

	pathsToCheck := []string{
		".microgit",
		".microgit/refs",
		".microgit/refs/tags",
		".microgit/refs/heads",
		".microgit/objects",
		".microgit/HEAD",
	}

	for _, path := range pathsToCheck {
		if _, err := os.Stat(path); os.IsNotExist(err) {
			t.Fatalf("%v does not exists but expected to exists", path)
		}
	}

	headContent, err := os.ReadFile(filepath.Join(".microgit", "HEAD"))
	if err != nil {
		t.Fatalf("Failed to read the contents of HEAD file, %v", err)
	}

	expectedHeadContent := "ref: refs/heads/master"
	if string(headContent) != expectedHeadContent {
		t.Fatalf("The content of HEAD file is not match. Expected: %v, got: %v", expectedHeadContent, string(headContent))
	}
}

func TestHashBlobObjectNotWriteToDisk(t *testing.T) {
	currWd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current working directory, %v", err)
	}

	tmpDir := testutil.CreateTestDir(t)
	defer os.RemoveAll(tmpDir)
	defer os.Chdir(currWd)

	os.WriteFile("test.txt", []byte("Hello"), 0o664)

	hexSum := createFileAndHashIt(t, "Hello", false)

	combined := append([]byte("blob"), []byte(" 5")...)
	combined = append(combined, '\x00')
	combined = append(combined, []byte("Hello")...)
	shaSum := sha1.Sum(combined)
	expected := hex.EncodeToString(shaSum[:])

	if expected != hexSum {
		t.Fatalf("HashObject return incorrect hash. Expected: %v, got: %v", expected, hexSum)
	}
}

func TestHashBlobObjectWriteToDisk(t *testing.T) {
	currWd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current working directory, %v", err)
	}

	tmpDir := testutil.CreateTestDir(t)
	defer os.RemoveAll(tmpDir)
	defer os.Chdir(currWd)

	err = object.InitDB()
	if err != nil {
		t.Fatalf("Failed to execute Init command, error: %v", err)
	}

	hexSum := createFileAndHashIt(t, "Hello", true)

	combined := append([]byte("blob"), []byte(" 5")...)
	combined = append(combined, '\x00')
	combined = append(combined, []byte("Hello")...)
	shaSum := sha1.Sum(combined)
	expected := hex.EncodeToString(shaSum[:])

	if expected != hexSum {
		t.Fatalf("HashObject return incorrect hash. Expected: %v, got: %v", expected, hexSum)
	}

	objectPath := filepath.Join(".microgit", "objects", hexSum[:2], hexSum[2:])
	fileContent, err := os.ReadFile(objectPath)
	if err != nil {
		t.Fatalf("Error when opening the object file: %v", err)
	}

	if string(fileContent) != string(combined) {
		t.Fatalf("Object file contains wrong content. Expected: %v, got: %v", fileContent, combined)
	}
}

func TestHashObjectInvalidObjectType(t *testing.T) {
	currWd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current working directory, %v", err)
	}

	tmpDir := testutil.CreateTestDir(t)
	defer os.RemoveAll(tmpDir)
	defer os.Chdir(currWd)

	os.WriteFile("test.txt", []byte("Hello"), 0o664)

	createFileAndHashIt(t, "Hello", false)
}

func TestCatFileReturnCorrectResult(t *testing.T) {
	currWd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current working directory, %v", err)
	}

	tmpDir := testutil.CreateTestDir(t)
	defer os.RemoveAll(tmpDir)
	defer os.Chdir(currWd)

	err = object.InitDB()
	if err != nil {
		t.Fatalf("Failed to execute Init command, error: %v", err)
	}

	hexSum := createFileAndHashIt(t, "Hello", true)

	objInfo, err := CatFile(hexSum)
	if err != nil {
		t.Fatalf("CatFile return error: %v", err)
	}

	if objInfo.Type != "blob" {
		t.Fatalf("CatFile return wrong content type. Expected: blob, got: %v", objInfo.Type)
	}
	if objInfo.Size != 5 {
		t.Fatalf("CatFile return wrong content size. Expected: 5, got: %v", objInfo.Size)
	}
	if string(objInfo.Content) != "Hello" {
		t.Fatalf("CatFile return wrong content. Expected: %v, got: %v", []byte("Hello"), objInfo.Content)
	}
}

func createFileAndHashIt(t *testing.T, content string, shouldWrite bool) string {
	os.WriteFile("test.txt", []byte(content), 0o664)

	hexSum, err := HashObject("test.txt", "blob", shouldWrite)
	if err != nil {
		t.Fatalf("HashObject return error: %v", err)
	}

	return hexSum
}
