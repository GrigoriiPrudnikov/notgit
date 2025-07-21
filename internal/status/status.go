package status

// const (
// 	None = iota
// 	Added
// 	Modified
// 	Deleted
// )
//
// type change struct {
// 	Staged   int
// 	Unstaged int
// }
//
// func GetChanges() map[string]change {
// 	indexfile, err := indexfile.Parse()
// 	if err != nil {
// 		return nil
// 	}
//
// 	worktree := tree.().GetFiles()
// 	index := tree.Staged(indexfile).GetFiles()
// 	headCommit := commit.ParseHead()
// 	var head map[string]string
// 	if headCommit != nil {
// 		head = headCommit.Tree.GetFiles()
// 	}
// 	println("head")
// 	utils.PrintStruct(head)
// 	println("index")
// 	utils.PrintStruct(tree.Staged(indexfile))
//
// 	difference := map[string]change{}
// 	all := union(worktree, index, head)
//
// 	for _, filePath := range all {
// 		workEntry := worktree[filePath]
// 		indexEntry := index[filePath]
// 		headEntry := head[filePath]
//
// 		// no changes
// 		if workEntry == indexEntry && indexEntry == headEntry {
// 			continue
// 		}
//
// 		ch := change{}
//
// 		if headEntry != indexEntry {
// 			if headEntry == "" {
// 				ch.Staged = Added
// 			} else if indexEntry == "" {
// 				ch.Staged = Deleted
// 			} else {
// 				ch.Staged = Modified
// 			}
// 		}
//
// 		if indexEntry != workEntry {
// 			if workEntry == "" {
// 				ch.Unstaged = Deleted
// 			} else if indexEntry == "" {
// 				ch.Unstaged = Added
// 			} else {
// 				ch.Unstaged = Modified
// 			}
// 		}
//
// 		difference[filePath] = ch
// 	}
//
// 	return difference
// }
//
// func union(m1, m2, m3 map[string]string) []string {
// 	keys := make(map[string]struct{})
//
// 	for _, m := range []map[string]string{m1, m2, m3} {
// 		for k := range m {
// 			keys[k] = struct{}{}
// 		}
// 	}
//
// 	result := make([]string, 0, len(keys))
// 	for k := range keys {
// 		result = append(result, k)
// 	}
//
// 	return result
// }
