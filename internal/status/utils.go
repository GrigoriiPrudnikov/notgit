package status

import (
	"notgit/internal/blob"
	"notgit/internal/tree"
)

func compare(a, b *tree.Tree) (difference []blob.Blob) {
	if a == nil && b == nil {
		return
	}
	if b == nil {
		for _, blob := range a.Blobs {
			difference = append(difference, blob)
		}
		return
	}
	if a == nil {
		for _, blob := range b.Blobs {
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
	for _, blob := range b.Blobs {
		found := findBlob(a.Blobs, blob.Path)
		if found == nil {
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

func extractPaths(blobs []blob.Blob) (paths []string) {
	for _, blob := range blobs {
		paths = append(paths, blob.Path)
	}
	return
}
