package commit

import (
	"notgit/utils"
	"strings"
)

func hash(c *Commit) {
	content := []string{
		"tree " + c.Tree,
		"author " + c.Author,
		"committer " + c.Author,
	}

	for _, parent := range c.Parents {
		content = append(content, "parent "+parent+"\n")
	}

	content = append(content, "\n"+c.Message)
	contentBytes := []byte(strings.Join(content, "\n"))

	c.Hash = utils.Hash("commit", contentBytes)
}
