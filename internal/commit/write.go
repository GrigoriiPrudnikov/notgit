package commit

import (
	"fmt"
	"notgit/internal/object"
	"notgit/internal/utils"
	"os"
	"path/filepath"
)

func (c *Commit) Write() error {
	err := c.Tree.Write()
	if err != nil {
		return err
	}

	content := c.GetContent()
	header := fmt.Sprintf("commit %d\x00\n", len(content))

	hash := c.Hash()
	compressed := utils.Compress(header, content)

	err = object.Write(hash, compressed)
	if err != nil {
		return err
	}

	return os.WriteFile(filepath.Join(".notgit", "HEAD"), []byte(hash), 0644)
}
