package commands

import (
	"errors"
	"notgit/utils"
)

func setValue(args []string, global bool) error {
	if len(args) != 2 {
		return errors.New("invalid arguments")
	}

	section, key, err := utils.GetSectionAndKey(args)
	if err != nil {
		return err
	}
	value := args[1]

	config, err := utils.ParseConfig(global)
	if err != nil {
		return err
	}

	config[section][key] = value

	err = utils.UpdateConfig(config, global)
	if err != nil {
		return err
	}

	return nil
}
