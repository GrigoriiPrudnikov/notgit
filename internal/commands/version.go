package commands

import (
	"fmt"
	"notgit/internal/version"
)

func Version(_ string) error {
	fmt.Println("notgit version", version.Version)

	return nil
}
