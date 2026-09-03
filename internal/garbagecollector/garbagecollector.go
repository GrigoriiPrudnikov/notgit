package garbagecollector

import (
	"notgit/internal/commit"
	"notgit/internal/indexfile"
	"notgit/internal/tree"
	"os"
	"path/filepath"
)

// Deletes all unused objects
func CollectGarbage(wd string) error {
	objectsDir := filepath.Join(wd, ".notgit", "objects")

	seenObjects := map[string]bool{}

	stagedFiles, err := indexfile.Parse(wd)
	if err != nil {
		return err
	}

	for _, hash := range stagedFiles {
		seenObjects[hash] = true
	}

	currentCommit := commit.ParseHead()
	if currentCommit == nil {
		return nil
	}

	seenObjects[currentCommit.Hash()] = true

	for currentCommit != nil {
		err := checkTree(currentCommit.Tree, seenObjects)
		if err != nil {
			return err
		}

		if currentCommit.Parent == "" {
			currentCommit = nil
		} else {
			currentCommit = commit.Parse(currentCommit.Parent)
		}
	}

	dirs, err := os.ReadDir(objectsDir)
	if err != nil {
		return err
	}

	for _, dir := range dirs {
		info, err := os.Stat(filepath.Join(objectsDir, dir.Name()))
		if err != nil {
			return err
		}
		if !info.IsDir() {
			continue
		}

		files, err := os.ReadDir(filepath.Join(objectsDir, dir.Name()))
		if err != nil {
			return err
		}

		for _, file := range files {
			fullHash := dir.Name() + file.Name()

			if _, ok := seenObjects[fullHash]; !ok {
				err := os.Remove(filepath.Join(objectsDir, filepath.Join(dir.Name(), file.Name())))
				if err != nil {
					return err
				}
			}
		}

		// Delete directory if it is empty
		files, err = os.ReadDir(filepath.Join(objectsDir, dir.Name()))
		if err != nil {
			return err
		}

		if len(files) == 0 {
			os.Remove(filepath.Join(objectsDir, dir.Name()))
		}

	}

	return nil
}

func checkTree(tree *tree.Tree, seenObjects map[string]bool) error {
	seenObjects[tree.Hash()] = true

	for _, subtree := range tree.SubTrees {
		err := checkTree(subtree, seenObjects)
		if err != nil {
			return err
		}
	}

	for _, blob := range tree.Blobs {
		seenObjects[blob] = true
	}

	return nil
}
