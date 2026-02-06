package commit

import (
	"notgit/internal/object"
	"notgit/internal/utils"
)

func (c *Commit) Write(store object.Store) (string, error) {
	if err := c.Tree.Write(store); err != nil {
		return "", err
	}

	content := c.GetContent()

	hash := c.Hash()
	compressed := utils.Compress("commit", content)

	if err := store.Write(hash, compressed); err != nil {
		return "", err
	}

	return hash, nil
}
