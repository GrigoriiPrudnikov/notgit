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

func Init() error {
	dir, err := os.Getwd()
	if err != nil {
		return err
	}

	if utils.RepoInitialized(dir) {
		fmt.Println("Notgit repository already initialized")
		return nil
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	config, err := utils.ParseConfig(homeDir + "/.notgitconfig")
	if err != nil {
		return err
	}

	configContent := []string{}

	for section := range config {
		configContent = append(configContent, "["+section+"]")

		for key, value := range config[section] {
			configContent = append(configContent, key+" = "+value)
		}
	}

	// dirs
	dirs := []string{
		dir + "/.notgit/",
		dir + "/.notgit/info/",
		dir + "/.notgit/objects/",
		dir + "/.notgit/refs/",
		dir + "/.notgit/refs/heads/",
		dir + "/.notgit/refs/tags/",
	}

	for _, d := range dirs {
		err := os.Mkdir(d, 0755)
		if err != nil {
			return err
		}
	}

	// files
	files := map[string]string{
		dir + "/.notgit/config":       strings.Join(configContent, "\n"),
		dir + "/.notgit/HEAD":         "ref: refs/heads/master",
		dir + "/.notgit/description":  "Unnamed repository; edit this file to name it for git",
		dir + "/.notgit/info/exclude": strings.Join(excludeContent, "\n"),
	}

	for f, content := range files {
		err := os.WriteFile(f, []byte(content), 0644)
		if err != nil {
			return err
		}
	}

	fmt.Println("empty NotGit repository initialized in " + dir + "/.notgit/")
	return nil
}
