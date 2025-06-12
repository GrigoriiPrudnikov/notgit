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

	modified, untracked, deleted := status.GetUnstaged()
	modifiedStaged, added, deletedStaged := status.GetStaged()

	totalChanges := len(added) + len(modified) + len(modifiedStaged) + len(untracked)
	if totalChanges == 0 {
		fmt.Println("nothing to commit, working tree clean")
		return nil
	}

	for _, path := range added {
		if slices.Contains(modified, path) {
			fmt.Println(green("A")+red("M"), path)
		} else {
			fmt.Println(green("A "), path)
		}
	}
	for _, path := range modifiedStaged {
		if slices.Contains(modified, path) {
			fmt.Println(green("M")+red("M"), path)
		} else {
			fmt.Println(green("M "), path)
		}
	}
	for _, path := range deletedStaged {
		fmt.Println(green("D "), path)
	}
	for _, path := range modified {
		if slices.Contains(modifiedStaged, path) {
			continue
		}
		fmt.Println(red("M "), path)
	}
	for _, path := range deleted {
		fmt.Println(red("D "), path)
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
