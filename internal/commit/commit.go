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
	t := time.Now()

	return &Commit{
		Time:      t.Unix(),
		Offset:    t.Format("-0700"),
		Tree:      root.Hash,
		Author:    author,
		Committer: author,
		Message:   message,
		Parents:   parents,
	}
}

func (c *Commit) Write() error {
	content := c.getContent()
	contentBytes := []byte(content)

	header := fmt.Sprintf("commit %d\x00\n", len(contentBytes))

	compressed := utils.Compress(header, contentBytes)
	hash := c.Hash()

	err := object.Write(hash, compressed)
	if err != nil {
		return err
	}

	err = os.WriteFile(filepath.Join(".notgit", "HEAD"), []byte(hash), 0644)

	return err
}

func ParseHead() *Commit {
	head, err := os.ReadFile(filepath.Join(".notgit", "HEAD"))
	if err != nil {
		return nil
	}

	return Parse(string(head))
}

func Parse(hash string) *Commit {
	if len(hash) == 0 || len(hash) != 64 {
		return nil
	}

	c := &Commit{}

	dir, file := hash[0:2], hash[2:]
	path := filepath.Join(".notgit", "objects", dir, file)

	_, err := os.Stat(path)
	if err != nil {
		return nil
	}

	content, err := os.ReadFile(path)
	if err != nil {
		return nil
	}
	content, err = utils.Decompress(content)
	if err != nil {
		return nil
	}

	if !strings.HasPrefix(string(content), "commit") {
		return nil
	}

	lines := strings.Split(string(content), "\n")
	for i, line := range lines {
		if i == 0 {
			continue
		}

		if i == len(lines)-1 {
			c.Message = line
		}

		prefix := strings.Split(line, " ")[0]
		values := strings.Split(line, " ")[1:]

		switch prefix {
		case "author", "committer":
			name, time, offset := parseNameTimeOffset(line)
			c.Time = time
			c.Offset = offset

			if prefix == "author" {
				c.Author = name
				continue
			}

			c.Committer = name

		case "tree":
			c.Tree = values[0]

		case "parent":
			c.Parents = append(c.Parents, values[0])
		}
	}

	return c
}

func (c *Commit) getContent() []byte {
	content := []string{
		"tree " + c.Tree,
		"author " + c.Author + " " + strconv.FormatInt(c.Time, 10) + " " + c.Offset,
		"committer " + c.Author + " " + strconv.FormatInt(c.Time, 10) + " " + c.Offset,
	}

	for _, parent := range c.Parents {
		content = append(content, "parent "+parent)
	}

	content = append(content, "\n"+c.Message)
	return []byte(strings.Join(content, "\n"))
}

func parseNameTimeOffset(line string) (name string, time int64, offset string) {
	values := strings.Split(line, " ")
	n := len(values)
	name = strings.Join(values[:n-2], " ")
	time, _ = strconv.ParseInt(values[n-2], 10, 64)
	offset = values[n-1]
	return
}
