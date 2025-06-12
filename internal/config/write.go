package config

import "gopkg.in/ini.v1"

func Write(c map[string]map[string]string, global bool) error {
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

	return config.SaveToIndent(path, "")
}
