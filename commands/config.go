package commands

import (
	"fmt"
	"notgit/utils"
	"os"
	"path/filepath"
	"slices"
)

func Config(args []string) {
	// if len(args) < 2 {
	// 	fmt.Println("usage: notgit config [<options>]")
	// 	return
	// }
	homeDir, err := os.UserHomeDir()

	if err != nil {
		fmt.Println("Error getting home directory:", err)
		return
	}

	if _, err := os.Stat(filepath.Join(homeDir, ".notgitconfig")); os.IsNotExist(err) {
		filePath := filepath.Join(homeDir, ".notgitconfig")

		file, err := os.Create(filePath)
		if err != nil {
			fmt.Println("Error creating file:", err)
			return
		}
		defer file.Close()
		return
	}

	file, _ := utils.ParseINI("~/.notgitconfig")

	fmt.Println(file)

	for _, arg := range args {
		if len(arg) < 2 {
			// add error message
			return
		}
	}

	global := slices.Contains(args, "--global")
	fmt.Println(args)

	fmt.Println(global)
}
