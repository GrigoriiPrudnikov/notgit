package commit

import (
	"notgit/internal/indexfile"
	"notgit/internal/tree"
	"notgit/internal/utils"
	"strconv"
	"strings"
	"time"
)

type Commit struct {
	Time      int64
	Offset    string
	Author    string
	Committer string
	Message   string
	Tree      string
	Parents   []*Commit
}

func NewCommit(message, author string, parents []string) *Commit {
	index, err := indexfile.Parse()
	if err != nil {
		return nil
	}
	root := tree.Staged(index)
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

func (c *Commit) Hash() string {
	content := []byte(c.GetContent())

	return utils.Hash("commit", content)
}
