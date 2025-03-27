package utils

import (
	"bufio"
	"os"
	"path/filepath"
)

func Ignored(path string) bool {
	dir, err := os.Getwd()
	if err != nil {
		return true
	}

	ignoreFile := filepath.Join(dir, ".notgitignore")
	if _, err := os.Stat(ignoreFile); err != nil {
		return false
	}

	ignored := readIgnoreFile(ignoreFile)
	for _, ignoredPath := range ignored {
		// If the ignored path starts with '/', treat it as an absolute path
		if ignoredPath[0] == byte('/') {
			ignoredPath = filepath.Join(dir, ignoredPath)
		}

		path = filepath.Join(dir, path)
		match, err := filepath.Match(ignoredPath, path)

		if err != nil {
			return false
		}

		return match
	}

	return false
}

func readIgnoreFile(path string) []string {
	file, err := os.Open(path)
	if err != nil {
		return nil
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" || line[0] == '#' {
			continue
		}

		lines = append(lines, line)
	}
	return lines
}
