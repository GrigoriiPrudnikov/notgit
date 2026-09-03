package commands

import (
	"fmt"
	"notgit/internal/repo"
)

func Commit(r *repo.Repo, params []string, options map[string]string) error {
	var message string

	if _, ok := options["m"]; ok {
		message = options["m"]
	} else {
		message = options["message"]
	}

	if message == "" {
		return fmt.Errorf("commit message is required")
	}

	return r.Commit(message)
}
