package core

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
)

const (
	MICROGIT_FOLDER_NAME = ".microgit"

	BLOB_OBJECT_TYPE   = "blob"
	TAG_OBJECT_TYPE    = "tag"
	COMMIT_OBJECT_TYPE = "commit"
	TREE_OBJECT_TYPE   = "tree"
)

func Init() error {
	err := os.Mkdir(MICROGIT_FOLDER_NAME, 0o774)
	if err != nil {
		errorMessage := fmt.Errorf("failed to initialize .microgit folder: %v", err)
		return errorMessage
	}

	refsFolderPath := filepath.Join(MICROGIT_FOLDER_NAME, "refs")
	err = os.Mkdir(refsFolderPath, 0o774)
	if err != nil {
		errorMessage := fmt.Errorf("failed to initialize refs folder: %v", err)
		return errorMessage
	}

	refsHeadsFolderPath := filepath.Join(MICROGIT_FOLDER_NAME, "refs", "heads")
	err = os.Mkdir(refsHeadsFolderPath, 0o774)
	if err != nil {
		errorMessage := fmt.Errorf("failed to initialize refs folder: %v", err)
		return errorMessage
	}

	refsTagsFolderPath := filepath.Join(MICROGIT_FOLDER_NAME, "refs", "tags")
	err = os.Mkdir(refsTagsFolderPath, 0o774)
	if err != nil {
		errorMessage := fmt.Errorf("failed to initialize refs folder: %v", err)
		return errorMessage
	}

	objectsFolderPath := filepath.Join(MICROGIT_FOLDER_NAME, "objects")
	err = os.Mkdir(objectsFolderPath, 0o774)
	if err != nil {
		errorMessage := fmt.Errorf("failed to initialize objects folder: %v", err)
		return errorMessage
	}

	headsFilePath := filepath.Join(MICROGIT_FOLDER_NAME, "HEAD")
	err = os.WriteFile(headsFilePath, []byte("ref: refs/heads/master"), 0o664)
	if err != nil {
		errorMessage := fmt.Errorf("failed to initialize HEAD: %v", err)
		return errorMessage
	}

	return nil
}

func HashObject(path, objectType string, shouldWrite bool) (string, error) {
	if objectType != BLOB_OBJECT_TYPE && objectType != TAG_OBJECT_TYPE && objectType != COMMIT_OBJECT_TYPE && objectType != TREE_OBJECT_TYPE {
		err := fmt.Errorf("invalid objectType supplied: %v", objectType)
		return "", err
	}

	fileContent, err := os.ReadFile(path)
	if err != nil {
		err := fmt.Errorf("failed to read file content, %v", err)
		return "", err
	}

	combined := append([]byte(objectType), '\x00')
	combined = append(combined, fileContent...)
	sha1Sum := sha1.Sum(combined)
	hexSum := hex.EncodeToString(sha1Sum[:])

	if !shouldWrite {
		return hexSum, nil
	}

	initial, fileId := hexSum[:2], hexSum[2:]
	folderName := filepath.Join(MICROGIT_FOLDER_NAME, "objects", initial)
	fileName := filepath.Join(folderName, fileId)

	err = os.MkdirAll(folderName, 0o774)
	if err != nil {
		err := fmt.Errorf("failed to create object folder, %v", err)
		return "", err
	}

	err = os.WriteFile(fileName, combined, 0o664)
	if err != nil {
		err := fmt.Errorf("failed to create the object file, %v", err)
		return "", err
	}

	return hexSum, nil
}
