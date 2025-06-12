package tree

import (
	"errors"
	"fmt"
	"maps"
	"notgit/internal/blob"
	"notgit/internal/utils"
	"os"
	"path/filepath"
	"strings"
)

type Tree struct {
	SubTrees map[string]*Tree
	Blobs    []blob.Blob
}

func NewTree() *Tree {
	return &Tree{SubTrees: map[string]*Tree{}}
}

// Returns tree with all staged files
func Staged(index []blob.Blob) *Tree {
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
func (t *Tree) Print(indent, treePath string) {
	fmt.Printf("%s- [Tree] %s (%s)\n", indent, treePath, t.Hash)

	for _, b := range t.Blobs {
		fmt.Printf("%s  • [Blob] %s (%s)\n", indent, b.Path, b.Hash)
	}

	for subpath, subtree := range t.SubTrees {
		subtree.Print(indent+"  ", subpath)
	}
}

func create(path string, files map[string][]blob.Blob) (*Tree, error) {
	root := NewTree()

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

			root.SubTrees[child.Name()] = subtree
		}
	}

	return root, err
}

func (t *Tree) Add(path, fullPath string) error {
	if utils.Ignored(fullPath) {
		return nil
	}

	info, err := os.Stat(fullPath)
	if os.IsNotExist(err) {
		return err
	}

	parts := strings.Split(path, "/")

	if len(parts) == 1 {
		if info.IsDir() {
			var subtree *Tree

			// it checks if subtree already exists
			for path, sub := range t.SubTrees {
				if path == parts[0] {
					subtree = sub
					break
				}
			}

			// if subtree doesn't exist, creates it
			if subtree == nil {
				subtree = &Tree{}
				if t.SubTrees == nil {
					t.SubTrees = map[string]*Tree{}
				}
				t.SubTrees[parts[0]] = subtree
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

		newBlob, err := blob.NewBlob(fullPath)
		if err != nil {
			return err
		}

		for i, blob := range t.Blobs {
			if blob.Path == newBlob.Path {
				t.Blobs[i].Content = newBlob.Content
				t.Blobs[i].Hash = newBlob.Hash
				return nil
			}
		}

		t.Blobs = append(t.Blobs, newBlob)
		return nil
	}

	for path, subtree := range t.SubTrees {
		if path == parts[0] {
			err := subtree.Add(filepath.Join(parts[1:]...), fullPath)
			if err != nil {
				return err
			}

			return nil
		}
	}

	subtree := Tree{}

	t.SubTrees[parts[0]] = &subtree

	return subtree.Add(filepath.Join(parts[1:]...), fullPath)
}

func (t *Tree) Hash() string {
	content := []byte{}

	for _, blob := range t.Blobs {
		content = append(content, []byte("blob "+blob.Hash+" "+blob.Path+"\n")...)
	}
	for path, subtree := range t.SubTrees {
		content = append(content, []byte("blob "+subtree.Hash()+" "+path+"\n")...)
	}

	hex := utils.Hash("tree", content)
	return hex
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

		dir := filepath.Dir(path)
		b, err := blob.NewBlob(path)

		files[dir] = append(files[dir], b)

		return nil
	})
	return files, err
}
