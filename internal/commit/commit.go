package commit

import (
	"errors"
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
	root := tree.Root()
	time := time.Now()

	return &Commit{
		Time:    time.Unix(),
		Offset:  time.Format("-0700"),
		Tree:    root.Hash,
		Author:  author,
		Message: message,
		Parents: parents,
	}
}

func (c *Commit) Write() error {
	content := c.getContent()
	contentBytes := []byte(content)

	header := fmt.Sprintf("commit %d\x00\n", len(contentBytes))

	compressed := utils.Compress(header, contentBytes)
	hash(c)

	err := object.Write(c.Hash, compressed)
	if err != nil {
		return err
	}

	err = os.WriteFile(filepath.Join(".notgit", "HEAD"), []byte(c.Hash), 0644)

	return err
}

func ParseHead() (*Commit, error) {
	head, err := os.ReadFile(filepath.Join(".notgit", "HEAD"))
	if err != nil {
		return nil, err
	}

	return Parse(string(head))
}

func Parse(hash string) (*Commit, error) {
	c := Commit{}
	if len(hash) != 64 {
		return nil, errors.New("invalid hash")
	}

	dir, file := hash[0:2], hash[2:]
	path := filepath.Join(".notgit", "objects", dir, file)

	_, err := os.Stat(path)
	if err != nil {
		return nil, err
	}

	content, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	content, err = utils.Decompress(content)
	if err != nil {
		return nil, err
	}

	lines := strings.Split(string(content), "\n")
	for i, line := range lines {
		if i == 0 {
			continue
		}
		if i == len(lines) {
			c.Message = line
			break
		}

		section := strings.Split(line, " ")
		if len(section) != 2 {
			return nil, errors.New("invalid commit structure")
		}
		key, value := section[0], section[1]

		switch key {
		case "tree":
			c.Tree = value
		case "author":
			c.Author = value
		case "parent":
			c.Parents = append(c.Parents, value)
		}
	}

	fmt.Println(c)

	return &c, nil
}

func (c *Commit) getContent() []byte {
	content := []string{
		strconv.FormatInt(c.Time, 10) + " " + c.Offset,
		"tree " + c.Tree,
		"author " + c.Author,
		"committer " + c.Author,
	}

	for _, parent := range c.Parents {
		content = append(content, "parent "+parent)
	}

	content = append(content, "\n"+c.Message)
	return []byte(strings.Join(content, "\n"))
}
