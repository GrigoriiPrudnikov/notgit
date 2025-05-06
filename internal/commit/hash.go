package commit

import "notgit/utils"

func (c *Commit) Hash() string {
	content := []byte(c.GetContent())

	return utils.Hash("commit", content)
}
