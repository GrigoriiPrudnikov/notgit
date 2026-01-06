package commit

import (
	"notgit/internal/tree"
	"notgit/internal/utils"
	"strconv"
	"strings"
	"time"
)

// TODO: change all types to strings (hash)
type Commit struct {
	Time    int64
	Offset  string
	Author  string
	Message string
	Tree    *tree.Tree
	Parent  string
}

func NewCommit(message, author string, parent string) *Commit {
	root, err := tree.LoadStaged(".")
	if err != nil {
		return nil
	}

	t := time.Now()

	c := &Commit{
		Time:    t.Unix(),
		Offset:  t.Format("-0700"),
		Author:  author,
		Message: message,
		Parent:  parent,
		Tree:    root,
	}

	return c
}

func (c *Commit) GetContent() []byte {
	content := []string{
		"tree " + c.Tree.Hash(),
		"author " + c.Author + " " + strconv.FormatInt(c.Time, 10) + " " + c.Offset,
		"committer " + c.Author + " " + strconv.FormatInt(c.Time, 10) + " " + c.Offset,
	}

	content = append(content, "parent "+c.Parent)

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

func (c *Commit) Hash() string {
	content := []byte(c.GetContent())

	return utils.Hash("commit", content)
}
