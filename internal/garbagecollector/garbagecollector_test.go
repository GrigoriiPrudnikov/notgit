package garbagecollector

import (
	"notgit/internal/commands"
	"os"
	"path/filepath"
	"testing"
)

func TestUnusedObject(t *testing.T) {
	dir := t.TempDir()

	err := commands.Init(dir)
	if err != nil {
		t.Fatal(err)
	}

	path := filepath.Join(dir, "test_file")

	err = os.WriteFile(path, []byte("test content"), 0644)
	if err != nil {
		t.Fatal(err)
	}

	// TODO
}

func TestEmptyDir(t *testing.T) {
	// TODO
}
