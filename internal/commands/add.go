package commands

import "notgit/internal/repo"

func Add(r *repo.Repo, params []string, options map[string]string) error {
	if options["u"] != "" {
		filesToStage := []string{}

		for key := range r.Index {
			filesToStage = append(filesToStage, key)
		}

		for _, file := range filesToStage {
			err := r.Stage(file)
			if err != nil {
				return err
			}
		}
	}

	for _, p := range params {
		err := r.Stage(p)
		if err != nil {
			return err
		}
	}

	return nil
}
