package status

import (
	"notgit/internal/blob"
	"notgit/internal/tree"
	"path/filepath"
)

func GetModifiedAndUntrackedFiles(all, staged *tree.Tree) (modified, untracked []string) {
	difference := compare(all, staged)

	for _, a := range difference {
		if staged == nil {
			untracked = append(untracked, a.Path)
			continue
		}

		b := findBlob(staged.Blobs, a.Path)
		if b == nil {
			untracked = append(untracked, a.Path)
		} else {
			modified = append(modified, a.Path)
		}
	}

	for _, t := range all.SubTrees {
		var found *tree.Tree
		if staged != nil {
			found = findTree(staged.SubTrees, t.Path)
		}
		modifiedSub, untrackedSub := GetModifiedAndUntrackedFiles(t, found)

		for _, path := range modifiedSub {
			modified = append(modified, filepath.Join(t.Path, path))
		}
		for _, path := range untrackedSub {
			untracked = append(untracked, filepath.Join(t.Path, path))
		}
	}

	return
}

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
