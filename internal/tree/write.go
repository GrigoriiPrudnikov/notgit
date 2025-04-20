package tree

import (
	"notgit/utils"
	"os"
	"path/filepath"
)

func (t *Tree) Write() error {
	content := []byte{}

	for _, blob := range t.Blobs {
		err := blob.Write()
		if err != nil {
			return err
		}

		content = append(content, []byte(blob.Permission+" blob "+blob.Hash+" "+blob.Path+"\n")...)
	}

	for _, subtree := range t.SubTrees {
		err := subtree.Write()
		if err != nil {
			return err
		}

		content = append(content, []byte(subtree.Permission+" tree "+subtree.Hash+" "+subtree.Path+"\n")...)
	}

	compressed := utils.Compress(content, "tree")

	wd, err := os.Getwd()
	if err != nil {
		return err
	}
	objects := filepath.Join(wd, ".notgit", "objects")

	hash := t.Hash
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

	err = os.WriteFile(file, compressed, 0644)

	return err
}
