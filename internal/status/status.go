package status

import (
	"fmt"
	"notgit/internal/blob"
	"notgit/internal/tree"
	"path/filepath"
)

func GetModifiedAndUntrackedFiles(all, staged *tree.Tree) (modified, untracked []string) {
	modified, untracked = compare(all, staged)

	for _, tree := range all.SubTrees {
		found := findTree(staged.SubTrees, tree.Path)
		modifiedSub, untrackedSub := GetModifiedAndUntrackedFiles(tree, found)

		for _, path := range modifiedSub {
			modified = append(modified, filepath.Join(tree.Path, path))
		}
		for _, path := range untrackedSub {
			untracked = append(untracked, filepath.Join(tree.Path, path))
		}
	}

	return
}

func compare(all, staged *tree.Tree) (modified, untracked []string) {
	if staged == nil {
		for _, blob := range all.Blobs {
			untracked = append(untracked, blob.Path)
		}
		return
	}
	for _, blob := range all.Blobs {
		found := findBlob(staged.Blobs, blob.Path)
		if found == nil {
			untracked = append(untracked, blob.Path)
			continue
		}

		fmt.Println("found:", blob.Path)
		fmt.Println(string(blob.Content))
		fmt.Println(string(found.Content))

		if !contentMatches(blob.Content, found.Content) {
			modified = append(modified, blob.Path)
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

func contentMatches(a, b []byte) bool {
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}

	return true
}
