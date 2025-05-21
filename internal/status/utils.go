package status

import (
	"notgit/internal/blob"
	"notgit/internal/tree"
)

func compare(a, b *tree.Tree) (difference []blob.Blob) {
	if b == nil {
		for _, blob := range a.Blobs {
			difference = append(difference, blob)
		}
		return
	}
	for _, blob := range a.Blobs {
		found := findBlob(b.Blobs, blob.Path)
		if found == nil || blob.Hash != found.Hash {
			difference = append(difference, blob)
		}
	}

	return
}

func findBlob(b []blob.Blob, name string) *blob.Blob {
	for _, blob := range b {
		if blob.Path == name {
			return &blob
		}
	}

	return nil
}

func findTree(t []*tree.Tree, name string) *tree.Tree {
	for _, tree := range t {
		if tree.Path == name {
			return tree
		}
	}
	return nil
}
