package testutil

import (
	"os"
	"testing"
)

func CreateTestDir(t *testing.T) string {
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
