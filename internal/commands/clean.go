package commands

import (
	"notgit/internal/repo"
	"os"
)

func Clean(r *repo.Repo, params []string, opts map[string]string) error {
	if r == nil {
		return nil
	}
	os.RemoveAll(r.Wd + "/.notgit")

	return nil
}
