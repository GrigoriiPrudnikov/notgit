package commands

import (
	"errors"
	"flag"
	"notgit/internal/commit"
	"notgit/internal/config"
	"os"
	"path/filepath"
)

func Commit() error {
	wd, err := os.Getwd()
	if err != nil {
		return err
	}

	notgitDir := filepath.Join(wd, ".notgit")
	if _, err := os.Stat(notgitDir); os.IsNotExist(err) {
		// TODO: add handling for parent directories (fatal: not a git repository (or any of the parent directories): .git)
		return errors.New("not a git repository")
	}

	config, err := config.Parse(true)

	var message, author string

	fs := flag.NewFlagSet("commit", flag.ExitOnError)

	fs.StringVar(&message, "m", "", "commit message")
	fs.StringVar(&message, "message", "", "commit message")
	fs.StringVar(&author, "author", config["user"]["name"], "author")

	fs.Parse(os.Args[2:])

	commit := commit.NewCommit(message, author, nil)
	if commit == nil {
		return errors.New("commit creation failed")
	}

	err = commit.Write()

	return err
}
