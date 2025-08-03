package object

import (
	"bytes"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"micro-git/root"
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

func GenInfo(objectType string, fileContent []byte) (*ObjectInfo, error) {
	if objectType != BLOB_OBJECT_TYPE &&
		objectType != TAG_OBJECT_TYPE &&
		objectType != COMMIT_OBJECT_TYPE &&
		objectType != TREE_OBJECT_TYPE {
		err := fmt.Errorf("invalid objectType supplied: %v", objectType)
		return nil, err
	}

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
	}, nil
}

func Write(objectType string, fileContent []byte) (string, error) {
	if objectType != BLOB_OBJECT_TYPE &&
		objectType != TAG_OBJECT_TYPE &&
		objectType != COMMIT_OBJECT_TYPE &&
		objectType != TREE_OBJECT_TYPE {
		err := fmt.Errorf("invalid objectType supplied: %v", objectType)
		return "", err
	}

	// error is not possible because the only error that can happened inside GenInfo
	// already handled in this function as well
	objectInfo, _ := GenInfo(objectType, fileContent)

	initial, fileId := objectInfo.Oid[:2], objectInfo.Oid[2:]
	folderName := filepath.Join(root.FOLDER_NAME, "objects", initial)
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

	path := filepath.Join(root.FOLDER_NAME, "objects", folderPrefix, fileName)
	fileContent, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	splitted := bytes.Split(fileContent, []byte("\x00"))

	// commit object file
	if len(splitted) == 1 {
		commitContent := splitted[0]
		return &ObjectInfo{
			Type:       COMMIT_OBJECT_TYPE,
			Size:       len(commitContent),
			Content:    commitContent,
			RawContent: commitContent,
			Oid:        oid,
		}, nil
	}

	// blob or tree object file
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

func ReadTree(oid string) error {
	treeEntries := make([]treeEntry, 0, 10)

	err := recursivelyReadTree(oid, ".", &treeEntries)
	if err != nil {
		return err
	}

	// @todo: How to make sure operation is atomic?
	for _, treeEntry := range treeEntries {
		filenamePath := treeEntry.filename
		isDir := treeEntry.objectType == TREE_OBJECT_TYPE
		entryOid := treeEntry.oid

		if isDir {
			err := os.MkdirAll(filenamePath, 0o777)
			if err != nil {
				fmt.Printf("dir %v failed to read, skipping...\n", filenamePath)
				continue
			}
		} else {
			objectInfo, err := Read(entryOid)
			if err != nil {
				fmt.Printf("file %v failed to read, skipping...\n", filenamePath)
				continue
			}

			path := filepath.Dir(filenamePath)

			err = os.MkdirAll(path, 0o777)
			if err != nil {
				fmt.Printf("file %v failed to write to directory, %v, skipping...\n", err, filenamePath)
				continue
			}

			err = os.WriteFile(filenamePath, objectInfo.Content, 0o664)
			if err != nil {
				fmt.Printf("file %v failed to write to directory, %v, skipping...\n", err, filenamePath)
				continue
			}
		}
	}

	return nil
}

func recursivelyReadTree(oid, prefix string, treeEntries *[]treeEntry) error {
	objectInfo, err := Read(oid)
	if err != nil {
		return err
	}

	fileContent := string(objectInfo.Content)
	lines := strings.Split(fileContent, "\n")

	for _, line := range lines {
		// handle empty line due to write-tree join everything with \n
		if line == "" {
			continue
		}

		segment := strings.Split(line, "\t")
		objectInfos := segment[0]
		filename := segment[1]

		subSegment := strings.Split(objectInfos, " ")
		objectType := subSegment[0]
		entryOid := subSegment[1]

		filenameJoinedByPath := filepath.Join(prefix, filename)

		*treeEntries = append(*treeEntries, treeEntry{
			filename:   filenameJoinedByPath,
			oid:        entryOid,
			objectType: objectType,
		})

		if objectType == TREE_OBJECT_TYPE {
			recursivelyReadTree(entryOid, filenameJoinedByPath, treeEntries)
		} else if objectType != BLOB_OBJECT_TYPE {
			return fmt.Errorf("read-tree found unidentifiable object type %v", objectType)
		}
	}

	return nil
}
