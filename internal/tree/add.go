package tree

import (
	"notgit/internal/blob"
	"notgit/utils"
	"os"
	"path/filepath"
	"strings"
)

func (t *Tree) Add(path, fullPath string) error {
	if utils.Ignored(fullPath) {
		return nil
	}

	info, err := os.Stat(fullPath)
	if os.IsNotExist(err) {
		return err
	}

	parts := strings.Split(path, "/")

	if len(parts) == 1 {
		if info.IsDir() {
			var subtree *Tree

			// it checks if subtree already exists
			for path, sub := range t.SubTrees {
				if path == parts[0] {
					subtree = sub
					break
				}
			}

			// if subtree doesn't exist, creates it
			if subtree == nil {
				subtree = &Tree{}
				if t.SubTrees == nil {
					t.SubTrees = map[string]*Tree{}
				}
				t.SubTrees[parts[0]] = subtree
			}

			entries, err := os.ReadDir(fullPath)
			if err != nil {
				return err
			}

			for _, entry := range entries {
				err := subtree.Add(entry.Name(), filepath.Join(fullPath, entry.Name()))
				if err != nil {
					return err
				}
			}

			return nil
		}

		newBlob, err := blob.NewBlob(fullPath)
		if err != nil {
			return err
		}

		for i, blob := range t.Blobs {
			if blob.Path == newBlob.Path {
				t.Blobs[i].Content = newBlob.Content
				t.Blobs[i].Hash = newBlob.Hash
				return nil
			}
		}

		t.Blobs = append(t.Blobs, newBlob)
		return nil
	}

	for path, subtree := range t.SubTrees {
		if path == parts[0] {
			err := subtree.Add(filepath.Join(parts[1:]...), fullPath)
			if err != nil {
				return err
			}

			return nil
		}
	}

	subtree := Tree{}

	t.SubTrees[parts[0]] = &subtree

	return subtree.Add(filepath.Join(parts[1:]...), fullPath)
}
