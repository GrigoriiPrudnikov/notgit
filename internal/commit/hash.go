package commit

import "notgit/utils"

func (c *Commit) Hash() string {
	content := c.getContent()
	contentBytes := []byte(content)

	return utils.Hash("commit", contentBytes)
}
