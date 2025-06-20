package blob

import (
	"notgit/internal/utils"
	"os"
	"path/filepath"
)

type Blob struct {
	Path    string
	Content []byte
}

func NewBlob(path string) (Blob, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return Blob{}, err
	}

	blob := Blob{
		Path:    filepath.Base(path),
		Content: b,
	}

	return blob, err
}

func (b *Blob) Hash() string {
	return utils.Hash("blob", b.Content)
}

func (b *Blob) exists() bool {
	wd, err := os.Getwd()
	if err != nil {
		return false
	}

	objects := filepath.Join(wd, ".notgit", "objects")
	dir := filepath.Join(objects, b.Hash()[:2])
	file := filepath.Join(dir, b.Hash()[2:])

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
