package commands

import (
	"errors"
	"flag"
	"os"
)

func Commit() error {
	var message, author string

	fs := flag.NewFlagSet("commit", flag.ExitOnError)

	fs.StringVar(&message, "m", "", "commit message")
	fs.StringVar(&message, "message", "", "commit message")
	fs.StringVar(&author, "author", "", "author")

	fs.Parse(os.Args[2:])
	args := fs.Args()

	if len(args) == 0 {
		return errors.New("no arguments")
	}

	return nil
}
