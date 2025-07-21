package commit

import (
	"notgit/internal/tree"
	"notgit/internal/utils"
	"sort"
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
	Tree      *tree.Tree
	Parents   []*Commit
}

func NewCommit(message, author string, parents []string) *Commit {
	root, err := tree.LoadStaged()
	if err != nil {
		return nil
	}
	t := time.Now()

	c := &Commit{
		Time:      t.Unix(),
		Offset:    t.Format("-0700"),
		Tree:      root,
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
		"tree " + c.Tree.Hash(),
		"author " + c.Author + " " + strconv.FormatInt(c.Time, 10) + " " + c.Offset,
		"committer " + c.Author + " " + strconv.FormatInt(c.Time, 10) + " " + c.Offset,
	}

	sort.Slice(c.Parents, func(i, j int) bool {
		return c.Parents[i].Time > c.Parents[j].Time
	})

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
