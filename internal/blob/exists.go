package blob

import (
	"os"
	"path/filepath"
)

func (b *Blob) exists() bool {
	wd, err := os.Getwd()
	if err != nil {
		return false
	}

	objects := filepath.Join(wd, ".notgit", "objects")
	dir := filepath.Join(objects, b.Hash[:2])
	file := filepath.Join(dir, b.Hash[2:])

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
