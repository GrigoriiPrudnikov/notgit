package blob

import (
	"os"
	"path/filepath"
)

// Add tree blob
func Create(path string) (Blob, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return Blob{}, err
	}

	hash := Hash(b)
	content := Compress(b, "blob")

	blob := Blob{
		Type:    "blob",
		Hash:    hash,
		Content: content,
		Size:    len(b),
	}

	return blob, err
}

func Write(blob Blob) error {
	wd, err := os.Getwd()
	if err != nil {
		return err
	}

	objects := filepath.Join(wd, ".notgit", "objects")

	if _, err := os.Stat(objects); os.IsNotExist(err) {
		err = os.MkdirAll(objects, 0755)
		if err != nil {
			return err
		}
	}

	hash := blob.Hash

	dir := filepath.Join(objects, hash[:2])
	file := filepath.Join(dir, hash[2:])

	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err = os.Mkdir(dir, 0755)
		if err != nil {
			return err
		}
	}

	if _, err := os.Stat(file); os.IsExist(err) {
		return nil
	}

	err = os.WriteFile(file, blob.Content, 0644)

	return err
}
