package tree

import (
	"fmt"
	"maps"
	"notgit/internal/blob"
	"notgit/internal/indexfile"
	"os"
	"path/filepath"
	"slices"
	"strings"
)

var dirsMap = map[string][]blob.Blob{}

// Returns tree with all staged files
func Root() Tree {
	index, err := indexfile.Parse()
	if err != nil {
		return Tree{}
	}

	for _, staged := range index {
		dir := filepath.Dir(staged.Path)
		staged.Path = filepath.Base(staged.Path)
		dirsMap[dir] = append(dirsMap[dir], staged)
	}

	root, err := create(".")
	if err != nil {
		return Tree{}
	}

	return root
}

func (t *Tree) Add(path, fullPath string) error {
	info, err := os.Stat(fullPath)
	if os.IsNotExist(err) {
		return err
	}

	parts := strings.Split(path, "/")

	if len(parts) == 1 {
		if info.IsDir() {
			var subtree *Tree

			// it checks if subtree already exists
			for _, sub := range t.SubTrees {
				if sub.Path == parts[0] {
					subtree = sub
					break
				}
			}

			// if subtree doesn't exist, creates it
			if subtree == nil {
				subtree = &Tree{
					Path:       parts[0],
					Permission: fmt.Sprintf("%o", info.Mode().Perm()),
				}
				t.SubTrees = append(t.SubTrees, subtree)
			}

			entries, err := os.ReadDir(fullPath)
			if err != nil {
				return err
			}

			for _, entry := range entries {
				err := subtree.Add(entry.Name(), filepath.Join(fullPath, entry.Name()))
				if err != nil {
					return err
				}
			}
			return nil
		}

		b, err := blob.Create(fullPath)
		if err != nil {
			return err
		}

		t.Blobs = append(t.Blobs, b)
		return nil
	}

	for _, subtree := range t.SubTrees {
		if subtree.Path == parts[0] {
			err := subtree.Add(filepath.Join(parts[1:]...), fullPath)
			if err != nil {
				return err
			}

			break
		}
	}

	subtree := Tree{
		Path:       parts[0],
		Permission: fmt.Sprintf("%o", info.Mode().Perm()),
	}

	err = subtree.Add(filepath.Join(parts[1:]...), fullPath)
	if err != nil {
		return err
	}

	t.SubTrees = append(t.SubTrees, &subtree)

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
