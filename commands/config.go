package commands

import (
	"errors"
	"flag"
	"fmt"
	"notgit/utils"
	"os"
	"path/filepath"
	"strings"
)

func Config() error {
	// TODO: add Usage func
	var global, get, unset bool

	flagSet := flag.NewFlagSet("config", flag.ExitOnError)
	flagSet.BoolVar(&global, "global", false, "Use global config")
	flagSet.BoolVar(&get, "get", false, "Get value")
	flagSet.BoolVar(&unset, "unset", false, "Unset value")

	flagSet.Parse(os.Args[2:])
	args := flagSet.Args()

	fmt.Println(len(args))
	fmt.Println(global)
	fmt.Println(get)
	fmt.Println(unset)

	var err error

	if get {
		err = getValue(args, global)
	} else if unset {
		err = unsetValue(args, global)
	} else {
		err = setValue(args, global)
	}

	return err
}

func setValue(args []string, global bool) error {
	var dir string
	var err error

	if global {
		dir, err = os.UserHomeDir()
	} else {
		dir, err = os.Getwd()
	}

	if err != nil {
		return err
	}

	var configPath string
	if global {
		configPath = filepath.Join(dir, ".notgitconfig")
	} else {
		configPath = filepath.Join(dir, "/.notgit/config")

		if _, err := os.Stat(configPath); os.IsNotExist(err) {
			return errors.New("not a notgit repository")
		}
	}

	parts := strings.Split(args[0], ".")
	section, key := parts[0], parts[1]

	value := args[1]

	config, err := utils.ParseConfig(configPath)
	if err != nil {
		return err
	}

	config[section][key] = value

	err = utils.UpdateConfig(configPath, config)
	if err != nil {
		return err
	}

	return nil
}

func getValue(args []string, global bool) error { return nil }

func unsetValue(args []string, global bool) error { return nil }
