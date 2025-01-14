package root

import (
	"os"
	"path/filepath"
	"testing"

	"micro-git/testutil"
)

func TestInitDB(t *testing.T) {
	currWd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current working directory, %v", err)
	}

	tmpDir := testutil.CreateTestDir(t)
	defer os.RemoveAll(tmpDir)
	defer os.Chdir(currWd)

	err = InitDB()
	if err != nil {
		t.Fatalf("Failed to execute Init command, error: %v", err)
	}

	pathsToCheck := []string{
		".microgit",
		".microgit/refs",
		".microgit/refs/tags",
		".microgit/refs/heads",
		".microgit/refs/heads/master",
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
