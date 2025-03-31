package utils

import "os"

func RepoInitialized(dir string) bool {
	notgitdir := dir + "/.notgit"

	if _, err := os.Stat(notgitdir); os.IsNotExist(err) {
		return false
	}
	return true
}
