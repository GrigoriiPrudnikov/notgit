package commands

import (
	"flag"
	"fmt"
	"os"
)

func Status() error {
	var short bool

	fs := flag.NewFlagSet("status", flag.ExitOnError)
	fs.BoolVar(&short, "s", false, "short")
	fs.BoolVar(&short, "short", false, "short")
	fs.Parse(os.Args[2:])

	fmt.Println("not implemented")

	// changes := status.GetChanges()
	//
	// if len(changes) == 0 {
	// fmt.Println("nothing to commit, working tree clean")
	// 	return nil
	// }
	//
	// for path, change := range changes {
	// 	statusString := [2]string{" ", " "}
	//
	// 	if change.Unstaged == status.Added {
	// 		statusString[0] = red("?")
	// 		statusString[1] = red("?")
	// 	}
	//
	// 	if change.Unstaged == status.Modified {
	// 		statusString[1] = red("M")
	// 	}
	// 	if change.Staged == status.Modified {
	// 		statusString[0] = statusString[1]
	// 		statusString[1] = green("M")
	// 	}
	// 	if change.Unstaged == status.Deleted {
	// 		statusString[1] = red("D")
	// 	}
	//
	// 	if change.Staged == status.Added {
	// 		statusString[1] = green("A")
	// 	}
	//
	// 	fmt.Println(statusString[0]+statusString[1], path)
	// }

	return nil
}

// func red(s string) string {
// 	return "\033[31m" + s + "\033[0m"
// }
//
// func green(s string) string {
// 	return "\033[32m" + s + "\033[0m"
// }
