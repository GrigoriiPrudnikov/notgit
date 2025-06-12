package config

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
)

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

func getConfigPath(global bool) (string, error) {
	var dir string
	var err error

	if global {
		dir, err = os.UserHomeDir()
	} else {
		dir, err = os.Getwd()
	}

	if err != nil {
		return "", err
	}

	configPath := filepath.Join(dir, ".notgitconfig")
	if !global {
		configPath = filepath.Join(dir, ".notgit/config")

		if _, err := os.Stat(configPath); os.IsNotExist(err) {
			return "", errors.New("not a notgit repository")
		}
	}

	return configPath, nil
}
