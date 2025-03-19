package commands

import (
	"errors"
	"fmt"
	"notgit/utils"
)

func getValue(args []string, global bool) error {
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

	value, ok := config[section][key]
	if !ok {
		return errors.New("key not found")
	}

	fmt.Println(value)

	return nil
}
