package object

import (
	"bytes"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
)

const (
	FOLDER_NAME = ".microgit"
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

type treeEntry struct {
	objectType string
	oid        string
	filename   string
}

func InitDB() error {
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

func GenInfo(objectType string, fileContent []byte) *ObjectInfo {
	combinedContent := append([]byte(objectType), []byte(fmt.Sprintf(" %v\x00", len(fileContent)))...)
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
	folderName := filepath.Join(FOLDER_NAME, "objects", initial)
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

	path := filepath.Join(FOLDER_NAME, "objects", folderPrefix, fileName)
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

func WriteTree(prefix string) (string, error) {
	files, err := os.ReadDir(prefix)
	if err != nil {
		return "", err
	}

	entries := []treeEntry{}

	for _, file := range files {
		if file.Name() == ".microgit" {
			continue
		}

		if file.IsDir() {
			recursiveOid, err := WriteTree(filepath.Join(prefix, file.Name()))
			if err != nil {
				err := fmt.Errorf("failed to recursively create tree object %v: %v", file.Name(), err)
				return "", err
			}

			entries = append(entries, treeEntry{TREE_OBJECT_TYPE, recursiveOid, file.Name()})
		} else {
			fileContent, err := os.ReadFile(filepath.Join(prefix, file.Name()))
			if err != nil {
				err := fmt.Errorf("failed to read file content %v: %v", file.Name(), err)
				return "", err
			}

			oid, err := Write(BLOB_OBJECT_TYPE, fileContent)
			if err != nil {
				return "", err
			}

			entries = append(entries, treeEntry{BLOB_OBJECT_TYPE, oid, file.Name()})
		}
	}

	var treeFileContent []byte
	for _, entry := range entries {
		row := fmt.Sprintf("%v %v\t%v\n", entry.objectType, entry.oid, entry.filename)
		treeFileContent = append(treeFileContent, []byte(row)...)
	}

	return Write(TREE_OBJECT_TYPE, treeFileContent)
}
