package commands

import (
	"flag"
	"notgit/internal/tree"
	"notgit/utils"
	"os"
)

func Status() error {
	var short bool

	fs := flag.NewFlagSet("status", flag.ExitOnError)
	fs.BoolVar(&short, "s", false, "short")
	fs.BoolVar(&short, "short", false, "short")
	fs.Parse(os.Args[2:])

	root := tree.Root()

	utils.PrintStruct(root)

	return nil
}
