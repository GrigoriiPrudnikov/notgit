package utils

import (
	"bufio"
	"os"
	"strings"
)

func ParseINI(filePath string) (map[string]map[string]string, error) {
	result := make(map[string]map[string]string)
	var currentSection string

	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		if line == "" || strings.HasPrefix(line, "#") || strings.HasPrefix(line, ";") {
			continue
		}

		if strings.HasPrefix(line, "[") && strings.HasSuffix(line, "]") {
			currentSection = strings.TrimSpace(line[1 : len(line)-1])
			if _, exists := result[currentSection]; !exists {
				result[currentSection] = make(map[string]string)
			}
			continue
		}

		if equalIndex := strings.Index(line, "="); equalIndex > 0 {
			key := strings.TrimSpace(line[:equalIndex])
			value := strings.TrimSpace(line[equalIndex+1:])
			if currentSection == "" {
				currentSection = "root"
				if _, exists := result[currentSection]; !exists {
					result[currentSection] = make(map[string]string)
				}
			}
			result[currentSection][key] = value
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return result, nil
}
