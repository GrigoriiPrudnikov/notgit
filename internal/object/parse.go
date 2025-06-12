package object

import (
	"notgit/internal/utils"
	"os"
	"path/filepath"
	"strings"
)

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
	content, err = utils.Decompress(content)
	content = []byte(strings.Join(strings.Split(string(content), "\n")[1:], ""))

	return content, err
}
