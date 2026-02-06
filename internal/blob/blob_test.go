package blob

import (
	"os"
	"path/filepath"
	"testing"

	"notgit/internal/object"
)

func TestCreateBlob(t *testing.T) {
	tempDir := t.TempDir()
	path := filepath.Join(tempDir, "test.txt")

	err := os.WriteFile(path, []byte("test"), 0644)
	if err != nil {
		t.Error(err)
	}

	b, err := NewBlob(path)
	if err != nil {
		t.Error(err)
	}

	if b.BasePath() != "test.txt" {
		t.Error("unexpected base path")
	}
}

func TestBlobExists(t *testing.T) {
	tempDir := t.TempDir()
	path := filepath.Join(tempDir, "test.txt")

	err := os.WriteFile(path, []byte("test"), 0644)
	if err != nil {
		t.Error(err)
	}

	b, err := NewBlob(path)
	if err != nil {
		t.Error(err)
	}

	store := object.NewFSStore(tempDir)

	if err := b.Write(store); err != nil {
		t.Error(err)
	}

	// Assert by checking the filesystem, since Blob no longer knows about wd.
	if _, err := os.Stat(filepath.Join(tempDir, ".notgit", "objects", b.Hash()[:2], b.Hash()[2:])); os.IsNotExist(err) {
		t.Error("blob should exist")
	}
}

func TestBlobWrite(t *testing.T) {
	tempDir := t.TempDir()
	path := filepath.Join(tempDir, "test.txt")

	err := os.WriteFile(path, []byte("test"), 0644)
	if err != nil {
		t.Error(err)
	}

	b, err := NewBlob(path)
	if err != nil {
		t.Error(err)
	}

	store := object.NewFSStore(tempDir)

	if err := b.Write(store); err != nil {
		t.Error(err)
	}

	if _, err := os.Stat(filepath.Join(tempDir, ".notgit", "objects", b.Hash()[:2], b.Hash()[2:])); os.IsNotExist(err) {
		t.Error("blob should exist")
	}
}
