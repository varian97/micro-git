package object

import (
	"bytes"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"micro-git/db"
	"os"
	"path/filepath"
	"strconv"
)

const (
	BLOB_OBJECT_TYPE   = "blob"
	TAG_OBJECT_TYPE    = "tag"
	COMMIT_OBJECT_TYPE = "commit"
	TREE_OBJECT_TYPE   = "tree"
)

type ObjectInfo struct {
	Type       string
	Size       int
	RawContent []byte
	Content    []byte
	Oid        string
}

func GenInfo(objectType string, fileContent []byte) *ObjectInfo {
	combinedContent := append([]byte(objectType), []byte(fmt.Sprintf(" %v", len(fileContent)))...)
	combinedContent = append(combinedContent, '\x00')
	combinedContent = append(combinedContent, fileContent...)
	sha1Sum := sha1.Sum(combinedContent)
	hexSum := hex.EncodeToString(sha1Sum[:])

	return &ObjectInfo{
		Type:       objectType,
		Size:       len(fileContent),
		Content:    fileContent,
		RawContent: combinedContent,
		Oid:        hexSum,
	}
}

func Write(objectType string, fileContent []byte) (string, error) {
	if objectType != BLOB_OBJECT_TYPE &&
		objectType != TAG_OBJECT_TYPE &&
		objectType != COMMIT_OBJECT_TYPE &&
		objectType != TREE_OBJECT_TYPE {
		err := fmt.Errorf("invalid objectType supplied: %v", objectType)
		return "", err
	}

	objectInfo := GenInfo(objectType, fileContent)

	initial, fileId := objectInfo.Oid[:2], objectInfo.Oid[2:]
	folderName := filepath.Join(db.FOLDER_NAME, "objects", initial)
	fileName := filepath.Join(folderName, fileId)

	err := os.MkdirAll(folderName, 0o774)
	if err != nil {
		err := fmt.Errorf("failed to create object folder, %v", err)
		return "", err
	}

	err = os.WriteFile(fileName, objectInfo.RawContent, 0o664)
	if err != nil {
		err := fmt.Errorf("failed to create the object file, %v", err)
		return "", err
	}

	return objectInfo.Oid, nil
}

func Read(oid string) (*ObjectInfo, error) {
	folderPrefix, fileName := oid[:2], oid[2:]

	path := filepath.Join(db.FOLDER_NAME, "objects", folderPrefix, fileName)
	fileContent, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	splitted := bytes.Split(fileContent, []byte("\x00"))
	header, content := splitted[0], splitted[1]

	headerSplitted := bytes.Split(header, []byte(" "))
	objectType, sizeByte := headerSplitted[0], headerSplitted[1]

	size, err := strconv.Atoi(string(sizeByte))
	if err != nil {
		return nil, fmt.Errorf("failed to parse the object file: %v", err)
	}
	return &ObjectInfo{
		Type:       string(objectType),
		Size:       size,
		Content:    content,
		RawContent: fileContent,
		Oid:        oid,
	}, nil
}
