package tree

import (
	"errors"
	"notgit/internal/object"
	"path/filepath"
	"strings"
)

func Parse(hash, path string) (*Tree, error) {
	_, content, err := object.Parse(hash)
	if err != nil {
		return nil, err
	}

	root := NewTree(path)

	for _, line := range strings.Split(string(content), "\n") {
		if line == "" {
			continue
		}

		parts := strings.Split(line, " ")
		if len(parts) != 3 {
			return nil, errors.New("invalid tree line")
		}

		kind, path, hash := parts[0], parts[1], parts[2]

		if kind == "blob" {
			root.Blobs[filepath.Base(path)] = hash
		}

		if kind == "tree" {
			tree, err := Parse(hash, filepath.Clean(filepath.Join(root.Path, path)))
			if err != nil {
				return nil, err
			}
			root.SubTrees[path] = tree
		}
	}

	return root, nil
}
