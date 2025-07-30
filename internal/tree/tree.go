package tree

import (
	"notgit/internal/blob"
	"notgit/internal/indexfile"
	"notgit/internal/utils"
	"os"
	"path/filepath"
	"strings"
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

func (t *Tree) Hash() string {
	content, err := t.GetContent()
	if err != nil {
		return ""
	}

	return utils.Hash("tree", content)
}

func (t Tree) BasePath() string {
	return filepath.Base(t.Path)
}

func (t *Tree) Add(path string) error {
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
		dir, err := os.ReadDir(path)
		if err != nil {
			return err
		}
		for _, entry := range dir {
			if err := t.Add(filepath.Join(path, entry.Name())); err != nil {
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
	return b.Write()
}

// Returns found hash and flag indicating whether the file was found
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

func LoadWorktree(path string) (*Tree, error) {
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
			tree, err := LoadWorktree(filepath.Join(path, entry.Name()))
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
		err = b.Write()
		if err != nil {
			return nil, err
		}
		root.Blobs[b.BasePath()] = b.Hash()
	}
	return root, nil
}

func LoadStaged() (*Tree, error) {
	root := NewTree(".")

	index, err := indexfile.Parse()
	if err != nil {
		println("error parsing index file:", err.Error())
		return nil, err
	}

	for path, hash := range index {
		err = root.addFile(path, hash)
		if err != nil {
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
