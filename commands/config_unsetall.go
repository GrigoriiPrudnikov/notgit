package commands

import (
	"errors"
)

func unsetAllValues(args []string) error {
	if len(args) != 1 {
		return errors.New("invalid arguments")
	}

	unsetValue(args, false)
	unsetValue(args, true)

	return nil
}
