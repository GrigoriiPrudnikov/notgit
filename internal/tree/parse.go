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

	for _, line := range strings.Split(string(content), "\n") {
		parts := strings.Split(line, " ")
		if len(parts) != 4 {
			return nil, err
		}

		kind, hash, path := parts[1], parts[2], parts[3]

		if kind == "blob" {
			content, err := object.Parse(hash)
			if err != nil {
				return nil, err
			}

			b := blob.Blob{
				Hash:    hash,
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
			root.SubTrees = append(root.SubTrees, tree)
		}
	}

	return root, nil
}
