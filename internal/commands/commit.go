package commands

import (
	"errors"
	"flag"
	"notgit/internal/commit"
	"notgit/internal/config"
	"notgit/internal/status"
	"notgit/internal/utils"
	"os"
	"path/filepath"
)

func Commit() error {
	wd, err := os.Getwd()
	if err != nil {
		return err
	}

	if !utils.RepoInitialized(wd) {
		return errors.New("not a notgit repository")
	}

	// TODO: make config always take data from local config and write it fomr global config on repo init
	config, err := config.Parse(true)
	defaultAuthor := config["user"]["name"] + " <" + config["user"]["email"] + ">"

	var message, author string
	var amend, allowEmpty bool

	fs := flag.NewFlagSet("commit", flag.ExitOnError)

	fs.StringVar(&message, "m", "", "commit message")
	fs.StringVar(&author, "author", defaultAuthor, "commit author")
	fs.BoolVar(&amend, "amend", false, "amend previous commit")
	fs.BoolVar(&allowEmpty, "allow-empty", false, "allow empty commit")

	fs.Parse(os.Args[2:])

	if message == "" {
		return errors.New("commit message is required")
	}
	if author == "" {
		return errors.New("author is required")
	}

	if amend {
		c := commit.ParseHead()
		if c == nil {
			return errors.New("nothing to amend.")
		}
		c.Author = author
		c.Message = message
		return c.Write()
	}

	var parents []string
	head, err := os.ReadFile(filepath.Join(wd, ".notgit", "HEAD"))
	if string(head) != "" {
		parents = append(parents, string(head))
	}

	worktreeAndIndexDiff, indexAndHeadDiff := status.GetRepoStatus()
	if len(worktreeAndIndexDiff)+len(indexAndHeadDiff) == 0 && !allowEmpty {
		return errors.New("nothing to commit, working tree clean")
	}

	c := commit.NewCommit(message, author, parents)
	if c == nil {
		return errors.New("commit creation failed")
	}

	return c.Write()
}
