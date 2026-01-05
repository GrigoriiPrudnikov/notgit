package repo

import (
	"errors"
	"fmt"
	"notgit/internal/blob"
	"notgit/internal/commit"
	"notgit/internal/indexfile"
	"notgit/internal/utils"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

type Repo struct {
	Wd      string
	Head    string
	Author  string
	Index   indexfile.IndexFile
	Objects []string
}

var (
	ErrNotGitRepo = errors.New("notgit: not a repository")
)

func ParseRepo(wd string) (*Repo, error) {
	r := &Repo{Wd: wd}

	if !utils.RepoInitialized(wd) {
		return nil, ErrNotGitRepo
	}

	b, err := os.ReadFile(filepath.Join(wd, ".notgit", "head"))
	if err != nil {
		return nil, errors.New("notgit: failed to read head file\n" + err.Error())
	}

	r.Head = string(b)

	index, err := indexfile.Parse(wd)
	if err != nil {
		return nil, errors.New("notgit: failed to parse index file\n" + err.Error())
	}

	r.Index = index

	objects, err := os.ReadDir(filepath.Join(wd, ".notgit", "objects"))
	if err != nil {
		return nil, errors.New("notgit: failed to read objects directory\n" + err.Error())
	}

	for _, sub := range objects {
		if !sub.IsDir() {
			return nil, errors.New("notgit: unexpected file in objects directory")
		}

		files, err := os.ReadDir(filepath.Join(wd, ".notgit", "objects", sub.Name()))
		if err != nil {
			return nil, errors.New("notgit: failed to read object subdirectory\n" + err.Error())
		}

		for _, file := range files {
			if file.IsDir() {
				return nil, errors.New("notgit: unexpected directory in objects subdirectory")
			}

			hash := sub.Name() + file.Name()
			r.Objects = append(r.Objects, hash)
		}
	}

	return r, nil
}

func (r *Repo) Init() error {
	if r.Wd == "" {
		return errors.New("no working directory provided")
	}

	dirs := []string{
		filepath.Join(r.Wd, ".notgit"),
		filepath.Join(r.Wd, ".notgit", "objects"),
	}

	for _, d := range dirs {
		err := os.Mkdir(d, 0755)
		if err != nil {
			return errors.New("notgit: failed to create directory\n" + err.Error())
		}
	}

	files := map[string]string{
		filepath.Join(r.Wd, ".notgit", "info"):        "",
		filepath.Join(r.Wd, ".notgit", "head"):        "",
		filepath.Join(r.Wd, ".notgit", "description"): "Unnamed repository; edit this file to name it for git",
	}

	for path, content := range files {
		err := os.WriteFile(path, []byte(content), 0644)
		if err != nil {
			return errors.New("notgit: failed to write file\n" + err.Error())
		}
	}

	return nil
}

func (r *Repo) Stage(path string) error {
	info, err := os.Stat(path)
	if err != nil {
		return errors.New("notgit: failed to get file info\n" + err.Error())
	}
	if os.IsNotExist(err) {
		return errors.New("notgit: file does not exist")
	}

	// Stage dir
	if info.IsDir() {
		dir, err := os.ReadDir(path)
		if err != nil {
			return errors.New("notgit: failed to read directory\n" + err.Error())
		}
		for _, entry := range dir {
			if err := r.Stage(filepath.Join(path, entry.Name())); err != nil {
				return err
			}
		}
		return nil
	}

	b, err := blob.NewBlob(path)
	if err != nil {
		return errors.New("notgit: failed to create blob\n" + err.Error())
	}
	err = b.Write()
	if err != nil {
		return errors.New("notgit: failed to write blob\n" + err.Error())
	}

	r.Index[path] = b.Hash()

	err = indexfile.Write(r.Index)
	if err != nil {
		return errors.New("notgit: failed to write index file\n" + err.Error())
	}

	return nil
}

func (r *Repo) Commit(message string) error {
	// homeDir, err := os.UserHomeDir()

	// cfg, err := config.Parse(filepath.Join(homeDir, ".notgitconfig"))
	// if err != nil {
	// 	return errors.New("notgit: please, run `notgit config` to configure your user name and email")
	// }

	author := "AUTHOR"

	cmt := commit.NewCommit(message, author, []string{r.Head})

	return cmt.Write()
}

func (r *Repo) Log() error {
	if r.Head == "" {
		return errors.New("notgit: no head commit")
	}

	head := commit.Parse(r.Head)

	current := head
	for current != nil {
		timestamp := time.Unix(current.Time, 0)

		offset := current.Offset
		offsetHours, _ := strconv.Atoi(offset[:3])              // "-07" -> -7
		offsetMins, _ := strconv.Atoi(offset[0:1] + offset[3:]) // "-00" -> 0
		totalOffset := offsetHours*3600 + offsetMins*60

		loc := time.FixedZone("commit zone", totalOffset)
		adjusted := timestamp.In(loc).Format("2006-01-02")

		headIndicator := ""
		if current == head {
			headIndicator = "\033[33m(\033[0m\033[35mHEAD\033[0m\033[33m)\033[0m "
		}

		fmt.Printf("\033[33m%s\033[0m \033[34m%s\033[0m %s%s\n", current.Hash()[:7], adjusted, headIndicator, current.Message)
		if len(current.Parents) > 0 {
			current = current.Parents[0]
		} else {
			current = nil
		}
	}

	return nil
}

func (r *Repo) Initialized() bool {
	return utils.RepoInitialized(r.Wd)
}
