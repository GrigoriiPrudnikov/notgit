package utils

import (
	"os"
	"path/filepath"
	"strings"
)

func InWorkingDirectory(path string) bool {
	wd, err := os.Getwd()
	if err != nil {
		return false
	}

	absPath, err := filepath.Abs(path)
	if err != nil {
		return false
	}

	if strings.HasPrefix(absPath, wd) {
		return true
	}

	return false
}
