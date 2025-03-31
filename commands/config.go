package commands

import (
	"errors"
	"flag"
	"fmt"
	"notgit/internal/config"
	"os"
)

type action struct {
	flag *bool
	fn   func() error
}

func Config() error {
	// TODO: somehow rewrite
	var global, local, get, getAll, unset, unsetAll, add, list, help bool

	fs := flag.NewFlagSet("config", flag.ExitOnError)
	fs.BoolVar(&global, "global", false, "Use global config")
	fs.BoolVar(&local, "local", false, "Use local config")
	fs.BoolVar(&get, "get", false, "Get value")
	fs.BoolVar(&getAll, "get-all", false, "Get all values")
	fs.BoolVar(&unset, "unset", false, "Unset value")
	fs.BoolVar(&unsetAll, "unset-all", false, "Unset all values")
	fs.BoolVar(&add, "add", false, "Add value")
	fs.BoolVar(&list, "list", false, "List values")
	fs.BoolVar(&help, "help", false, "Show help")

	fs.Parse(os.Args[2:])
	args := fs.Args()

	if global && local {
		return errors.New("--global and --local flags can't be used together")
	}

	if local {
		global = false
	}

	actions := []action{
		{&get, func() error { return getValue(args, global) }},
		{&getAll, func() error { return getAllValues(args) }},
		{&unset, func() error { return unsetValue(args, global) }},
		{&unsetAll, func() error { return unsetAllValues(args) }},
		{&add, func() error { return setValue(args, global) }},
		{&list, func() error { return ListValues(global) }},
	}

	for _, a := range actions {
		if *a.flag {
			return a.fn()
		}
	}

	return setValue(args, global)
}

func getValue(args []string, global bool) error {
	if len(args) != 1 {
		return errors.New("invalid arguments")
	}

	section, key, err := config.GetSectionAndKey(args)
	if err != nil {
		return err
	}

	config, err := config.Parse(global)
	if err != nil {
		return err
	}

	value, ok := config[section][key]
	if !ok {
		return nil
	}

	fmt.Println(value)

	return nil
}

func getAllValues(args []string) error {
	if len(args) != 1 {
		return errors.New("invalid arguments")
	}

	getValue(args, false)
	getValue(args, true)

	return nil
}

func ListValues(global bool) error {
	config, err := config.Parse(global)
	if err != nil {
		return err
	}

	for section, values := range config {
		for key, value := range values {
			fmt.Printf("%s.%s = %s\n", section, key, value)
		}
	}

	return nil
}

func setValue(args []string, global bool) error {
	if len(args) != 2 {
		return errors.New("invalid arguments")
	}

	section, key, err := config.GetSectionAndKey(args)
	if err != nil {
		return err
	}
	value := args[1]

	cfg, err := config.Parse(global)
	if err != nil {
		return err
	}

	cfg[section][key] = value

	err = config.Set(cfg, global)
	if err != nil {
		return err
	}

	return nil
}

func unsetValue(args []string, global bool) error {
	if len(args) != 1 {
		return errors.New("invalid arguments")
	}

	section, key, err := config.GetSectionAndKey(args)
	if err != nil {
		return err
	}

	cfg, err := config.Parse(global)
	if err != nil {
		return err
	}

	delete(cfg[section], key)

	err = config.Set(cfg, global)
	if err != nil {
		return err
	}

	return nil
}

func unsetAllValues(args []string) error {
	if len(args) != 1 {
		return errors.New("invalid arguments")
	}

	unsetValue(args, false)
	unsetValue(args, true)

	return nil
}
