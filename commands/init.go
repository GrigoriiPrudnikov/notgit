package commands

import (
	"fmt"
	"notgit/utils"
	"os"
	"strings"
)

var excludeContent = []string{
	"# git ls-files --others --exclude-from=.git/info/exclude",
	"# Lines that start with '#' are comments.",
	"# For a project mostly in C, the following would be a good set of",
	"# exclude patterns (uncomment them if you want to use them):",
	"# *.[oa]",
	"# *~",
}

func Init(_args []string) {
	dir, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if !utils.RepoInitialized(dir) {
		err := os.Mkdir(dir+"/.notgit/", 0755)
		if err != nil {
			fmt.Println(err)
			return
		}

		_, err = os.Create(dir + "/.notgit/config")
		if err != nil {
			fmt.Println(err)
			return
		}

		homeDir, err := os.UserHomeDir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		config, err := utils.ParseConfig(homeDir + "/.notgitconfig")
		if err != nil {
			fmt.Println(err)
			return
		}

		configContent := []string{}

		for section := range config {
			configContent = append(configContent, "["+section+"]")

			for key, value := range config[section] {
				configContent = append(configContent, key+" = "+value)
			}
		}

		// dirs
		err = os.Mkdir(dir+"/.notgit/info/", 0755)
		if err != nil {
			fmt.Println(err)
			return
		}

		err = os.Mkdir(dir+"/.notgit/objects/", 0755)
		if err != nil {
			fmt.Println(err)
			return
		}

		err = os.Mkdir(dir+".notgit/refs", 0755)
		if err != nil {
			fmt.Println(err)
			return
		}

		err = os.Mkdir(dir+".notgit/refs/heads", 0755)
		if err != nil {
			fmt.Println(err)
			return
		}

		err = os.Mkdir(dir+".notgit/refs/tags", 0755)
		if err != nil {
			fmt.Println(err)
			return
		}

		// files
		err = os.WriteFile(dir+"/.notgit/config", []byte(strings.Join(configContent, "\n")), 0644)
		if err != nil {
			fmt.Println(err)
			return
		}

		err = os.WriteFile(dir+"/.notgit/HEAD", []byte("ref: refs/heads/master"), 0644)
		if err != nil {
			fmt.Println(err)
			return
		}

		err = os.WriteFile(dir+"/.notgit/description", []byte("Unnamed repository; edit this file to name it for git"), 0644)
		if err != nil {
			fmt.Println(err)
			return
		}

		err = os.WriteFile(dir+"/.notgit/info/exclude", []byte(strings.Join(excludeContent, "\n")), 0644)
		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Println("empty NotGit repository initialized in " + dir + "/.notgit/")
		return
	}
	fmt.Println("Notgit repository already initialized")
}
