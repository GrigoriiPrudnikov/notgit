package commands

import (
	"fmt"
	"notgit/internal/repo"
)

func Init(r *repo.Repo, params []string, options map[string]string) error {
	err := r.Init()
	if err != nil {
		return err
	}

	fmt.Println("Initialized empty repository in", r.Wd)

	return nil
}
