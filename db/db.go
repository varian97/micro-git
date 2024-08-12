package db

import (
	"fmt"
	"os"
	"path/filepath"
)

const (
	FOLDER_NAME = ".microgit"
)

func Init() error {
	err := os.Mkdir(FOLDER_NAME, 0o774)
	if err != nil {
		errorMessage := fmt.Errorf("failed to initialize .microgit folder: %v", err)
		return errorMessage
	}

	refsFolderPath := filepath.Join(FOLDER_NAME, "refs")
	err = os.Mkdir(refsFolderPath, 0o774)
	if err != nil {
		errorMessage := fmt.Errorf("failed to initialize refs folder: %v", err)
		return errorMessage
	}

	refsHeadsFolderPath := filepath.Join(FOLDER_NAME, "refs", "heads")
	err = os.Mkdir(refsHeadsFolderPath, 0o774)
	if err != nil {
		errorMessage := fmt.Errorf("failed to initialize refs folder: %v", err)
		return errorMessage
	}

	refsTagsFolderPath := filepath.Join(FOLDER_NAME, "refs", "tags")
	err = os.Mkdir(refsTagsFolderPath, 0o774)
	if err != nil {
		errorMessage := fmt.Errorf("failed to initialize refs folder: %v", err)
		return errorMessage
	}

	objectsFolderPath := filepath.Join(FOLDER_NAME, "objects")
	err = os.Mkdir(objectsFolderPath, 0o774)
	if err != nil {
		errorMessage := fmt.Errorf("failed to initialize objects folder: %v", err)
		return errorMessage
	}

	headsFilePath := filepath.Join(FOLDER_NAME, "HEAD")
	err = os.WriteFile(headsFilePath, []byte("ref: refs/heads/master"), 0o664)
	if err != nil {
		errorMessage := fmt.Errorf("failed to initialize HEAD: %v", err)
		return errorMessage
	}

	return nil
}
