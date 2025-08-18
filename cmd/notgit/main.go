package main

import (
	"fmt"
	"notgit/internal/commands"
	"os"
)

var command = map[string]func() error{
	"add":     commands.Add,
	"commit":  commands.Commit,
	"config":  commands.Config,
	"init":    commands.Init,
	"log":     commands.Log,
	"status":  commands.Status,
	"version": commands.Version,
}

func main() {
	args := os.Args
	if len(args) == 1 {
		// print help
		return
	}

	if len(args) == 2 && (args[1] == "-v" || args[1] == "--version") {
		command["version"]()
		return
	}

	execute, exists := command[os.Args[1]]
	if !exists {
		// TODO: add help like this:
		// The most similar commands are
		//    diff
		//    fsck

		fmt.Printf("notgit: %s is not a git command. See 'notgit --help'.\n", os.Args[1])
		return
	}

	err := execute()
	if err != nil {
		fmt.Println(err)
	}
}
