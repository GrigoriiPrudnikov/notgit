package indexfile

import (
	"bytes"
	"notgit/internal/blob"
	"os"
	"path/filepath"
)

// TODO: rewrite to function AddToIndex with garbage collection
func Write(files []blob.Blob) error {
	wd, err := os.Getwd()
	if err != nil {
		return err
	}

	indexPath := filepath.Join(wd, ".notgit", "index")
	if _, err := os.Stat(indexPath); os.IsNotExist(err) {
		os.WriteFile(indexPath, []byte(""), 0644)
	}

	var b bytes.Buffer
	for _, staged := range files {
		b.WriteString(staged.Permission + " " + staged.Hash + " " + staged.Path + "\n")
	}

	err = os.WriteFile(indexPath, b.Bytes(), 0644)

	return err
}
