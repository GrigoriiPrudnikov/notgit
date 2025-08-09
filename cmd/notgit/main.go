package main

import (
	"fmt"
	"notgit/internal/commands"
	"os"
)

var command = map[string]func() error{
	"add":    commands.Add,
	"commit": commands.Commit,
	"config": commands.Config,
	"init":   commands.Init,
	"log":    commands.Log,
	"status": commands.Status,
}

func main() {
	if len(os.Args) < 2 {
		// print help
		return
	}

	// todo: add support for flags like --version

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
