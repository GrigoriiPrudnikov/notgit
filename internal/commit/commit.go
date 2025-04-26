package commit

import (
	"fmt"
	"notgit/internal/object"
	"notgit/internal/tree"
	"notgit/utils"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

func NewCommit(message, author string, parents []string) *Commit {
	root := tree.Root()
	time := time.Now()

	return &Commit{
		Time:    time.Unix(),
		Offset:  time.Format("-0700"),
		Tree:    root.Hash,
		Author:  author,
		Message: message,
		Parents: parents,
	}
}

func (c *Commit) Write() error {
	content := []string{
		strconv.FormatInt(c.Time, 10) + " " + c.Offset,
		"tree " + c.Tree,
		"author " + c.Author,
		"committer " + c.Author,
	}

	for _, parent := range c.Parents {
		content = append(content, "parent "+parent)
	}

	content = append(content, "\n", c.Message)
	contentBytes := []byte(strings.Join(content, "\n"))

	header := fmt.Sprintf("tree %d\x00\n", len(contentBytes))

	fmt.Println(header)
	fmt.Println(string(contentBytes))

	compressed := utils.Compress(header, contentBytes)
	hash(c)

	err := object.Write(c.Hash, compressed)
	if err != nil {
		return err
	}

	err = os.WriteFile(filepath.Join(".notgit", "HEAD"), []byte(c.Hash), 0644)

	return err
}
