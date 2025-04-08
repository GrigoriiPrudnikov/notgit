package tree

import (
	"notgit/internal/blob"
	"os"
	"path/filepath"
)

func Create(path string) (Tree, error) {
	wd, err := os.Getwd()
	if err != nil {
		return Tree{}, err
	}

	absolutePath := filepath.Join(wd, path)

	dir, err := os.ReadDir(absolutePath)
	if err != nil {
		return Tree{}, err
	}

	root := Tree{
		Path: filepath.Base(path),
	}

	// if files != nil {
	// 	// check if file in subtree
	// 	for _, file := range files {
	// 		parts := strings.SplitN(file, "/", 2)
	// 		if len(parts) == 2 {
	// 			subtree, err := Create(parts[0], []string{parts[1]})
	// 			if err != nil {
	// 				return Tree{}, err
	// 			}
	//
	// 			root.SubTrees = append(root.SubTrees, &subtree)
	// 			continue
	// 		}
	//
	// 		blob, err := blob.Create(filepath.Join(path, file))
	// 		if err != nil {
	// 			return Tree{}, err
	// 		}
	//
	// 		root.Blobs = append(root.Blobs, blob)
	// 		continue
	// 	}
	//
	// 	return Tree{}, errors.New("not implemented")
	// }

	for _, entry := range dir {
		if entry.IsDir() {
			subTree, err := Create(filepath.Join(path, entry.Name()))
			if err != nil {
				return Tree{}, err
			}

			root.SubTrees = append(root.SubTrees, &subTree)
			continue
		}

		file, err := blob.Create(filepath.Join(path, entry.Name()))
		if err != nil {
			return Tree{}, err
		}

		root.Blobs = append(root.Blobs, file)
	}

	if !blob.Exists(root.Hash) {
		err = root.write()
	}

	return root, err
}

func (t *Tree) write() error {
	return nil
}
