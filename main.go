package main

import (
	"fmt"
	"notgit/commands"
	"os"
)

var command = map[string]func() error{
	"config": commands.Config,
	"init":   commands.Init,
}

func main() {
	if len(os.Args) < 2 {
		// print help
		return
	}

	_, exists := command[os.Args[1]]
	if !exists {
		fmt.Println("Command not found")
		return
	}

	err := command[os.Args[1]]()
	if err != nil {
		fmt.Println(err)
	}
}
