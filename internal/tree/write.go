package tree

import (
	"fmt"
	"notgit/internal/blob"
	"notgit/internal/object"
	"notgit/internal/utils"
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

		content = append(content, []byte("blob "+blob.Hash()+" "+blob.Path+"\n")...)
	}

	for path, subtree := range t.SubTrees {
		err := subtree.Write()
		if err != nil {
			return err
		}

		content = append(content, []byte("tree "+subtree.Hash()+" "+path+"\n")...)
	}

	header := fmt.Sprintf("tree %d\x00\n", len(content))
	compressed := utils.Compress(header, content)

	return object.Write(t.Hash(), compressed)
}

func (t *Tree) WriteIndex() error {
	content := []byte{}

	for _, entry := range t.getEntries("") {
		fmt.Println(entry.path)
		path := entry.path
		blob := entry.blob
		content = append(content, []byte(blob.Hash()+" "+path+"\n")...)
	}

	wd, err := os.Getwd()
	if err != nil {
		return err
	}
	index := filepath.Join(wd, ".notgit", "index")

	return os.WriteFile(index, content, 0644)
}

type indexEntry struct {
	path string
	blob blob.Blob
}

func (t *Tree) getEntries(path string) []indexEntry {
	entries := []indexEntry{}

	for _, blob := range t.Blobs {
		entries = append(entries, indexEntry{
			path: filepath.Join(path, blob.Path),
			blob: blob,
		})
	}

	for subpath, subtree := range t.SubTrees {
		fullSubPath := filepath.Join(path, subpath)
		entries = append(entries, subtree.getEntries(fullSubPath)...)
	}

	return entries
}
