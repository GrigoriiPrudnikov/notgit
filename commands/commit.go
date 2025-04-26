package commands

import (
	"errors"
	"flag"
	"notgit/internal/commit"
	"os"
)

func Commit() error {
	var message, author string

	fs := flag.NewFlagSet("commit", flag.ExitOnError)

	fs.StringVar(&message, "m", "", "commit message")
	fs.StringVar(&message, "message", "", "commit message")
	fs.StringVar(&author, "author", "", "author")

	fs.Parse(os.Args[2:])

	commit := commit.NewCommit(message, author, nil)
	if commit == nil {
		return errors.New("commit creation failed")
	}

	err := commit.Write()

	return err
}
