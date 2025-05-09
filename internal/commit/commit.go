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
	root := tree.Staged()
	t := time.Now()

	c := &Commit{
		Time:      t.Unix(),
		Offset:    t.Format("-0700"),
		Tree:      root.Hash(),
		Author:    author,
		Committer: author,
		Message:   message,
	}

	for _, parent := range parents {
		p := Parse(parent)
		if p.Tree == c.Tree {
			return nil
		}
		c.Parents = append(c.Parents, p)
	}

	return c
}

func (c *Commit) Write() error {
	content := c.GetContent()
	header := fmt.Sprintf("commit %d\x00\n", len(content))

	compressed := utils.Compress(header, content)
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
	if len(hash) != 64 {
		return nil
	}

	c := &Commit{}

	dir, file := hash[0:2], hash[2:]
	path := filepath.Join(".notgit", "objects", dir, file)

	_, err := os.Stat(path)
	if os.IsNotExist(err) {
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

		if line == "" {
			break
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
			parent := Parse(values[0])
			if parent != nil {
				c.Parents = append(c.Parents, parent)
			}
		}
	}

	c.Message = lines[len(lines)-1]

	return c
}

func (c *Commit) GetContent() []byte {
	content := []string{
		"tree " + c.Tree,
		"author " + c.Author + " " + strconv.FormatInt(c.Time, 10) + " " + c.Offset,
		"committer " + c.Author + " " + strconv.FormatInt(c.Time, 10) + " " + c.Offset,
	}

	for _, parent := range c.Parents {
		content = append(content, "parent "+parent.Hash())
	}

	content = append(content, "", c.Message)
	return []byte(strings.Join(content, "\n"))
}

func parseNameTimeOffset(line string) (name string, time int64, offset string) {
	values := strings.Split(line, " ")
	n := len(values)
	name = strings.Join(values[1:n-2], " ")
	time, _ = strconv.ParseInt(values[n-2], 10, 64)
	offset = values[n-1]
	return
}
