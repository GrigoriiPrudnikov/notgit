package utils

import (
	"errors"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/ini.v1"
)

func ParseConfig(global bool) (map[string]map[string]string, error) {
	path, err := getConfigPath(global)
	if err != nil {
		return nil, err
	}

	cfg, err := ini.Load(path)
	if err != nil {
		return nil, err
	}

	result := make(map[string]map[string]string)

	for _, section := range cfg.Sections() {
		sectionName := section.Name()
		if sectionName == ini.DefaultSection {
			sectionName = "root"
		}
		result[sectionName] = make(map[string]string)

		for _, key := range section.Keys() {
			result[sectionName][key.Name()] = key.String()
		}
	}

	return result, nil
}

func UpdateConfig(c map[string]map[string]string, global bool) error {
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
