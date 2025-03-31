package blob

import (
	"bytes"
	"compress/zlib"
	"os"
	"path/filepath"
)

// TODO: add blob type (type, size, name, content, etc)
// TODO: add compress util function
func Create(b []byte) error {
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

	hash := Hash(b)

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

	var compressed bytes.Buffer
	w := zlib.NewWriter(&compressed)
	w.Write(b)
	w.Close()

	err = os.WriteFile(file, compressed.Bytes(), 0644)

	return err
}
