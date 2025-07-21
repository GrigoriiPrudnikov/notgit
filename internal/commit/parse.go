package commit

import (
	"notgit/internal/tree"
	"notgit/internal/utils"
	"os"
	"path/filepath"
	"strings"
)

func ParseHead() *Commit {
	head, err := os.ReadFile(filepath.Join(".notgit", "HEAD"))
	if err != nil {
		return nil
	}

	return Parse(string(head))
}

func Parse(hash string) *Commit {
	if len(hash) != 64 {
		return nil
	}

	c := &Commit{}

	dir, file := hash[0:2], hash[2:]
	path := filepath.Join(".notgit", "objects", dir, file)

	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		return nil
	}

	content, err := os.ReadFile(path)
	if err != nil {
		return nil
	}
	content, err = utils.Decompress(content)
	if err != nil {
		return nil
	}

	if !strings.HasPrefix(string(content), "commit") {
		return nil
	}

	lines := strings.Split(string(content), "\n")
	for i, line := range lines {
		if i == 0 {
			continue
		}

		if line == "" {
			break
		}

		prefix := strings.Split(line, " ")[0]
		values := strings.Split(line, " ")[1:]

		switch prefix {
		case "author", "committer":
			name, time, offset := parseNameTimeOffset(line)
			c.Time = time
			c.Offset = offset

			if prefix == "author" {
				c.Author = name
				continue
			}

			c.Committer = name

		case "tree":
			t, err := tree.Parse(values[0])
			if err != nil {
				println("here")
				println(err.Error())
				return nil
			}
			c.Tree = t

		case "parent":
			parent := Parse(values[0])
			if parent != nil {
				c.Parents = append(c.Parents, parent)
			}
		}
	}

	c.Message = lines[len(lines)-1]

	return c
}
