package filesystem

import (
	"fmt"
	"os"
	"path/filepath"
)

const (
	MICROGIT_FOLDER_NAME = ".microgit"
)

func Init() {
	err := os.Mkdir(MICROGIT_FOLDER_NAME, 0o774)
	if err != nil {
		errorMessage := fmt.Errorf("failed to initialize .microgit folder: %v", err)
		panic(errorMessage)
	}

	refsFolderPath := filepath.Join(MICROGIT_FOLDER_NAME, "refs")
	err = os.Mkdir(refsFolderPath, 0o774)
	if err != nil {
		errorMessage := fmt.Errorf("failed to initialize refs folder: %v", err)
		panic(errorMessage)
	}

	refsHeadsFolderPath := filepath.Join(MICROGIT_FOLDER_NAME, "refs", "heads")
	err = os.Mkdir(refsHeadsFolderPath, 0o774)
	if err != nil {
		errorMessage := fmt.Errorf("failed to initialize refs folder: %v", err)
		panic(errorMessage)
	}

	refsTagsFolderPath := filepath.Join(MICROGIT_FOLDER_NAME, "refs", "tags")
	err = os.Mkdir(refsTagsFolderPath, 0o774)
	if err != nil {
		errorMessage := fmt.Errorf("failed to initialize refs folder: %v", err)
		panic(errorMessage)
	}

	objectsFolderPath := filepath.Join(MICROGIT_FOLDER_NAME, "objects")
	err = os.Mkdir(objectsFolderPath, 0o774)
	if err != nil {
		errorMessage := fmt.Errorf("failed to initialize objects folder: %v", err)
		panic(errorMessage)
	}

	headsFilePath := filepath.Join(MICROGIT_FOLDER_NAME, "HEAD")
	err = os.WriteFile(headsFilePath, []byte("ref: refs/heads/master"), 0o664)
	if err != nil {
		errorMessage := fmt.Errorf("failed to initialize HEAD: %v", err)
		panic(errorMessage)
	}
}
