package status

import (
	"notgit/internal/blob"
	"notgit/internal/commit"
	"notgit/internal/tree"
	"path/filepath"
)

func GetStaged() (modified, added, deleted []string) {
	stagedTree := tree.Staged()
	head := commit.ParseHead()

	var headTree *tree.Tree
	var err error

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
	stagedTree := tree.Staged()

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

	for _, t := range all.SubTrees {
		var found *tree.Tree
		if staged != nil {
			found = findTree(staged.SubTrees, t.Path)
		}
		modifiedSub, untrackedSub, deletedSub := getModifiedUntrackedDeleted(t, found)

		for _, mod := range modifiedSub {
			mod.Path = filepath.Join(t.Path, mod.Path)
			modified = append(modified, mod)
		}
		for _, untrack := range untrackedSub {
			untrack.Path = filepath.Join(t.Path, untrack.Path)
			untracked = append(untracked, untrack)
		}
		for _, del := range deletedSub {
			del.Path = filepath.Join(t.Path, del.Path)
			deleted = append(deleted, del)
		}
	}

	return
}
