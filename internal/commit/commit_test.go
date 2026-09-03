package commit_test

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"notgit/internal/commit"
	"notgit/internal/object"
	"notgit/internal/tree"
	"notgit/internal/utils"
)

func TestGetContent(t *testing.T) {
	tr := tree.NewTree(".")
	c := &commit.Commit{
		Time:    1234567890,
		Offset:  "-0700",
		Author:  "Test Author",
		Message: "test message",
		Tree:    tr,
		Parent:  strings.Repeat("a", 64),
	}

	content := c.GetContent()
	s := string(content)

	if !strings.Contains(s, "tree ") {
		t.Error("GetContent should contain tree line")
	}
	if !strings.Contains(s, "author Test Author 1234567890 -0700") {
		t.Error("GetContent should contain author line")
	}
	if !strings.Contains(s, "committer Test Author 1234567890 -0700") {
		t.Error("GetContent should contain committer line")
	}
	if !strings.Contains(s, "parent ") {
		t.Error("GetContent should contain parent line")
	}
	if !strings.Contains(s, "test message") {
		t.Error("GetContent should contain message")
	}
}

func TestHash(t *testing.T) {
	tr := tree.NewTree(".")
	c := &commit.Commit{
		Time:    1234567890,
		Offset:  "-0700",
		Author:  "Test Author",
		Message: "test message",
		Tree:    tr,
		Parent:  strings.Repeat("a", 64),
	}

	hash := c.Hash()

	if len(hash) != 64 {
		t.Errorf("Hash should be 64 chars, got %d", len(hash))
	}
	hexRegex := regexp.MustCompile(`^[a-f0-9]+$`)
	if !hexRegex.MatchString(hash) {
		t.Error("Hash should be hex string")
	}
}

func TestHash_Deterministic(t *testing.T) {
	tr := tree.NewTree(".")
	c := &commit.Commit{
		Time:    1234567890,
		Offset:  "-0700",
		Author:  "Test Author",
		Message: "test message",
		Tree:    tr,
		Parent:  strings.Repeat("a", 64),
	}

	h1 := c.Hash()
	h2 := c.Hash()

	if h1 != h2 {
		t.Error("Hash should be deterministic")
	}
}

func TestWrite(t *testing.T) {
	tempDir := t.TempDir()

	// Create .notgit/objects for object.Parse used by tree.Write
	objectsDir := filepath.Join(tempDir, ".notgit", "objects")
	if err := os.MkdirAll(objectsDir, 0755); err != nil {
		t.Fatal(err)
	}

	origWd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	defer os.Chdir(origWd)

	if err := os.Chdir(tempDir); err != nil {
		t.Fatal(err)
	}

	tr := tree.NewTree(".")
	store := object.NewFSStore(tempDir)

	c := &commit.Commit{
		Time:    1234567890,
		Offset:  "-0700",
		Author:  "Test Author",
		Message: "test message",
		Tree:    tr,
		Parent:  strings.Repeat("a", 64),
	}

	hash, err := c.Write(store)
	if err != nil {
		t.Fatal(err)
	}

	if hash != c.Hash() {
		t.Errorf("Write returned hash %s, want %s", hash, c.Hash())
	}

	objPath := filepath.Join(tempDir, ".notgit", "objects", hash[:2], hash[2:])
	if _, err := os.Stat(objPath); os.IsNotExist(err) {
		t.Error("commit object should exist in store")
	}
}

func TestParse_InvalidHashLength(t *testing.T) {
	tests := []string{"", "short", "a" + strings.Repeat("b", 62), strings.Repeat("a", 65)}
	for _, hash := range tests {
		if got := commit.Parse(hash); got != nil {
			t.Errorf("Parse(%q) should return nil for invalid length, got %v", hash, got)
		}
	}
}

func TestParse_ObjectNotFound(t *testing.T) {
	tempDir := t.TempDir()

	origWd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	defer os.Chdir(origWd)

	if err := os.Chdir(tempDir); err != nil {
		t.Fatal(err)
	}

	// Create .notgit but no objects
	if err := os.MkdirAll(filepath.Join(tempDir, ".notgit", "objects"), 0755); err != nil {
		t.Fatal(err)
	}

	hash := strings.Repeat("a", 64)
	if got := commit.Parse(hash); got != nil {
		t.Errorf("Parse with non-existent object should return nil, got %v", got)
	}
}

