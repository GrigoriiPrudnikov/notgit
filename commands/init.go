package commands

import (
	"fmt"
	"notgit/utils"
	"os"
)

func Init(_args []string) {
	dir, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if !utils.RepoInitialized(dir) {
		os.Mkdir(dir+"/.notgit/", 0755)
		fmt.Println("Empty NotGit repository initialized in " + dir + "/.notgit/")
	} else {
		fmt.Println("Notgit repository already initialized")
		return
	}

	fmt.Println(dir)
}
