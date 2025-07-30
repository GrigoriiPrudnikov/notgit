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
	"sort"
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

// Returns patterns from .notgitignore file
func parseIgnoreFile() ([]string, error) {
	wd, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	data, err := os.ReadFile(filepath.Join(wd, ".notgitignore"))
	if os.IsNotExist(err) {
		return nil, nil
	}

	lines := strings.Split(string(data), "\n")
	var patterns []string
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" || trimmed[0] == '#' {
			continue
		}
		patterns = append(patterns, line)
	}

	return patterns, nil
}

var alwaysIgnored = []string{".git", ".notgit"}

// Returns true if given path is ignored (file is in .notgitignore)
func Ignored(path string) bool {
	if slices.Contains(alwaysIgnored, path) {
		return true
	}

	for _, pattern := range alwaysIgnored {
		if strings.HasPrefix(path, pattern) {
			return true
		}
	}

	ignored := false

	patterns, err := parseIgnoreFile()
	if err != nil {
		return false
	}

	for _, pattern := range patterns {
		negated := strings.HasPrefix(pattern, "!")
		if negated {
			pattern = pattern[1:]
		}
		fullMatch, _ := filepath.Match(pattern, path)
		baseMatch, _ := filepath.Match(pattern, filepath.Base(path))
		dirMatch := strings.HasSuffix(pattern, "/") && strings.HasPrefix(path, pattern)

		if fullMatch || baseMatch || dirMatch {
			ignored = !negated
		}
	}

	return ignored
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

// FOR DEBUGGING ONLY
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
		if dir == "/" {
			return false
		}

		return RepoInitialized(filepath.Dir(dir))
	}
	return true
}

func GetSortedKeys[T any](m map[string]T) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
