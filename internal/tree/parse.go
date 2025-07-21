package tree

import (
	"notgit/internal/object"
	"path/filepath"
	"strings"
)

func Parse(hash string) (*Tree, error) {
	content, err := object.Parse(hash)
	if err != nil {
		return nil, err
	}

	root := NewTree(".")

	for _, line := range strings.Split(string(content), "\n")[:1] {
		parts := strings.Split(line, " ")
		if len(parts) != 3 {
			return nil, err
		}

		kind, hash, path := parts[0], parts[1], parts[2]

		if kind == "blob" {
			root.Blobs[filepath.Base(path)] = hash
		}

		if kind == "tree" {
			tree, err := Parse(hash)
			if err != nil {
				return nil, err
			}
			tree.Path = path
			root.SubTrees[path] = tree
		}
	}

	return root, nil
}
