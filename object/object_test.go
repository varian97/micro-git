package object

import (
	"crypto/sha1"
	"encoding/hex"
	"testing"
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
