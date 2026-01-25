package main

import (
	"fmt"
	"notgit/internal/commands"
	"notgit/internal/repo"
	"notgit/internal/utils"
	"os"
	"slices"
)

var commandsAvailableWithoutRepo = []string{"version", "init", "config"}

func main() {
	args := os.Args
	if len(args) == 1 {
		// print help
		return
	}

	wd, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
		return
	}

	command := os.Args[1]
	parsedArgs, err := utils.ParseArgs(os.Args[2:])
	if err != nil {
		fmt.Println("notgit: " + err.Error())
		os.Exit(1)
	}
	params := parsedArgs.Params
	opts := parsedArgs.Opts

	r, err := repo.ParseRepo(wd)

	if err == repo.ErrNotGitRepo {
		if !slices.Contains(commandsAvailableWithoutRepo, command) {
			fmt.Println("notgit: not a git repository")
			os.Exit(1)
		}

		action := commands.Commands[command].Run
		r := &repo.Repo{Wd: wd}
		err = action(r, params, opts)

		if err != nil {
			fmt.Println("notgit: " + err.Error())
			os.Exit(1)
		}
		return
	}

	if err != nil {
		fmt.Println("notgit: " + err.Error())
		os.Exit(1)
	}

	action, ok := commands.Commands[command]
	if !ok {
		fmt.Println("notgit: unknown command")
		os.Exit(1)
	}
	err = action.Run(r, params, opts)

	if err != nil {
		fmt.Println("notgit: " + err.Error())
		os.Exit(1)
	}

	// err = garbagecollector.CollectGarbage(wd)
	// if err != nil {
	// 	fmt.Println("notgit: failed to collect garbage\n" + err.Error())
	// }
}
