package commands

import (
	"errors"
)

func getAllValues(args []string) error {
	if len(args) != 1 {
		return errors.New("invalid arguments")
	}

	getValue(args, false)
	getValue(args, true)

	return nil
}
