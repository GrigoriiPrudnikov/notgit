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

	staged, modified, untracked := status.GetStatus()

	fmt.Println("staged:")
	fmt.Println(staged)
	fmt.Println("modified:")
	fmt.Println(modified)
	fmt.Println("untracked:")
	fmt.Println(untracked)

	return nil
}
