package tree

import (
	"notgit/internal/blob"
	"notgit/internal/object"
	"strings"
)

func Parse(hash string) (*Tree, error) {
	content, err := object.Parse(hash)
	if err != nil {
		return nil, err
	}

	root := &Tree{}

	for _, line := range strings.Split(string(content), "\n")[:1] {
		parts := strings.Split(line, " ")
		if len(parts) != 3 {
			return nil, err
		}

		kind, hash, path := parts[0], parts[1], parts[2]

		if kind == "blob" {
			content, err := object.Parse(hash)
			if err != nil {
				return nil, err
			}

			b := blob.Blob{
				Path:    path,
				Content: content,
			}

			root.Blobs = append(root.Blobs, b)
		}

		if kind == "tree" {
			tree, err := Parse(hash)
			if err != nil {
				return nil, err
			}
			root.SubTrees[path] = tree
		}
	}

	return root, nil
}
