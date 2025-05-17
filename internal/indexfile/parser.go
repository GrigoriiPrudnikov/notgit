package indexfile

import (
	"bytes"
	"notgit/internal/blob"
	"notgit/internal/object"
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

		content, err := object.Parse(parts[1])
		if err != nil {
			return nil, err
		}

		b := blob.Blob{
			Permission: parts[0],
			Path:       parts[2],
			Hash:       parts[1],
			Content:    content,
		}

		stagedFiles = append(stagedFiles, b)
	}

	return stagedFiles, nil
}
