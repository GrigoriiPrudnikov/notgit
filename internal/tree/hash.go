package tree

import "notgit/utils"

func (t *Tree) Hash() string {
	content := []byte{}

	for _, blob := range t.Blobs {
		content = append(content, []byte("blob "+blob.Hash+" "+blob.Path+"\n")...)
	}
	for path, subtree := range t.SubTrees {
		content = append(content, []byte("blob "+subtree.Hash()+" "+path+"\n")...)
	}

	hex := utils.Hash("tree", content)
	return hex
}
