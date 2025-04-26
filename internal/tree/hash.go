package tree

import "notgit/utils"

func Hash(t *Tree) {
	content := []byte{}

	for _, blob := range t.Blobs {
		content = append(content, []byte(blob.Permission+" blob "+blob.Hash+" "+blob.Path+"\n")...)
	}
	for _, subtree := range t.SubTrees {
		content = append(content, []byte(subtree.Permission+" blob "+subtree.Hash+" "+subtree.Path+"\n")...)
	}

	hex := utils.Hash("tree", content)
	t.Hash = hex
}
