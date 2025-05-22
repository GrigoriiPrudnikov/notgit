package commands

import (
	"flag"
	"fmt"
	"notgit/internal/status"
	"os"
)

func Status() error {
	var short bool

	fs := flag.NewFlagSet("status", flag.ExitOnError)
	fs.BoolVar(&short, "s", false, "short")
	fs.BoolVar(&short, "short", false, "short")
	fs.Parse(os.Args[2:])

	stagedModified, stagedUntracked, modified, untracked := status.GetStatus()

	var indent string
	if len(stagedUntracked) > 0 || len(untracked) > 0 {
		indent = " "
	}

	for _, path := range stagedUntracked {
		fmt.Println(green("A")+indent, path)
	}
	for _, path := range stagedModified {
		fmt.Println(green("M")+indent, path)
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
