package blob

import (
	"notgit/internal/object"
	"notgit/internal/utils"
)

func (blob *Blob) Write(store object.Store) error {
	content := blob.Content
	compressed := utils.Compress("blob", content)

	return store.Write(blob.Hash(), compressed)
}
