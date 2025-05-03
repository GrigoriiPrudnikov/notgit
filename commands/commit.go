package commands

import (
	"errors"
	"flag"
	"fmt"
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

	// TODO: make config always take data from local config and write it fomr global config on repo init
	config, err := config.Parse(true)
	defaultAuthor := config["user"]["name"] + " <" + config["user"]["email"] + ">"

	var message, author string

	fs := flag.NewFlagSet("commit", flag.ExitOnError)

	fs.StringVar(&message, "m", "", "commit message")
	fs.StringVar(&message, "message", "", "commit message")
	fs.StringVar(&author, "author", defaultAuthor, "commit author")

	fs.Parse(os.Args[2:])

	if message == "" {
		return errors.New("commit message is required")
	}
	if author == "" {
		return errors.New("author is required")
	}

	c := commit.NewCommit(message, author, nil)
	if c == nil {
		return errors.New("commit creation failed")
	}

	parsed, err := commit.ParseHead()
	if err != nil {
		return err
	}
	fmt.Println("parsed:", *parsed)

	return nil // c.Write()
}
