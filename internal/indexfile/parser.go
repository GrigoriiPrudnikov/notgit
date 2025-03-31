package indexfile

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
)

func Parse() ([]StagedFile, error) {
	var stagedFiles []StagedFile

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

		stagedFiles = append(stagedFiles, StagedFile{
			Permission: parts[0],
			Hash:       parts[1],
			Name:       parts[2],
		})
	}

	return stagedFiles, nil
}
