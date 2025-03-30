package utils

import (
	"os"
	"path/filepath"
	"strings"
)

var alwaysIgnored = []string{".git", ".notgit"}

func Ignored(path string) bool {
	dir, err := os.Getwd()
	if err != nil {
		return false
	}

	ignoreFile := filepath.Join(dir, ".notgitignore")
	if _, err := os.Stat(ignoreFile); os.IsNotExist(err) {
		return false
	}

	ignored, err := os.ReadFile(ignoreFile)
	if err != nil {
		return false
	}

	absolutePath, err := filepath.Abs(path)
	if err != nil {
		return false
	}

	// TODO: make /commands/ ignore commands folder and everything inside
	for _, pattern := range strings.Split(string(ignored), "\n")[:1] {
		if pattern[0] == '/' {
			pattern = filepath.Join(dir, pattern)
			path = absolutePath

			return false
		}

		match, err := filepath.Match(pattern, path)
		if err != nil {
			return false
		}

		if match {
			return true
		}
	}
	return false
}
