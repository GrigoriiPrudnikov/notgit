package commands

import (
	"flag"
	"os"
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

	if get {
		err := getValue(args, global)
		return err
	}

	if unset {
		err := unsetValue(args, global)
		return err
	}

	err := setValue(args, global)
	return err
}
