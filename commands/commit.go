package commands

import (
	"errors"
	"flag"
	"notgit/internal/commit"
	"notgit/internal/config"
	"notgit/utils"
	"os"
)

func Commit() error {
	wd, err := os.Getwd()
	if err != nil {
		return err
	}

	if !utils.RepoInitialized(wd) {
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

	c := commit.NewCommit(message, author, nil)
	if c == nil {
		return errors.New("commit creation failed")
	}

	return c.Write()
}
