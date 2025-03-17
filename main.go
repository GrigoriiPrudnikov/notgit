package main

import (
	"fmt"
	"notgit/commands"
	"os"
)

var actions = map[string]func([]string) error{
	"config": commands.Config,
	"init":   commands.Init,
}

func main() {
	if len(os.Args) < 2 {
		// print help
		return
	}

	args := []string{}
	if len(os.Args) > 2 {
		args = os.Args[2:]
	}

	_, exists := actions[os.Args[1]]
	if !exists {
		fmt.Println("Command not found")
		return
	}

	err := actions[os.Args[1]](args)
	if err != nil {
		fmt.Println("error:", err)
	}
}
