// Package tree contains the tree-related functions.
package tree

import (
	"os"
	"path/filepath"
	"strings"

	"notgit/internal/blob"
	"notgit/internal/indexfile"
	"notgit/internal/object"
	"notgit/internal/utils"
)

type Tree struct {
	Path     string
	SubTrees map[string]*Tree
	Blobs    map[string]string // key is relative path, value is hash
}

func NewTree(path string) *Tree {
	return &Tree{
		Path:     path,
		SubTrees: make(map[string]*Tree),
		Blobs:    make(map[string]string),
	}
}

func (t *Tree) Hash() (string, error) {
	content, err := t.GetContent()
	if err != nil {
		return "", err
	}

	return utils.Hash("tree", content), nil
}

func (t Tree) BasePath() string {
	return filepath.Base(t.Path)
}

// Add records a path in the tree and writes its blob via the provided Store,
// so callers no longer have to thread the working directory into tree logic.
func (t *Tree) Add(path string, store object.Store) error {
	if utils.Ignored(path) {
		return nil
	}

	info, err := os.Stat(path)
	if err != nil {
		return err
	}
	if os.IsNotExist(err) {
		return nil
	}

	if info.IsDir() {
		dir, error := os.ReadDir(path)
		if error != nil {
			return error
		}
		for _, entry := range dir {
			if error = t.Add(filepath.Join(path, entry.Name()), store); error != nil {
				return err
			}
		}
		return nil
	}

	b, err := blob.NewBlob(path)
	if err != nil {
		return err
	}

	err = t.addFile(path, b.Hash())
	if err != nil {
		return err
	}
	return b.Write(store)
}

// Find returns found hash and flag indicating whether the file was found
func (t Tree) Find(path string) (string, bool) {
	parts := strings.Split(path, string(filepath.Separator))

	if len(parts) == 1 {
		hash, ok := t.Blobs[path]
		return hash, ok
	}

	subdir := parts[0]
	subTree, ok := t.SubTrees[subdir]
	if !ok {
		return "", false
	}

	return subTree.Find(filepath.Join(parts[1:]...))
}

// LoadWorktree walks the working tree, writing blobs through the provided Store
// instead of constructing paths using a wd string directly.
func LoadWorktree(path string, store object.Store) (*Tree, error) {
	if utils.Ignored(path) {
		return nil, nil
	}

	dir, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}

	root := NewTree(path)

	for _, entry := range dir {
		if entry.IsDir() {
			tree, err := LoadWorktree(filepath.Join(path, entry.Name()), store)
			if err != nil {
				return nil, err
			}
			if tree == nil {
				continue
			}
			root.SubTrees[tree.BasePath()] = tree
			continue
		}

		if utils.Ignored(path) {
			continue
		}

		if utils.Ignored(filepath.Join(path, entry.Name())) {
			continue
		}

		b, err := blob.NewBlob(filepath.Join(path, entry.Name()))
		if err != nil {
			return nil, err
		}
		if err := b.Write(store); err != nil {
			return nil, err
		}
		root.Blobs[b.BasePath()] = b.Hash()
	}
	return root, nil
}

func LoadStaged(wd string) (*Tree, error) {
	root := NewTree(".")

	index, err := indexfile.Parse(wd)
	if err != nil {
		println("error parsing index file:", err.Error())
		return nil, err
	}

	for path, hash := range index {
		if err := root.addFile(path, hash); err != nil {
			return nil, err
		}
	}

	return root, nil
}

func (t *Tree) addFile(path, hash string) error {
	if path == filepath.Base(path) {
		t.Blobs[path] = hash
		return nil
	}

	parts := strings.Split(path, string(filepath.Separator))

	subdir := parts[0]
	subTree, ok := t.SubTrees[subdir]
	if !ok {
		subTree = NewTree(filepath.Join(t.Path, subdir))
		t.SubTrees[subdir] = subTree
	}

	return subTree.addFile(filepath.Join(parts[1:]...), hash)
}
