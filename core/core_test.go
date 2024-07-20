package core

import (
	"crypto/sha1"
	"encoding/hex"
	"os"
	"path/filepath"
	"testing"
)

func TestInit(t *testing.T) {
	currWd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current working directory, %v", err)
	}

	tmpDir := createTestDir(t)
	defer os.RemoveAll(tmpDir)
	defer os.Chdir(currWd)

	err = Init()
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

	tmpDir := createTestDir(t)
	defer os.RemoveAll(tmpDir)
	defer os.Chdir(currWd)

	os.WriteFile("test.txt", []byte("Hello"), 0o664)

	hexSum, err := HashObject("test.txt", "blob", false)
	if err != nil {
		t.Fatalf("HashObject return error: %v", err)
	}

	combined := append([]byte("blob"), '\x00')
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

	tmpDir := createTestDir(t)
	defer os.RemoveAll(tmpDir)
	defer os.Chdir(currWd)

	err = Init()
	if err != nil {
		t.Fatalf("Failed to execute Init command, error: %v", err)
	}

	os.WriteFile("test.txt", []byte("Hello"), 0o664)

	hexSum, err := HashObject("test.txt", "blob", true)
	if err != nil {
		t.Fatalf("HashObject return error: %v", err)
	}

	combined := append([]byte("blob"), '\x00')
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

	tmpDir := createTestDir(t)
	defer os.RemoveAll(tmpDir)
	defer os.Chdir(currWd)

	os.WriteFile("test.txt", []byte("Hello"), 0o664)

	_, err = HashObject("test.txt", "invalid", false)
	if err == nil {
		t.Fatalf("HashObject should return error, but it doesn't")
	}
}

func createTestDir(t *testing.T) string {
	tmpDir, err := os.MkdirTemp("", "microgit-test")
	if err != nil {
		os.RemoveAll(tmpDir)
		t.Fatalf("Failed to create temp dir for testing, %v", err)
	}

	err = os.Chdir(tmpDir)
	if err != nil {
		os.RemoveAll(tmpDir)
		t.Fatalf("Failed to switch working directory to temp dir, %v", err)
	}

	return tmpDir
}
