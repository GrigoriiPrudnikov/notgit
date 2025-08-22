package commands

import (
	"errors"
	"flag"
	"notgit/internal/tree"
	"notgit/internal/utils"
	"os"
)

func Reset(wd string) error {
	if !utils.RepoInitialized(wd) {
		return errors.New("not a notgit repository")
	}

	var hard bool

	fs := flag.NewFlagSet("reset", flag.ExitOnError)
	fs.BoolVar(&hard, "hard", false, "Hard reset")

	fs.Parse(os.Args[2:])
	args := fs.Args()

	var paths []string
	for _, a := range args {
		if a[0] == '-' {
			continue
		}
		paths = append(paths, a)
	}

	tree, err := tree.LoadStaged()
	if err != nil {
		return err
	}
	utils.PrintStruct(tree)

	// for _, path := range paths {
	// 	err := tree.Remove(path)
	// 	if err != nil {
	// 		return err
	// 	}
	// }

	return nil
}
