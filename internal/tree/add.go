package tree

import (
	"fmt"
	"notgit/internal/blob"
	"os"
	"path/filepath"
	"strings"
)

func (t *Tree) Add(path, fullPath string) error {
	info, err := os.Stat(fullPath)
	if os.IsNotExist(err) {
		return err
	}

	parts := strings.Split(path, "/")

	if len(parts) == 1 {
		if info.IsDir() {
			var subtree *Tree

			// it checks if subtree already exists
			for _, sub := range t.SubTrees {
				if sub.Path == parts[0] {
					subtree = sub
					break
				}
			}

			// if subtree doesn't exist, creates it
			if subtree == nil {
				subtree = &Tree{
					Path:       parts[0],
					Permission: fmt.Sprintf("%o", info.Mode().Perm()),
				}
				t.SubTrees = append(t.SubTrees, subtree)
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

			Hash(subtree)
			return nil
		}

		b, err := blob.NewBlob(fullPath)
		if err != nil {
			return err
		}

		for _, blob := range t.Blobs {
			if blob.Path == b.Path {
				return nil
			}
		}

		t.Blobs = append(t.Blobs, b)
		return nil
	}

	for _, subtree := range t.SubTrees {
		if subtree.Path == parts[0] {
			err := subtree.Add(filepath.Join(parts[1:]...), fullPath)
			if err != nil {
				return err
			}

			Hash(subtree)
			return nil
		}
	}

	subtree := Tree{
		Path:       parts[0],
		Permission: fmt.Sprintf("%o", info.Mode().Perm()),
	}

	err = subtree.Add(filepath.Join(parts[1:]...), fullPath)
	if err != nil {
		return err
	}

	Hash(&subtree)
	t.SubTrees = append(t.SubTrees, &subtree)

	return nil
}
