package tree

import "notgit/internal/blob"

type Tree struct {
	Hash     string
	Path     string
	SubTrees []*Tree
	Blobs    []blob.Blob
}
