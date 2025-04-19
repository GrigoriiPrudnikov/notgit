package indexfile

import (
	"bytes"
	"notgit/internal/blob"
	"os"
	"path/filepath"
	"strings"
)

func Parse() ([]blob.Blob, error) {
	var stagedFiles []blob.Blob

	wd, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	indexPath := filepath.Join(wd, ".notgit", "index")
	if _, err := os.Stat(indexPath); os.IsNotExist(err) {
		os.WriteFile(indexPath, []byte(""), 0644)
	}

	b, err := os.ReadFile(indexPath)
	if err != nil {
		return nil, err
	}

	lines := bytes.Split(b, []byte("\n"))

	for _, line := range lines {
		parts := strings.Split(string(line), " ")

		if len(parts) != 3 {
			continue
		}

		b, err := blob.NewBlob(filepath.Join(parts[2]))
		if err != nil {
			return nil, err
		}

		b.Path = parts[2]
		stagedFiles = append(stagedFiles, b)
	}

	return stagedFiles, nil
}
