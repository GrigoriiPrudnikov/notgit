package status

import (
	"notgit/internal/blob"
	"notgit/internal/commit"
	"notgit/internal/tree"
	"path/filepath"
)

func GetStatus() (staged, modified, untracked []string) {
	all := tree.Root()
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

	modifiedBlobs, untrackedBlobs := getModifiedAndUntracked(all, stagedTree)
	stagedBlobs := getStaged(headTree, stagedTree)

	staged = extractPaths(stagedBlobs)
	untracked = extractPaths(untrackedBlobs)
	modified = filterModified(modifiedBlobs, stagedBlobs)

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

func getStaged(head, stagedTree *tree.Tree) (staged []blob.Blob) {
	diff := compare(stagedTree, head)
	for _, d := range diff {
		staged = append(staged, d)
	}

	for _, sub := range stagedTree.SubTrees {
		var headSubTree *tree.Tree
		if head != nil {
			headSubTree = findTree(head.SubTrees, sub.Path)
		}
		stagedSub := getStaged(headSubTree, sub)

		for _, s := range stagedSub {
			s.Path = filepath.Join(sub.Path, s.Path)
			staged = append(staged, s)
		}
	}

	return
}
