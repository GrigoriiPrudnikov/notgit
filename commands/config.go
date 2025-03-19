package commands

import (
	"errors"
	"flag"
	"fmt"
	"notgit/utils"
	"os"
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

// TODO: move these functions to its own files
func setValue(args []string, global bool) error {
	if len(args) != 2 {
		return errors.New("invalid arguments")
	}

	parts := strings.Split(args[0], ".")
	section, key := parts[0], parts[1]
	value := args[1]

	config, err := utils.ParseConfig(global)
	if err != nil {
		return err
	}

	config[section][key] = value

	err = utils.UpdateConfig(config, global)
	if err != nil {
		return err
	}

	return nil
}

func getValue(args []string, global bool) error {
	if len(args) != 1 {
		return errors.New("invalid arguments")
	}

	parts := strings.Split(args[0], ".")
	section, key := parts[0], parts[1]

	config, err := utils.ParseConfig(global)
	if err != nil {
		return err
	}

	value, ok := config[section][key]
	if !ok {
		return errors.New("key not found")
	}

	fmt.Println(value)

	return nil
}

func unsetValue(args []string, global bool) error {
	if len(args) != 1 {
		return errors.New("invalid arguments")
	}

	parts := strings.Split(args[0], ".")
	section, key := parts[0], parts[1]

	config, err := utils.ParseConfig(global)
	if err != nil {
		return err
	}

	delete(config[section], key)

	err = utils.UpdateConfig(config, global)
	if err != nil {
		return err
	}

	return nil
}
