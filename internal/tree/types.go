package tree

import "notgit/internal/blob"

type Tree struct {
	SubTrees map[string]*Tree
	Blobs    []blob.Blob
}
