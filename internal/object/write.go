package object

import (
	"os"
	"path/filepath"
)

func Write(hash string, content []byte) error {
	wd, err := os.Getwd()
	if err != nil {
		return err
	}
	objects := filepath.Join(wd, ".notgit", "objects")

	dirPath := filepath.Join(objects, hash[:2])
	filePath := filepath.Join(dirPath, hash[2:])

	// create objects dir if not exists
	if _, err := os.Stat(dirPath); os.IsNotExist(err) {
		if err := os.MkdirAll(dirPath, 0755); err != nil {
			return err
		}
	}

	if _, err := os.Stat(filePath); err == nil {
		return nil
	} else if !os.IsNotExist(err) {
		return err
	}

	return os.WriteFile(filePath, content, 0644)
}
