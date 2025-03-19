package commands

import (
	"errors"
	"notgit/utils"
)

func unsetValue(args []string, global bool) error {
	if len(args) != 1 {
		return errors.New("invalid arguments")
	}

	section, key, err := utils.GetSectionAndKey(args)
	if err != nil {
		return err
	}

	config, err := utils.ParseConfig(global)
	if err != nil {
		return err
	}

	delete(config[section], key)

	err = utils.UpdateConfig(config, global)
	if err != nil {
		return err
	}

	return nil
}
