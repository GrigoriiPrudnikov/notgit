package commands

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"notgit/utils"
)

var allowedFlags = []string{"--get", "--global"}

func Config(args []string) error {
	// TODO: rewrite using flag package
	flags := []string{}

	for _, arg := range args {
		if arg[0] == '-' {
			flags = append(flags, arg)
		}
	}

	for _, flag := range flags {
		if !slices.Contains(allowedFlags, flag) {
			fmt.Println("invalid flag:", flag)
			return errors.New("invalid flag:" + flag)
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
		return err
	}

	var configPath string
	if scope == "global" {
		configPath = filepath.Join(dir, ".notgitconfig")
	} else {
		configPath = filepath.Join(dir, "/.notgit/config")
	}

	_, err = os.Stat(configPath)
	configFileExists := err == nil

	if !configFileExists && scope == "global" {
		file, err := os.Create(configPath)
		if err != nil {
			return err
		}
		defer file.Close()
	}

	if !configFileExists && scope == "local" {
		fmt.Println("")
		return errors.New("not in a notgit directory")
	}

	fileMap, err := utils.ParseConfig(filepath.Join(dir, ".notgitconfig"))
	if err != nil {
		return err
	}

	n := len(args)
	if mode == "set" {
		arg1, arg2 := args[n-2], args[n-1]

		if arg1[0] == '-' || arg2[0] == '-' {
			return errors.New("invalid arguments")
		}

		toSet := strings.Split(arg1, ".")
		if len(toSet) != 2 {
			return errors.New("invalid arguments")
		}

		section, key := toSet[0], toSet[1]

		err = utils.UpdateConfig(configPath, section, key, arg2)
		if err != nil {
			return err
		}

		return nil
	}

	arg := args[n-1]

	if arg[0] == '-' {
		return errors.New("invalid arguments")
	}

	toGet := strings.Split(arg, ".")
	if len(toGet) != 2 {
		return errors.New("invalid arguments")
	}

	section, key := toGet[0], toGet[1]

	value, ok := fileMap[section][key]
	if !ok {
		return errors.New("key not found")
	}

	fmt.Println(value)
	return nil
}
