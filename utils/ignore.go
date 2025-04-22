package utils

import (
	"os"
	"path/filepath"
	"slices"
	"strings"
)

var alwaysIgnored = []string{".git", ".notgit"}

// Checks if a file or directory path should be ignored based on rules.
func Ignored(path string) bool {
	base := filepath.Base(path)
	if slices.Contains(alwaysIgnored, base) {
		return true
	}

	wd, err := os.Getwd()
	if err != nil {
		return false
	}

	ignoreFile := filepath.Join(wd, ".notgitignore")
	data, err := os.ReadFile(ignoreFile)
	if err != nil {
		return false
	}

	absPath, err := filepath.Abs(path)
	if err != nil {
		return false
	}

	relPath, err := filepath.Rel(wd, absPath)
	if err != nil {
		return false
	}
	relPath = filepath.ToSlash(relPath)

	lines := strings.Split(string(data), "\n")
	for _, pattern := range lines {
		pattern = strings.TrimSpace(pattern)
		if pattern == "" || strings.HasPrefix(pattern, "#") {
			continue
		}

		if strings.HasSuffix(pattern, "/") {
			// Directory pattern: check if path starts with it
			prefix := strings.TrimSuffix(pattern, "/")
			if strings.HasPrefix(relPath, prefix+"/") || relPath == prefix {
				return true
			}
		} else if strings.HasPrefix(pattern, "/") {
			// Root-relative pattern
			match, _ := filepath.Match(pattern[1:], relPath)
			if match || relPath == pattern[1:] {
				return true
			}
		} else {
			// General pattern
			match, _ := filepath.Match(pattern, relPath)
			if match || filepath.Base(relPath) == pattern {
				return true
			}
		}
	}

	return false
}
