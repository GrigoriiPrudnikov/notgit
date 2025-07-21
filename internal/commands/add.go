package commands

import (
	"errors"
	"flag"
	"notgit/internal/tree"
	"notgit/internal/utils"
	"os"
	"path/filepath"
)

func Add() error {
	wd, err := os.Getwd()
	if err != nil {
		return err
	}

	if !utils.RepoInitialized(wd) {
		return errors.New("not a notgit repository")
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

	root, err := tree.LoadStaged()
	if err != nil {
		return err
	}

	for _, path := range args {
		if !utils.InWorkingDirectory(path) {
			return errors.New("path is not in working directory")
		}

		if filepath.Clean(path) == "." {
			entries, err := os.ReadDir(path)
			if err != nil {
				return err
			}

			for _, entry := range entries {
				root.Add(entry.Name())
			}
		}

		err := root.Add(path)
		if err != nil {
			return err
		}
	}

	utils.PrintStruct(root)

	return root.WriteIndex()
}