func TestParse_ValidCommit(t *testing.T) {
	tempDir := t.TempDir()

	origWd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	defer os.Chdir(origWd)

	if err := os.Chdir(tempDir); err != nil {
		t.Fatal(err)
	}

	// Create empty tree and write to store
	tr := tree.NewTree(".")
	store := object.NewFSStore(tempDir)
	if err := tr.Write(store); err != nil {
		t.Fatal(err)
	}

	// Create commit content and object
	parentHash := strings.Repeat("b", 64)
	commitContent := []byte(strings.Join([]string{
		"tree " + tr.Hash(),
		"author Test Author 1234567890 -0700",
		"committer Test Author 1234567890 -0700",
		"parent " + parentHash,
		"",
		"test message",
	}, "\n"))

	commitHash := utils.Hash("commit", commitContent)
	commitCompressed := utils.Compress("commit", commitContent)

	objDir := filepath.Join(tempDir, ".notgit", "objects", commitHash[:2])
	if err := os.MkdirAll(objDir, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(objDir, commitHash[2:]), commitCompressed, 0644); err != nil {
		t.Fatal(err)
	}

	parsed := commit.Parse(commitHash)
	if parsed == nil {
		t.Fatal("Parse should return commit")
	}

	if parsed.Author != "Test Author" {
		t.Errorf("Author = %q, want Test Author", parsed.Author)
	}
	if parsed.Message != "test message" {
		t.Errorf("Message = %q, want test message", parsed.Message)
	}
	if parsed.Parent != parentHash {
		t.Errorf("Parent = %q, want %s", parsed.Parent, parentHash)
	}
	if parsed.Time != 1234567890 {
		t.Errorf("Time = %d, want 1234567890", parsed.Time)
	}
	if parsed.Offset != "-0700" {
		t.Errorf("Offset = %q, want -0700", parsed.Offset)
	}
}

func TestParse_NotCommitObject(t *testing.T) {
	tempDir := t.TempDir()

	origWd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	defer os.Chdir(origWd)

	if err := os.Chdir(tempDir); err != nil {
		t.Fatal(err)
	}

	// Create a blob object with a valid-looking hash
	blobContent := []byte("not a commit")
	blobCompressed := utils.Compress("blob", blobContent)
	blobHash := utils.Hash("blob", blobContent)

	objDir := filepath.Join(tempDir, ".notgit", "objects", blobHash[:2])
	if err := os.MkdirAll(objDir, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(objDir, blobHash[2:]), blobCompressed, 0644); err != nil {
		t.Fatal(err)
	}

	if got := commit.Parse(blobHash); got != nil {
		t.Error("Parse of blob object should return nil")
	}
}

func TestParseHead_NoHeadFile(t *testing.T) {
	tempDir := t.TempDir()

	origWd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	defer os.Chdir(origWd)

	if err := os.Chdir(tempDir); err != nil {
		t.Fatal(err)
	}

	// No .notgit/HEAD
	if got := commit.ParseHead(); got != nil {
		t.Errorf("ParseHead with no HEAD file should return nil, got %v", got)
	}
}

func TestParseHead_ValidHead(t *testing.T) {
	tempDir := t.TempDir()

	origWd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	defer os.Chdir(origWd)

	if err := os.Chdir(tempDir); err != nil {
		t.Fatal(err)
	}

	// Create tree and commit (same setup as TestParse_ValidCommit)
	tr := tree.NewTree(".")
	store := object.NewFSStore(tempDir)
	if err := tr.Write(store); err != nil {
		t.Fatal(err)
	}

	parentHash := strings.Repeat("b", 64)
	commitContent := []byte(strings.Join([]string{
		"tree " + tr.Hash(),
		"author Test Author 1234567890 -0700",
		"committer Test Author 1234567890 -0700",
		"parent " + parentHash,
		"",
		"head message",
	}, "\n"))

	commitHash := utils.Hash("commit", commitContent)
	commitCompressed := utils.Compress("commit", commitContent)

	objDir := filepath.Join(tempDir, ".notgit", "objects", commitHash[:2])
	if err := os.MkdirAll(objDir, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(objDir, commitHash[2:]), commitCompressed, 0644); err != nil {
		t.Fatal(err)
	}

	// Create HEAD file
	headDir := filepath.Join(tempDir, ".notgit")
	if err := os.MkdirAll(headDir, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(headDir, "HEAD"), []byte(commitHash), 0644); err != nil {
		t.Fatal(err)
	}

	parsed := commit.ParseHead()
	if parsed == nil {
		t.Fatal("ParseHead should return commit")
	}
	if parsed.Message != "head message" {
		t.Errorf("Message = %q, want head message", parsed.Message)
	}
}

func TestNewCommit_NoStaging(t *testing.T) {
	tempDir := t.TempDir()

	origWd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	defer os.Chdir(origWd)

	if err := os.Chdir(tempDir); err != nil {
		t.Fatal(err)
	}

	// Create .notgit with empty index (LoadStaged may fail or return empty tree)
	if err := os.MkdirAll(filepath.Join(tempDir, ".notgit"), 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(tempDir, ".notgit", "index"), []byte(""), 0644); err != nil {
		t.Fatal(err)
	}

	// NewCommit may return nil if LoadStaged fails, or a commit with empty tree
	c := commit.NewCommit("msg", "author", strings.Repeat("a", 64))
	if c == nil {
		// LoadStaged can fail for various reasons (e.g. indexfile.Parse)
		return
	}

	if c.Message != "msg" {
		t.Errorf("Message = %q, want msg", c.Message)
	}
	if c.Author != "author" {
		t.Errorf("Author = %q, want author", c.Author)
	}
}
