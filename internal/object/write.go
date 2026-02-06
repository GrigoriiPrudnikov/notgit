package object

import (
	"os"
	"path/filepath"
)

type Store interface {
	Write(hash string, content []byte) error
}

type FSStore struct {
	wd string
}

func NewFSStore(wd string) Store {
	return &FSStore{wd: wd}
}

func (s *FSStore) Write(hash string, content []byte) error {
	objects := filepath.Join(s.wd, ".notgit", "objects")

	dirPath := filepath.Join(objects, hash[:2])
	filePath := filepath.Join(dirPath, hash[2:])

	if err := os.MkdirAll(dirPath, 0755); err != nil {
		return err
	}

	if _, err := os.Stat(filePath); err == nil {
		// already exists
		return nil
	} else if !os.IsNotExist(err) {
		return err
	}

	return os.WriteFile(filePath, content, 0644)
}
