package commands

import (
	"fmt"
	"notgit/utils"
)

func ListValues(global bool) error {
	config, err := utils.ParseConfig(global)
	if err != nil {
		return err
	}

	for section, values := range config {
		for key, value := range values {
			fmt.Printf("%s.%s = %s\n", section, key, value)
		}
	}

	return nil
}
