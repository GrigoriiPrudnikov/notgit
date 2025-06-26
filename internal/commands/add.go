package commands

import (
	"errors"
	"flag"
	"notgit/internal/blob"
	"notgit/internal/indexfile"
	"notgit/internal/tree"
	"notgit/internal/utils"
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
		return errors.New("no path provided")
	}

	index, err := indexfile.Parse()
	if err != nil {
		return err
	}
	root := tree.Staged(index)

	if all || slices.Contains(args, ".") {
		root = tree.Staged(checkForDeleted(index))
		dir, err := os.ReadDir(wd)
		if err != nil {
			return err
		}
		for _, child := range dir {
			root.Add(child.Name(), child.Name())
		}
	}

	for _, arg := range args {
		arg = filepath.Clean(filepath.ToSlash(arg))

		if arg == "." {
			dir, err := os.ReadDir(wd)
			if err != nil {
				return err
			}
			for _, child := range dir {
				root.Add(child.Name(), child.Name())
			}
		}

		if !utils.InWorkingDirectory(arg) {
			return errors.New("'" + arg + "' is not in the working directory at '" + wd + "'")
		}

		err := root.Add(arg, arg)
		if err != nil {
			return err
		}
	}

	err = root.Write()
	if err != nil {
		return err
	}

	return root.WriteIndex()
}

func checkForDeleted(index []blob.Blob) []blob.Blob {
	var result []blob.Blob

	for _, blob := range index {
		_, err := os.Stat(blob.Path)
		if !os.IsNotExist(err) {
			result = append(result, blob)
		}
	}
	return result
}
