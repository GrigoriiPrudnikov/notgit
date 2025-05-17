package commands

import (
	"flag"
	"fmt"
	"notgit/internal/status"
	"notgit/internal/tree"
	"os"
)

func Status() error {
	var short bool

	fs := flag.NewFlagSet("status", flag.ExitOnError)
	fs.BoolVar(&short, "s", false, "short")
	fs.BoolVar(&short, "short", false, "short")
	fs.Parse(os.Args[2:])

	root := tree.Root()
	stagedTree := tree.Staged()

	modified, untracked := status.GetModifiedAndUntrackedFiles(root, stagedTree)

	fmt.Println("modified:")
	fmt.Println(modified)
	fmt.Println("untracked:")
	fmt.Println(untracked)
	fmt.Println("staged:")

	return nil
}
