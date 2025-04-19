package blob

import (
	"fmt"
	"notgit/utils"
	"os"
	"path/filepath"
)

func NewBlob(path string) (Blob, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return Blob{}, err
	}

	info, err := os.Stat(path)
	if err != nil {
		return Blob{}, err
	}

	permission := fmt.Sprintf("%o", info.Mode().Perm())

	blob := Blob{
		Permission: permission,
		Path:       filepath.Base(path),
		Content:    b,
	}

	Hash(&blob)

	return blob, err
}

func (blob *Blob) Write() error {
	if blob.exists() {
		return nil
	}

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

	if _, err := os.Stat(file); err == nil {
		return nil
	}

	content := blob.Content
	header := fmt.Sprintf("blob %d\x00\n", len(content))
	compressed := utils.Compress(content, header)

	err = os.WriteFile(file, compressed, 0644)

	return err
}
