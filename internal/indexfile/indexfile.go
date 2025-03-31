package indexfile

import (
	"bytes"
	"os"
	"path/filepath"
)

// TODO: rewrite to function AddToIndex with garbage collection
func Set(stagedFiles []StagedFile) error {
	wd, err := os.Getwd()
	if err != nil {
		return err
	}

	indexPath := filepath.Join(wd, ".notgit", "index")
	if _, err := os.Stat(indexPath); os.IsNotExist(err) {
		os.WriteFile(indexPath, []byte(""), 0644)
	}

	var b bytes.Buffer
	for _, stagedFile := range stagedFiles {
		b.WriteString(stagedFile.Permission + " " + stagedFile.Hash + " " + stagedFile.Name + "\n")
	}

	err = os.WriteFile(indexPath, b.Bytes(), 0644)

	return err
}
