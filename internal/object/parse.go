package object

import (
	"notgit/internal/utils"
	"os"
	"path/filepath"
	"strings"
)

// Returns header, content, error
func Parse(hash string) ([]byte, []byte, error) {
	wd, err := os.Getwd()
	if err != nil {
		return nil, nil, err
	}
	objects := filepath.Join(wd, ".notgit", "objects")

	dir := hash[:2]
	file := hash[2:]
	path := filepath.Join(objects, dir, file)

	content, err := os.ReadFile(path)
	if err != nil {
		println("here")
		return nil, nil, err
	}
	content, err = utils.Decompress(content)
	header := []byte(strings.Split(string(content), "\n")[0])
	header = header[:len(header)-1]
	content = []byte(strings.Join(strings.Split(string(content), "\n")[1:], "\n"))

	return header, content, err
}
