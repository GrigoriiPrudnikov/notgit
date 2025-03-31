package commands

import (
	"errors"
	"flag"
	"fmt"
	"notgit/internal/blob"
	"notgit/internal/indexfile"
	"notgit/utils"
	"os"
	"path/filepath"
	"slices"
)

func Add() error {
	wd, err := os.Getwd()
	if err != nil {
		return err
	}

	notgitDir := filepath.Join(wd, ".notgit")
	if _, err := os.Stat(notgitDir); os.IsNotExist(err) {
		// TODO: add handling for parent directories (fatal: not a git repository (or any of the parent directories): .git)
		return errors.New("not a git repository")
	}

	var all, force bool

	fs := flag.NewFlagSet("add", flag.ExitOnError)

	fs.BoolVar(&all, "all", false, "add all")
	fs.BoolVar(&all, "a", false, "add all")

	fs.BoolVar(&force, "force", false, "force")
	fs.BoolVar(&force, "f", false, "force")

	// fs.BoolVar(&update, "update", false, "update")
	// fs.BoolVar(&update, "u", false, "update")

	fs.Parse(os.Args[2:])
	args := fs.Args()

	if len(args) == 0 {
		return errors.New("no arguments")
	}

	for _, arg := range args {
		err := add(arg, force)

		if err != nil {
			return err
		}
	}

	return nil
}

func add(path string, force bool) error {
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		return err
	}

	if utils.Ignored(path) && !force {
		fmt.Println("ignored", path)
		return nil
	}

	if info.IsDir() {
		return addDir(path, force)
	}

	return addFile(path, force)
}

func addDir(path string, force bool) error {
	entries, err := os.ReadDir(path)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		entryPath := filepath.Join(path, entry.Name())
		add(entryPath, force)
	}

	return nil
}

func addFile(path string, force bool) error {
	wd, err := os.Getwd()
	if err != nil {
		return err
	}

	b, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	err = blob.Create(b)
	if err != nil {
		return err
	}

	indexPath := filepath.Join(wd, ".notgit", "index")

	_, err = os.Stat(indexPath)
	if os.IsNotExist(err) {
		err = os.WriteFile(indexPath, []byte(""), 0644)
		if err != nil {
			return err
		}
	}

	hash := blob.Hash(b)

	stagedFiles, err := indexfile.Parse()
	if err != nil {
		return err
	}

	stagedFile := indexfile.StagedFile{
		Permission: "100644",
		Hash:       hash,
		Name:       path,
	}

	// check for missing files and update if exists
	for i, file := range stagedFiles {
		if _, err := os.Stat(file.Name); os.IsNotExist(err) {
			stagedFiles = slices.Delete(stagedFiles, i, i+1)
			continue
		}

		if file.Name == path {
			stagedFiles[i] = stagedFile
			err = indexfile.Set(stagedFiles)
			return err
		}
	}

	stagedFiles = append(stagedFiles, stagedFile)
	err = indexfile.Set(stagedFiles)

	return err
}
