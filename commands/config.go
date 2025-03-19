package commands

import (
	"flag"
	"fmt"
	"os"
)

func Config() error {
	// TODO: add Usage func
	// TODO: add --local flag
	// TODO: add --add flag
	// TODO: add --get-all flag
	// TODO: add --unset-all flag
	var help, global, get, unset bool

	flagSet := flag.NewFlagSet("config", flag.ExitOnError)
	flagSet.BoolVar(&global, "global", false, "Use global config")
	flagSet.BoolVar(&get, "get", false, "Get value")
	flagSet.BoolVar(&unset, "unset", false, "Unset value")
	flagSet.BoolVar(&help, "help", false, "Show help")

	flagSet.Parse(os.Args[2:])
	args := flagSet.Args()

	if help {
		usage()
		return nil
	}
	if get {
		return getValue(args, global)
	}
	if unset {
		return unsetValue(args, global)
	}

	return setValue(args, global)
}

func usage() {
	fmt.Println(`usage: notgit config [<options>]

Config file location:
    --global    use global config file
    --local     use local config file
Action:
    --get       get value
    --unset     unset value`)
}
