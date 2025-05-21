package main

import (
	"fmt"
	"notgit/commands"
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

	_, exists := command[os.Args[1]]
	if !exists {
		// TODO: add help like this:
		// The most similar commands are
		//    diff
		//    fsck

		fmt.Printf("notgit: %s is not a git command. See 'notgit --help'.\n", os.Args[1])
		return
	}

	err := command[os.Args[1]]()
	if err != nil {
		fmt.Println("error:", err)
	}
}
