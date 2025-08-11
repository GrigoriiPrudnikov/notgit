package commands

import (
	"errors"
	"flag"
	"notgit/internal/blob"
	"notgit/internal/indexfile"
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
		return nil
	}

	if slices.Contains(args, ".") {
		all = true
	}

	index, err := indexfile.Parse()
	if err != nil {
		return err
	}

	for _, path := range args {
		info, err := os.Stat(path)
		if err != nil {
			return err
		}

		if info.IsDir() {
			if utils.Ignored(path) {
				continue
			}

			paths := []string{}
			err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
				if err != nil {
					return err
				}
				if info.IsDir() {
					return nil
				}
				paths = append(paths, path)
				return nil
			})
			if err != nil {
				return err
			}

			for _, path := range paths {
				if utils.Ignored(path) {
					continue
				}

				// todo: here add check for --update flag

				b, err := blob.NewBlob(path)
				if err != nil {
					return err
				}
				index[path] = b.Hash()
			}
			continue
		}

		b, err := blob.NewBlob(path)
		if err != nil {
			return err
		}

		index[path] = b.Hash()
	}

	// here add update logic
	if all {
		for path := range index {
			_, err := os.Stat(path)

			if os.IsNotExist(err) {
				delete(index, path)
			}
		}
	}

	return indexfile.Write(index)
}
