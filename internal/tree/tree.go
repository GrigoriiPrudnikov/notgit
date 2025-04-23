package tree

import (
	"errors"
	"fmt"
	"maps"
	"notgit/internal/blob"
	"notgit/internal/indexfile"
	"os"
	"path/filepath"
	"strings"
)

var dirsMap = map[string][]blob.Blob{}

// Returns tree with all staged files
func Root() *Tree {
	index, err := indexfile.Parse()
	if err != nil {
		return nil
	}

	for _, staged := range index {
		dir := filepath.Dir(staged.Path)
		staged.Path = filepath.Base(staged.Path)
		dirsMap[dir] = append(dirsMap[dir], staged)
	}

	root, err := create(".")
	if err != nil {
		return nil
	}

	return root
}

// For debug, remove later
func (t *Tree) Print(indent string) {
	fmt.Printf("%s- [Tree] %s %s (%s)\n", indent, t.Permission, t.Path, t.Hash)

	for _, b := range t.Blobs {
		fmt.Printf("%s  • [Blob] %s %s (%s)\n", indent, b.Permission, b.Path, b.Hash)
	}

	for _, subtree := range t.SubTrees {
		subtree.Print(indent + "  ")
	}
}

func create(path string) (*Tree, error) {
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		return nil, err
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
		return nil, err
	}

	for _, child := range children {
		if !child.IsDir() {
			continue
		}

		childPath := filepath.Join(path, child.Name())
		childPath = filepath.Clean(childPath)
		subdirs := []string{}

		for k := range maps.Keys(dirsMap) {
			subdirs = append(subdirs, k)
		}

		hasSubdir := false
		for path := range dirsMap {
			if strings.HasPrefix(path, childPath+string(os.PathSeparator)) {
				hasSubdir = true
				break
			}
		}

		if hasSubdir || dirsMap[childPath] != nil {
			subtree, err := create(childPath)
			if errors.Is(err, os.ErrNotExist) {
				continue
			}
			if err != nil {
				return nil, err
			}

			root.SubTrees = append(root.SubTrees, subtree)
		}
	}

	Hash(&root)

	return &root, err
}
