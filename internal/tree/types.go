package tree

import "notgit/internal/blob"

type Tree struct {
	Path     string
	SubTrees []*Tree
	Blobs    []blob.Blob
}
