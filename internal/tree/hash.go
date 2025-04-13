package tree

import (
	"crypto/sha256"
	"fmt"
)

func Hash(t *Tree) string {
	content := []byte{}

	for _, blob := range t.Blobs {
		content = append(content, []byte(blob.Permission+" blob "+blob.Hash+" "+blob.Path+"\n")...)
	}

	for _, subtree := range t.SubTrees {
		content = append(content, []byte(subtree.Permission+" blob "+subtree.Hash+" "+subtree.Path+"\n")...)
	}

	header := fmt.Sprintf("tree %d\x00\n", len(content))
	blob := append([]byte(header), content...)
	hash := sha256.Sum256(blob)
	hex := fmt.Sprintf("%x", hash)

	return hex
}
