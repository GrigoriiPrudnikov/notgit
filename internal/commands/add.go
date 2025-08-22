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

func Add(wd string) error {
	if !utils.RepoInitialized(wd) {
		return errors.New("not a notgit repository")
	}

	var all, force, update bool

	fs := flag.NewFlagSet("add", flag.ExitOnError)

	fs.BoolVar(&all, "all", false, "add all")
	fs.BoolVar(&all, "a", false, "add all")

	fs.BoolVar(&force, "force", false, "force")
	fs.BoolVar(&force, "f", false, "force")

	fs.BoolVar(&update, "update", false, "update")
	fs.BoolVar(&update, "u", false, "update")

	fs.Parse(os.Args[2:])
	args := fs.Args()

	if slices.Contains(args, ".") || update {
		all = true
	}

	index, err := indexfile.Parse()
	if err != nil {
		return err
	}

	if update {
		err := addPathToIndex(".", true, &index)
		if err != nil {
			return err
		}
	}

	if all {
		args = []string{"."}
	}

	for _, path := range args {
		err := addPathToIndex(path, update, &index)
		if err != nil {
			return err
		}
	}

	if all || update {
		for path := range index {
			_, err := os.Stat(path)

			if os.IsNotExist(err) {
				delete(index, path)
			}
		}
	}

	return indexfile.Write(index)
}

func addPathToIndex(path string, update bool, index *map[string]string) error {
	info, err := os.Stat(path)
	if err != nil {
		return err
	}

	if info.IsDir() {
		if utils.Ignored(path) {
			return nil
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

			if update {
				// check if path is already in index
				if _, ok := (*index)[path]; !ok {
					continue
				}
			}

			b, err := blob.NewBlob(path)
			if err != nil {
				return err
			}
			(*index)[path] = b.Hash()
		}
		return nil
	}

	b, err := blob.NewBlob(path)
	if err != nil {
		return err
	}

	(*index)[path] = b.Hash()

	return nil
}
