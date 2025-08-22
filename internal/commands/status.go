package commands

import (
	"errors"
	"flag"
	"fmt"
	"notgit/internal/status"
	"notgit/internal/utils"
	"os"
)

func Status(wd string) error {
	if !utils.RepoInitialized(wd) {
		return errors.New("not a notgit repository")
	}

	var short bool

	fs := flag.NewFlagSet("status", flag.ExitOnError)
	fs.BoolVar(&short, "s", false, "short")
	fs.BoolVar(&short, "short", false, "short")
	fs.Parse(os.Args[2:])

	worktreeAndIndexDiff, indexAndHeadDiff := status.GetRepoStatus()

	if len(worktreeAndIndexDiff)+len(indexAndHeadDiff) == 0 {
		fmt.Println("nothing to commit, working tree clean")
		return nil
	}

	// print files that have staged and unstaged changes
	for path, unstagedStatus := range worktreeAndIndexDiff {
		if stagedStatus, ok := indexAndHeadDiff[path]; ok {
			fmt.Printf("%s%s %s\n", getUnstagedSign(unstagedStatus), getStagedSign(stagedStatus), path)
			delete(worktreeAndIndexDiff, path)
			delete(indexAndHeadDiff, path)
		}
	}

	for path, status := range worktreeAndIndexDiff {
		fmt.Printf(" %s %s\n", getUnstagedSign(status), path)
	}
	for path, status := range indexAndHeadDiff {
		fmt.Printf("%s  %s\n", getStagedSign(status), path)
	}

	return nil
}

func red(s string) string {
	return "\033[31m" + s + "\033[0m"
}

func green(s string) string {
	return "\033[32m" + s + "\033[0m"
}

func getStagedSign(s status.Status) string {
	switch s {
	case status.Added:
		return green("A")
	case status.Modified:
		return green("M")
	case status.Deleted:
		return green("D")
	}

	return " "
}

func getUnstagedSign(s status.Status) string {
	switch s {
	case status.Added:
		return red("?")
	case status.Modified:
		return red("M")
	case status.Deleted:
		return red("D")
	}

	return " "
}
