package indexfile

import (
	"bytes"
	"notgit/internal/utils"
	"os"
	"path/filepath"
	"strings"
)

func Parse() (map[string]string, error) {
	stagedFiles := make(map[string]string)

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

		if len(parts) != 2 {
			continue
		}

		stagedFiles[parts[0]] = parts[1]
	}

	return stagedFiles, nil
}

func Write(stagedFiles map[string]string) error {
	content := []byte{}
	paths := utils.GetSortedKeys(stagedFiles)

	for _, path := range paths {
		content = append(content, []byte(path+" "+stagedFiles[path]+"\n")...)
	}

	wd, err := os.Getwd()
	if err != nil {
		return err
	}

	indexPath := filepath.Join(wd, ".notgit", "index")
	err = os.WriteFile(indexPath, content, 0644)

	return err
}
