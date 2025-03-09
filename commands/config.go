package commands

import (
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"notgit/utils"
)

var allowedFlags = []string{"--get", "--global"}

func Config(args []string) {
	flags := []string{}

	for _, arg := range args {
		if arg[0] == '-' {
			flags = append(flags, arg)
		}
	}

	for _, flag := range flags {
		if !slices.Contains(allowedFlags, flag) {
			fmt.Println("invalid flag:", flag)
			return
		}
	}

	var dir string
	var err error
	mode := "set"
	scope := "local"

	if slices.Contains(flags, "--get") {
		mode = "get"
	}

	if slices.Contains(flags, "--global") {
		scope = "global"
		dir, err = os.UserHomeDir()
	} else {
		dir, err = os.Getwd()
	}

	if err != nil {
		fmt.Println("error:", err)
		return
	}

	configPath := filepath.Join(dir, ".notgitconfig")
	_, err = os.Stat(configPath)
	configFileExists := err == nil

	if !configFileExists && scope == "global" {
		file, err := os.Create(configPath)
		if err != nil {
			fmt.Println(err)
			return
		}
		defer file.Close()
	}

	if !configFileExists && scope == "local" {
		fmt.Println("notgit repository not initialized")
		return
	}
	fmt.Println(dir)

	fileMap, err := utils.ParseConfig(filepath.Join(dir, ".notgitconfig"))
	if err != nil {
		fmt.Println(err)
		return
	}

	n := len(args)
	if mode == "set" {
		arg1, arg2 := args[n-2], args[n-1]

		if arg1[0] == '-' || arg2[0] == '-' {
			fmt.Println("invalid arguments")
			return
		}

		toSet := strings.Split(arg1, ".")
		if len(toSet) != 2 {
			fmt.Println("invalid arguments")
			return
		}

		section, key := toSet[0], toSet[1]

		if fileMap[section] == nil {
			fileMap[section] = make(map[string]string)
		}

		fileMap[section][key] = arg2
		content := []string{}

		for section := range fileMap {
			content = append(content, "["+section+"]")

			for key, value := range fileMap[section] {
				content = append(content, key+" = "+value)
			}
		}

		os.WriteFile(configPath, []byte(strings.Join(content, "\n")), 0644)
		return
	}

	arg := args[n-1]

	if arg[0] == '-' {
		fmt.Println("invalid arguments")
		return
	}

	toGet := strings.Split(arg, ".")
	if len(toGet) != 2 {
		fmt.Println("invalid arguments")
		return
	}

	section, key := toGet[0], toGet[1]

	value, ok := fileMap[section][key]
	if !ok {
		return
	}

	fmt.Println(value)
}
