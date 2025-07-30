package tree

import (
	"fmt"
	"maps"
	"notgit/internal/blob"
	"notgit/internal/object"
	"notgit/internal/utils"
	"os"
	"path/filepath"
)

func (t *Tree) Write() error {
	for _, subtree := range t.SubTrees {
		err := subtree.Write()
		if err != nil {
			return err
		}
	}
	for _, hash := range t.Blobs {
		_, content, err := object.Parse(hash)
		if err != nil {
			return err
		}

		b := &blob.Blob{
			Content: content,
		}
		if b.Write() != nil {
			return err
		}
	}

	content, err := t.GetContent()
	if err != nil {
		return err
	}

	header := fmt.Sprintf("tree %d\x00\n", len(content))
	compressed := utils.Compress(header, content)

	return object.Write(t.Hash(), compressed)
}

func (t *Tree) WriteIndex() error {
	content := []byte{}

	entries := t.getEntries()
	for _, path := range utils.GetSortedKeys(entries) {
		content = append(content, []byte(path+" "+entries[path]+"\n")...)
	}

	wd, err := os.Getwd()
	if err != nil {
		return err
	}
	index := filepath.Join(wd, ".notgit", "index")

	return os.WriteFile(index, content, 0644)
}

func (t *Tree) getEntries() map[string]string {
	entries := make(map[string]string)

	for path, hash := range t.Blobs {
		entries[filepath.Clean(filepath.Join(t.Path, path))] = hash
	}

	for _, subtree := range t.SubTrees {
		maps.Copy(entries, subtree.getEntries())
	}

	return entries
}

func (t Tree) GetContent() ([]byte, error) {
	content := []byte{}

	subtreesPaths := utils.GetSortedKeys(t.SubTrees)
	filesPaths := utils.GetSortedKeys(t.Blobs)

	for _, treePath := range subtreesPaths {
		subtree := t.SubTrees[treePath]
		line := "tree " + treePath + " " + subtree.Hash() + "\n"
		content = append(content, []byte(line)...)
	}

	for _, filePath := range filesPaths {
		hash := t.Blobs[filePath]
		line := "blob " + filePath + " " + hash + "\n"
		content = append(content, []byte(line)...)
	}

	return content, nil
}
