package commands

import (
	"fmt"
	"notgit/internal/version"
)

func Version() error {
	fmt.Println("notgit version", version.Version)

	return nil
}
