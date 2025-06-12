package commit

import (
	"fmt"
	"notgit/internal/object"
	"notgit/internal/utils"
	"os"
	"path/filepath"
)

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
