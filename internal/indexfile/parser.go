package indexfile

import (
	"bytes"
	"notgit/internal/blob"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func Parse() ([]blob.File, error) {
	var stagedFiles []blob.File

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

		mode, err := strconv.ParseInt(parts[0], 10, 32)
		if err != nil {
			return nil, err
		}

		stagedFiles = append(stagedFiles, blob.File{
			Mode: os.FileMode(mode),
			Hash: parts[1],
			Name: parts[2],
		})
	}

	return stagedFiles, nil
}
