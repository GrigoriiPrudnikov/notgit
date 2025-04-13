package commands

import (
	"errors"
	"flag"
	"notgit/internal/blob"
	"notgit/internal/indexfile"
	"notgit/utils"
	"os"
	"path/filepath"
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

	if all {
		return add(".", force)
	}

	for _, arg := range args {
		err := add(arg, force)

		if err != nil {
			return err
		}
	}

	return nil
}

// TODO: rewrite to trees
func add(path string, force bool) error {
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		return err
	}

	if utils.Ignored(path) && !force {
		return nil
	}

	if info.IsDir() {
		children, err := os.ReadDir(path)
		if err != nil {
			return err
		}

		for _, child := range children {
			childPath := filepath.Join(path, child.Name())

			err := add(childPath, force)
			if err != nil {
				return err
			}
		}

		return nil
	}

	return addFile(path)
}

func addFile(path string) error {
	wd, err := os.Getwd()
	if err != nil {
		return err
	}

	b, err := blob.Create(path)
	b.Path = path
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

	stagedFiles, err := indexfile.Parse()
	if err != nil {
		return err
	}

	var updated []blob.Blob

	for _, staged := range stagedFiles {
		if _, err := os.Stat(filepath.Join(wd, staged.Path)); os.IsNotExist(err) {
			continue
		}

		if staged.Path == path {
			continue
		}

		updated = append(updated, staged)
	}

	updated = append(updated, b)

	err = indexfile.Write(updated)

	return err
}
