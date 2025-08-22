package main

import (
	"fmt"
	"notgit/internal/commands"
	"notgit/internal/garbagecollector"
	"os"
)

var command = map[string]func(wd string) error{
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
		command["version"]("")
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

	wd, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
		return
	}

	err = execute(wd)
	if err != nil {
		fmt.Println(err)
	}

	err = garbagecollector.CollectGarbage(wd)
	if err != nil {
		fmt.Println(err)
	}
}
