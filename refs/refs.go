package refs

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"

	"micro-git/root"
)

const (
	HEAD_FILE  = "HEAD"
	REF_PREFIX = "ref"
)

func GetCurrentHead() (string, error) {
	/*
		returns the content of HEAD file.
		e.g: refs/heads/master
	*/
	content, err := os.ReadFile(filepath.Join(root.FOLDER_NAME, "HEAD"))
	if err != nil {
		return "", err
	}

	splitted := bytes.Fields(content)
	if len(splitted) != 2 {
		return "", fmt.Errorf("HEAD file content is corrupted")
	}

	return string(splitted[1]), nil
}

func GetRefContent(refPath string) (string, error) {
	/*
		reads the content of ref pointed by HEAD (an OID to tree object file)
		and returns it.
	*/
	commitOid, err := os.ReadFile(filepath.Join(root.FOLDER_NAME, refPath))
	if err != nil {
		return "", err
	}

	return string(commitOid), nil
}

func SetRefContent(refPath string, commitOid string) (string, error) {
	/*
		set the value of a ref to commitOid provided.
	*/
	err := os.WriteFile(filepath.Join(root.FOLDER_NAME, refPath), []byte(commitOid), 0o664)
	if err != nil {
		return "", err
	}

	return commitOid, nil
}
