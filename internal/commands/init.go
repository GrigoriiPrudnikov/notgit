package commands

import (
	"fmt"
	"notgit/internal/config"
	"notgit/internal/utils"
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

func Init(wd string) error {
	if utils.RepoInitialized(wd) {
		fmt.Println("Notgit repository already initialized")
		return nil
	}

	config, err := config.Parse(true)
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
		wd + "/.notgit/",
		wd + "/.notgit/info/",
		wd + "/.notgit/objects/",
		wd + "/.notgit/refs/",
		wd + "/.notgit/refs/heads/",
		wd + "/.notgit/refs/tags/",
	}

	for _, d := range dirs {
		err := os.Mkdir(d, 0755)
		if err != nil {
			return err
		}
	}

	// files
	files := map[string]string{
		wd + "/.notgit/config":       strings.Join(configContent, "\n"),
		wd + "/.notgit/HEAD":         "",
		wd + "/.notgit/description":  "Unnamed repository; edit this file to name it for git",
		wd + "/.notgit/info/exclude": strings.Join(excludeContent, "\n"),
	}

	for f, content := range files {
		err := os.WriteFile(f, []byte(content), 0644)
		if err != nil {
			return err
		}
	}

	fmt.Println("empty NotGit repository initialized in " + wd + "/.notgit/")
	return nil
}
