package commands

import (
	"notgit/internal/repo"
)

type Command struct {
	Description      string
	Usage            string
	AvailableOptions map[string]string
	Run              func(r *repo.Repo, params []string, opts map[string]string) error
}

var Commands = map[string]Command{
	"init": {
		Description:      "Initialize a new repository",
		Usage:            "notgit init [options]",
		AvailableOptions: map[string]string{},
		Run:              Init,
	},
	"add": {
		Description:      "Add files to the staging area",
		Usage:            "notgit add [files] [options]",
		AvailableOptions: map[string]string{},
		Run:              Add,
	},
	"commit": {
		Description: "Record changes to the repository",
		Usage:       "notgit commit -m <message> [options]",
		AvailableOptions: map[string]string{
			"-m, --message": "commit message",
		},
		Run: Commit,
	},
	"log": {
		Description:      "Show the commit history",
		Usage:            "notgit log [options]",
		AvailableOptions: map[string]string{},
		Run:              Log,
	},
	"version": {
		Description:      "Show the version",
		Usage:            "notgit version",
		AvailableOptions: map[string]string{},
		Run:              Version,
	},
	"clean": {
		Description:      "FOR TESTING ONLY",
		Usage:            "notgit clean",
		AvailableOptions: map[string]string{},
		Run:              Clean,
	},
	// "status": {
	// 	Description:      "Show the status of the repository",
	// 	Usage:            "notgit status [options]",
	// 	AvailableOptions: map[string]string{"-s": "short"},
	// 	Run:              Status,
	// },
	// "config": {
	// 	Description:      "Show or set configuration values",
	// 	Usage:            "notgit config [options]",
	// 	AvailableOptions: map[string]string{},
	// 	Run:              Config,
	// },
}
