package tree

import (
	"fmt"
	"maps"
	"notgit/internal/blob"
	"notgit/internal/indexfile"
	"os"
	"path/filepath"
	"slices"
)

var dirsMap = map[string][]blob.Blob{}

// Returns tree with all staged files
func Root() Tree {
	index, err := indexfile.Parse()
	if err != nil {
		return Tree{}
	}

	for _, staged := range index {
		staged.Path = filepath.Base(staged.Path)
		dirsMap[staged.Path] = append(dirsMap[staged.Path], staged)
	}

	root, err := create("")
	if err != nil {
		return Tree{}
	}

	return root
}

func (t *Tree) Write() error {
	return nil
}

func create(path string) (Tree, error) {
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		return Tree{}, err
	}

	root := Tree{
		Path:       filepath.Base(path),
		Permission: fmt.Sprintf("%o", info.Mode().Perm()),
	}

	blobs := dirsMap[path]
	for _, blob := range blobs {
		root.Blobs = append(root.Blobs, blob)
	}

	children, err := os.ReadDir(path)
	if err != nil {
		return Tree{}, err
	}

	for _, child := range children {
		if !child.IsDir() {
			continue
		}

		childPath := filepath.Join(path, child.Name())
		subdirs := []string{}

		for k := range maps.Keys(dirsMap) {
			subdirs = append(subdirs, k)
		}

		if slices.Contains(subdirs, childPath) {
			subtree, err := create(childPath)
			if err != nil {
				return Tree{}, err
			}

			root.SubTrees = append(root.SubTrees, &subtree)
		}
	}

	Hash(&root)

	return root, err
}
