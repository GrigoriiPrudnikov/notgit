package config

import (
	"errors"
	"strings"

	"gopkg.in/ini.v1"
)

func Set(c map[string]map[string]string, global bool) error {
	path, err := getConfigPath(global)
	if err != nil {
		return err
	}

	config := ini.Empty()

	for section, keys := range c {
		for key, value := range keys {
			config.Section(section).Key(key).SetValue(value)
		}
	}

	err = config.SaveToIndent(path, "  ")

	return err
}

// TODO: move it somewhere else
func GetSectionAndKey(args []string) (string, string, error) {
	if len(args) == 0 {
		return "", "", errors.New("invalid arguments")
	}

	parts := strings.Split(args[0], ".")
	if len(parts) != 2 {
		return "", "", errors.New("invalid arguments")
	}

	section, key := parts[0], parts[1]
	if section == "" || key == "" {
		return "", "", errors.New("invalid arguments")
	}

	return section, key, nil
}
