package commands

import (
	"flag"
	"fmt"
	"notgit/internal/status"
	"os"
	"slices"
)

func Status() error {
	var short bool

	fs := flag.NewFlagSet("status", flag.ExitOnError)
	fs.BoolVar(&short, "s", false, "short")
	fs.BoolVar(&short, "short", false, "short")
	fs.Parse(os.Args[2:])

	modifiedStaged, untrackedStaged, modified, untracked := status.GetStatus()

	var indent string
	if len(untrackedStaged) > 0 || len(untracked) > 0 {
		indent = " "
	}

	for _, path := range untrackedStaged {
		fmt.Println(green("A")+indent, path)
	}
	for _, path := range modifiedStaged {
		if slices.Contains(modified, path) {
			fmt.Println(green("M")+red("M"), path)

		} else {
			fmt.Println(green("M")+indent, path)
		}
	}
	for _, path := range modified {
		fmt.Println(red("M")+indent, path)
	}
	for _, path := range untracked {
		fmt.Println(red("??"), path)
	}

	return nil
}

func red(s string) string {
	return "\033[31m" + s + "\033[0m"
}

func green(s string) string {
	return "\033[32m" + s + "\033[0m"
}
