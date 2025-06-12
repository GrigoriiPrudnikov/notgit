package config

import (
	"gopkg.in/ini.v1"
)

func Parse(global bool) (map[string]map[string]string, error) {
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
