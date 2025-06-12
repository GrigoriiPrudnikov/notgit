package status

import (
	"notgit/internal/blob"
	"notgit/internal/commit"
	"notgit/internal/indexfile"
	"notgit/internal/tree"
	"path/filepath"
)

func GetStaged() (modified, added, deleted []string) {
	index, err := indexfile.Parse()
	if err != nil {
		return nil, nil, nil
	}
	stagedTree := tree.Staged(index)
	head := commit.ParseHead()

	var headTree *tree.Tree

	if head != nil {
		headTree, err = tree.Parse(head.Tree)
		if err != nil {
			return nil, nil, nil
		}
	}

	modifiedStagedBlobs, addedBlobs, deletedStagedBlobs := getModifiedUntrackedDeleted(stagedTree, headTree)

	modified = extractPaths(modifiedStagedBlobs)
	added = extractPaths(addedBlobs)
	deleted = extractPaths(deletedStagedBlobs)

	return
}

func GetUnstaged() (modified, untracked, deleted []string) {
	all := tree.Root()
	index, err := indexfile.Parse()
	if err != nil {
		return nil, nil, nil
	}
	stagedTree := tree.Staged(index)

	modifiedBlobs, untrackedBlobs, deletedBlobs := getModifiedUntrackedDeleted(all, stagedTree)

	modified = extractPaths(modifiedBlobs)
	untracked = extractPaths(untrackedBlobs)
	deleted = extractPaths(deletedBlobs)

	return
}

func getModifiedUntrackedDeleted(all, staged *tree.Tree) (modified, untracked, deleted []blob.Blob) {
	diff := compare(all, staged)

	for _, b := range diff {
		var foundA, foundB *blob.Blob
		if all != nil {
			foundA = findBlob(all.Blobs, b.Path)
		}
		if staged != nil {
			foundB = findBlob(staged.Blobs, b.Path)
		}

		if foundA == nil && foundB != nil {
			deleted = append(deleted, b)
			continue
		}
		if foundA != nil && foundB == nil {
			untracked = append(untracked, b)
			continue
		}

		if foundA != nil && foundB != nil && foundA.Hash != foundB.Hash {
			modified = append(modified, b)
		}
	}

	for path, t := range all.SubTrees {
		var found *tree.Tree
		if staged != nil {
			found = staged.SubTrees[path]
		}
		modifiedSub, untrackedSub, deletedSub := getModifiedUntrackedDeleted(t, found)

		for _, mod := range modifiedSub {
			mod.Path = filepath.Join(path, mod.Path)
			modified = append(modified, mod)
		}
		for _, untrack := range untrackedSub {
			untrack.Path = filepath.Join(path, untrack.Path)
			untracked = append(untracked, untrack)
		}
		for _, del := range deletedSub {
			del.Path = filepath.Join(path, del.Path)
			deleted = append(deleted, del)
		}
	}

	return
}

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

func extractPaths(blobs []blob.Blob) (paths []string) {
	for _, blob := range blobs {
		paths = append(paths, blob.Path)
	}
	return
}
