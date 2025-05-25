package tree

import (
	"fmt"
	"notgit/internal/blob"
	"notgit/internal/object"
	"notgit/utils"
	"os"
	"path/filepath"
)

func (t *Tree) Write() error {
	content := []byte{}

	for _, blob := range t.Blobs {
		err := blob.Write()
		if err != nil {
			return err
		}

		content = append(content, []byte("blob "+blob.Hash+" "+blob.Path+"\n")...)
	}

	for _, subtree := range t.SubTrees {
		err := subtree.Write()
		if err != nil {
			return err
		}

		content = append(content, []byte("tree "+subtree.Hash()+" "+subtree.Path+"\n")...)
	}

	header := fmt.Sprintf("tree %d\x00\n", len(content))
	compressed := utils.Compress(header, content)

	return object.Write(t.Hash(), compressed)
}

func (t *Tree) WriteIndex() error {
	content := []byte{}

	for _, entry := range t.getEntries() {
		path := entry.path
		blob := entry.blob
		content = append(content, []byte(blob.Hash+" "+path+"\n")...)
	}

	wd, err := os.Getwd()
	if err != nil {
		return err
	}
	index := filepath.Join(wd, ".notgit", "index")

	err = os.WriteFile(index, content, 0644)

	return nil
}

type indexEntry struct {
	path string
	blob blob.Blob
}

func (t *Tree) getEntries() []indexEntry {
	entries := []indexEntry{}

	for _, blob := range t.Blobs {
		entries = append(entries, indexEntry{
			path: filepath.Join(t.Path, blob.Path),
			blob: blob,
		})
	}

	for _, subtree := range t.SubTrees {
		for _, entry := range subtree.getEntries() {
			entries = append(entries, indexEntry{
				path: filepath.Join(t.Path, entry.path),
				blob: entry.blob,
			})
		}
	}

	return entries
}
