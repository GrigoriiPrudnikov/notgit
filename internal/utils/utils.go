package utils

import (
	"bytes"
	"compress/zlib"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"slices"
	"strings"
)

func Compress(header string, b []byte) []byte {
	content := append([]byte(header), b...)

	var compressed bytes.Buffer
	w := zlib.NewWriter(&compressed)
	w.Write(content)
	w.Close()

	return compressed.Bytes()
}

func Decompress(b []byte) ([]byte, error) {
	r, err := zlib.NewReader(bytes.NewReader(b))
	if err != nil {
		return nil, err
	}
	defer r.Close()

	var out bytes.Buffer
	_, err = io.Copy(&out, r)
	if err != nil {
		return nil, err
	}

	return out.Bytes(), nil
}

func Hash(kind string, content []byte) string {
	header := fmt.Sprintf("%s %d\x00\n", kind, len(content))
	data := append([]byte(header), content...)
	sum := sha256.Sum256(data)

	return fmt.Sprintf("%x", sum)
}

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

func PrintStruct(v any) {
	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	enc.SetEscapeHTML(false)
	enc.SetIndent("", "  ")
	if err := enc.Encode(v); err != nil {
		fmt.Println("PrintStruct: failed to encode:", err)
		return
	}
	fmt.Print(buf.String())
}

func RepoInitialized(dir string) bool {
	notgitdir := dir + "/.notgit"

	if _, err := os.Stat(notgitdir); os.IsNotExist(err) {
		return false
	}
	return true
}
