package filesystem

import (
	"os"
	"path/filepath"
	"testing"
)

func TestInit(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "microgit-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir for testing, %v", err)
	}
	defer os.RemoveAll(tmpDir)

	currWd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current working directory, %v", err)
	}

	err = os.Chdir(tmpDir)
	if err != nil {
		t.Fatalf("Failed to switch working directory to temp dir, %v", err)
	}
	defer os.Chdir(currWd)

	Init()

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
