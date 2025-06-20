package status

import (
	"notgit/internal/commit"
	"notgit/internal/indexfile"
	"notgit/internal/tree"
)

const (
	None = iota
	Added
	Modified
	Deleted
)

type change struct {
	Staged   int
	Unstaged int
}

func GetChanges() map[string]change {
	indexfile, err := indexfile.Parse()
	if err != nil {
		return nil
	}

	worktree := tree.Root().GetFiles()
	index := tree.Staged(indexfile).GetFiles()
	head := commit.ParseHeadTree().GetFiles()

	difference := map[string]change{}
	all := union(worktree, index, head)

	for _, filePath := range all {
		workEntry := worktree[filePath]
		indexEntry := index[filePath]
		headEntry := head[filePath]

		// no changes
		if workEntry == indexEntry && workEntry == headEntry {
			continue
		}

		ch := change{}

		if headEntry != indexEntry {
			if headEntry == "" {
				ch.Staged = Added
			} else if indexEntry == "" {
				ch.Staged = Deleted
			} else {
				ch.Staged = Modified
			}
		}

		if indexEntry != workEntry {
			if workEntry == "" {
				ch.Unstaged = Deleted
			} else if indexEntry == "" {
				ch.Unstaged = Added
			} else {
				ch.Unstaged = Modified
			}
		}

		difference[filePath] = ch
	}

	return difference
}

func union(m1, m2, m3 map[string]string) []string {
	keys := make(map[string]struct{})

	for _, m := range []map[string]string{m1, m2, m3} {
		for k := range m {
			keys[k] = struct{}{}
		}
	}

	result := make([]string, 0, len(keys))
	for k := range keys {
		result = append(result, k)
	}

	return result
}
