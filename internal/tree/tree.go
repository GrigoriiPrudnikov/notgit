package tree

import (
	"errors"
	"fmt"
	"maps"
	"notgit/internal/blob"
	"notgit/internal/indexfile"
	"notgit/utils"
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

			Hash(subtree)
			return nil
		}

		b, err := blob.NewBlob(fullPath)
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

			Hash(subtree)
			return nil
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

	Hash(&subtree)
	t.SubTrees = append(t.SubTrees, &subtree)

	return nil
}

func (t *Tree) Write() error {
	content := []byte{}

	for _, blob := range t.Blobs {
		err := blob.Write()
		if err != nil {
			return err
		}

		content = append(content, []byte(blob.Permission+" blob "+blob.Hash+" "+blob.Path+"\n")...)
	}

	for _, subtree := range t.SubTrees {
		err := subtree.Write()
		if err != nil {
			return err
		}

		content = append(content, []byte(subtree.Permission+" tree "+subtree.Hash+" "+subtree.Path+"\n")...)
	}

	compressed := utils.Compress(content, "tree")

	wd, err := os.Getwd()
	if err != nil {
		return err
	}
	objects := filepath.Join(wd, ".notgit", "objects")

	hash := t.Hash
	dir := filepath.Join(objects, hash[:2])
	file := filepath.Join(dir, hash[2:])

	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err = os.Mkdir(dir, 0755)
		if err != nil {
			return err
		}
	}

	if _, err := os.Stat(file); os.IsExist(err) {
		return nil
	}

	err = os.WriteFile(file, compressed, 0644)

	return err
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
			if errors.Is(err, os.ErrNotExist) {
				continue
			}
			if err != nil {
				return Tree{}, err
			}

			root.SubTrees = append(root.SubTrees, &subtree)
		}
	}

	Hash(&root)

	return root, err
}
