package object

import (
	"notgit/utils"
	"os"
	"path/filepath"
)

func Write(hash string, content []byte) error {
	wd, err := os.Getwd()
	if err != nil {
		return err
	}
	objects := filepath.Join(wd, ".notgit", "objects")

	dir := filepath.Join(objects, hash[:2])
	file := filepath.Join(dir, hash[2:])

	// create objects dir if not exists
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err = os.Mkdir(dir, 0755)
		if err != nil {
			return err
		}
	}

	if _, err := os.Stat(file); os.IsExist(err) {
		return nil
	}

	return os.WriteFile(file, content, 0644)
}

func Parse(hash string) ([]byte, error) {
	wd, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	objects := filepath.Join(wd, ".notgit", "objects")

	dir := filepath.Join(objects, hash[:2])
	file := filepath.Join(dir, hash[2:])

	content, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}

	return utils.Decompress(content)
}
