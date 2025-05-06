package tree

import "notgit/internal/blob"

type Tree struct {
	Permission string
	Path       string
	SubTrees   []*Tree
	Blobs      []blob.Blob
}
