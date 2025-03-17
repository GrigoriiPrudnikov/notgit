package utils

import (
	"os"

	"gopkg.in/ini.v1"
)

func ParseConfig(path string) (map[string]map[string]string, error) {
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

func UpdateConfig(path string, section string, key string, value string) error {
	configExists := false
	if _, err := os.Stat(path); os.IsExist(err) {
		configExists = true
	}

	var config *ini.File
	var err error

	if configExists {
		config, err = ini.Load(path)
		if err != nil {
			return err
		}
	} else {
		config = ini.Empty()
	}

	config.Section(section).Key(key).SetValue(value)
	err = config.SaveTo(path)
	if err != nil {
		return err
	}

	return nil
}
