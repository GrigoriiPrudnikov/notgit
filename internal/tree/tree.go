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
	"strings"
)

// Returns tree with all staged files
func Staged() *Tree {
	index, err := indexfile.Parse()
	if err != nil {
		return nil
	}

	files := map[string][]blob.Blob{}

	for _, staged := range index {
		dir := filepath.Dir(staged.Path)
		staged.Path = filepath.Base(staged.Path)
		files[dir] = append(files[dir], staged)
	}

	root, err := create(".", files)
	if err != nil {
		return nil
	}

	return root
}

func Root() *Tree {
	files := map[string][]blob.Blob{}

	files, err := getAllFiles()
	if err != nil {
		return nil
	}

	root, err := create(".", files)
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

func create(path string, files map[string][]blob.Blob) (*Tree, error) {
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		return nil, err
	}

	root := Tree{
		Path:       filepath.Base(path),
		Permission: fmt.Sprintf("%o", info.Mode().Perm()),
	}

	blobs := files[path]
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

		for k := range maps.Keys(files) {
			subdirs = append(subdirs, k)
		}

		hasSubdir := false
		for path := range files {
			if strings.HasPrefix(path, childPath+string(os.PathSeparator)) {
				hasSubdir = true
				break
			}
		}

		if hasSubdir || files[childPath] != nil {
			subtree, err := create(childPath, files)
			if errors.Is(err, os.ErrNotExist) {
				continue
			}
			if err != nil {
				return nil, err
			}

			root.SubTrees = append(root.SubTrees, subtree)
		}
	}

	return &root, err
}

func getAllFiles() (map[string][]blob.Blob, error) {
	files := map[string][]blob.Blob{}

	err := filepath.WalkDir(".", func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if utils.Ignored(path) {
			if d.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}
		if d.IsDir() {
			return nil
		}

		fmt.Println(path)
		dir := filepath.Dir(path)
		b, err := blob.NewBlob(path)

		files[dir] = append(files[dir], b)

		return nil
	})
	return files, err
}
