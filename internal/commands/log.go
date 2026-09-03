package commands

import "notgit/internal/repo"

func Log(r *repo.Repo, params []string, options map[string]string) error {
	return r.Log()
}
