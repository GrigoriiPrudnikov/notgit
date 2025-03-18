package utils

import (
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

func UpdateConfig(path string, c map[string]map[string]string) error {
	config := ini.Empty()

	for section, keys := range c {
		for key, value := range keys {
			config.Section(section).Key(key).SetValue(value)
		}
	}

	err := config.SaveToIndent(path, "  ")

	return err
}
