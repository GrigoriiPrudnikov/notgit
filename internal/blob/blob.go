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
		Path:    path,
		Content: b,
	}

	return blob, err
}

func (b Blob) BasePath() string {
	return filepath.Base(b.Path)
}

func (b *Blob) Hash() string {
	return utils.Hash("blob", b.Content)
}

