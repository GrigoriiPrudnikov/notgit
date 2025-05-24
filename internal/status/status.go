package status

import (
	"notgit/internal/blob"
	"notgit/internal/commit"
	"notgit/internal/tree"
	"path/filepath"
)

// TODO: split to GetStaged and GetUnstaged
func GetStatus() (modifiedStaged, untrackedStaged, modified, untracked []string) {
	all := tree.Root()
	stagedTree := tree.Staged()
	head := commit.ParseHead()

	var headTree *tree.Tree
	var err error

	if head != nil {
		headTree, err = tree.Parse(head.Tree)
		if err != nil {
			return nil, nil, nil, nil
		}
	}

	modifiedBlobs, untrackedBlobs := getModifiedAndUntracked(all, stagedTree)
	modifiedStagedBlobs, untrackedStagedBlobs := getModifiedAndUntracked(stagedTree, headTree)

	untracked = extractPaths(untrackedBlobs)
	modified = extractPaths(modifiedBlobs)
	modifiedStaged = extractPaths(modifiedStagedBlobs)
	untrackedStaged = extractPaths(untrackedStagedBlobs)

	return
}

func GetStaged() (modified, added []string) {
	stagedTree := tree.Staged()
	head := commit.ParseHead()

	var headTree *tree.Tree
	var err error

	if head != nil {
		headTree, err = tree.Parse(head.Tree)
		if err != nil {
			return nil, nil
		}
	}

	modifiedStagedBlobs, addedBlobs := getModifiedAndUntracked(stagedTree, headTree)

	modified = extractPaths(modifiedStagedBlobs)
	added = extractPaths(addedBlobs)

	return
}

func GetUnstaged() (modified, untracked []string) {
	all := tree.Root()
	stagedTree := tree.Staged()

	modifiedBlobs, untrackedBlobs := getModifiedAndUntracked(all, stagedTree)

	modified = extractPaths(modifiedBlobs)
	untracked = extractPaths(untrackedBlobs)

	return
}

func getModifiedAndUntracked(all, staged *tree.Tree) (modified, untracked []blob.Blob) {
	diff := compare(all, staged)

	for _, b := range diff {
		if staged == nil || findBlob(staged.Blobs, b.Path) == nil {
			untracked = append(untracked, b)
		} else {
			modified = append(modified, b)
		}
	}

	for _, t := range all.SubTrees {
		var found *tree.Tree
		if staged != nil {
			found = findTree(staged.SubTrees, t.Path)
		}
		modifiedSub, untrackedSub := getModifiedAndUntracked(t, found)

		for _, mod := range modifiedSub {
			mod.Path = filepath.Join(t.Path, mod.Path)
			modified = append(modified, mod)
		}
		for _, untrack := range untrackedSub {
			untrack.Path = filepath.Join(t.Path, untrack.Path)
			untracked = append(untracked, untrack)
		}
	}

	return
}
