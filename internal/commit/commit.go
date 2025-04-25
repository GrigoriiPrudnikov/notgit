package commit

import (
	"notgit/internal/tree"
	"time"
)

func NewCommit(message, author string) *Commit {
	root := tree.Root()

	return &Commit{
		Date:    time.Now(),
		Author:  author,
		Message: message,
		Tree:    root.Hash,
	}
}
