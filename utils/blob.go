package utils

import (
	"bytes"
	"compress/zlib"
	"crypto/sha256"
	"fmt"
	"os"
	"path/filepath"
)

func CreateBlob(b []byte) error {
	header := fmt.Sprintf("blob %d\x00\n", len(b))
	blob := append([]byte(header), b...)
	hash := sha256.Sum256(blob)
	hex := fmt.Sprintf("%x", hash)

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

	dir := filepath.Join(objects, hex[:2])
	file := filepath.Join(dir, hex[2:])

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
