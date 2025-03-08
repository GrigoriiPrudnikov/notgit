package main

import (
	"fmt"
	"notgit/commands"
	"os"
)

var actions = map[string]func([]string){
	"init":   commands.Init,
	"config": commands.Config,
}

func main() {
	if len(os.Args) < 2 {
		// print help
		return
	}

	// fmt.Println(os.Args[1:]) // print all args

	args := []string{}
	if len(os.Args) > 2 {
		args = os.Args[2:]
	}

	_, exists := actions[os.Args[1]]
	if !exists {
		fmt.Println("Command not found")
		return
	}

	actions[os.Args[1]](args)
}
