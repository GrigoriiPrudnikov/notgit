package commands

import (
	"errors"
	"flag"
	"os"
)

type action struct {
	flag *bool
	fn   func() error
}

func Config() error {
	var global, local, get, getAll, unset, unsetAll, add, list, help bool

	flagSet := flag.NewFlagSet("config", flag.ExitOnError)
	flagSet.BoolVar(&global, "global", false, "Use global config")
	flagSet.BoolVar(&local, "local", false, "Use local config")
	flagSet.BoolVar(&get, "get", false, "Get value")
	flagSet.BoolVar(&getAll, "get-all", false, "Get all values")
	flagSet.BoolVar(&unset, "unset", false, "Unset value")
	flagSet.BoolVar(&unsetAll, "unset-all", false, "Unset all values")
	flagSet.BoolVar(&add, "add", false, "Add value")
	flagSet.BoolVar(&list, "list", false, "List values")
	flagSet.BoolVar(&help, "help", false, "Show help")

	flagSet.Parse(os.Args[2:])
	args := flagSet.Args()

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
