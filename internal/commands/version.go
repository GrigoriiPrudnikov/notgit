package commands

import (
	"fmt"
	"notgit/internal/repo"
	"notgit/internal/version"
)

func Version(_ *repo.Repo, params []string, options map[string]string) error {
	fmt.Println("notgit version", version.Version)

	return nil
}
