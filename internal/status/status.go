package status

import (
	"fmt"
	"maps"
	"notgit/internal/commit"
	"notgit/internal/tree"
	"path/filepath"
)

const (
	None = iota
	Added
	Modified
	Deleted
)

type Status int

// Returns differences between worktree and index and between index and head
func GetRepoStatus() (map[string]Status, map[string]Status) {
	worktree, err := tree.LoadWorktree(".")
	if err != nil {
		return nil, nil
	}
	staged, err := tree.LoadStaged()
	if err != nil {
		fmt.Println("error loading staged tree:", err)
		return nil, nil
	}
	headCommit := commit.ParseHead()
	head := tree.NewTree(".")
	if headCommit != nil {
		head = headCommit.Tree
	}

	worktreeAndIndexDiff := compareTrees(worktree, staged)
	indexAndHeadDiff := compareTrees(staged, head)

	return worktreeAndIndexDiff, indexAndHeadDiff
}

func compareTrees(tree1, tree2 *tree.Tree) map[string]Status {
	result := map[string]Status{}

	// check for added and modified files
	for path, hash := range tree1.Blobs {
		hash2, found := tree2.Blobs[path]
		fullPath := filepath.Clean(filepath.Join(tree1.Path, path))
		if !found {
			result[fullPath] = Added
			continue
		}

		if hash != hash2 {
			result[fullPath] = Modified
		}
	}

	// check for deleted files
	for path := range tree2.Blobs {
		_, found := tree1.Blobs[path]
		if !found {
			result[path] = Deleted
		}
	}

	// check subtrees
	for path, subtree := range tree1.SubTrees {
		subtree2 := tree2.SubTrees[path]
		if subtree2 == nil {
			subtree2 = tree.NewTree(subtree.Path)
		}

		maps.Copy(result, compareTrees(subtree, subtree2))
	}

	// check for deleted subtrees/files
	for path, subtree := range tree2.SubTrees {
		subtree2 := tree1.SubTrees[path]
		if subtree2 == nil {
			subtree2 = tree.NewTree(subtree.Path)
		}

		maps.Copy(result, compareTrees(subtree, subtree2))
	}

	return result
}
