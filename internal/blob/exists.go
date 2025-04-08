package blob

import (
	"os"
	"path/filepath"
)

func Exists(hash string) bool {
	wd, err := os.Getwd()
	if err != nil {
		return false
	}

	objects := filepath.Join(wd, ".notgit", "objects")
	dir := filepath.Join(objects, hash[:2])
	file := filepath.Join(dir, hash[2:])

	_, err = os.Stat(dir)
	if os.IsNotExist(err) {
		return false
	}

	_, err = os.Stat(file)
	if os.IsNotExist(err) {
		return false
	}

	return true
}
