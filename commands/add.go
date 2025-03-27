package commands

import (
	"bytes"
	"compress/zlib"
	"errors"
	"flag"
	"fmt"
	"notgit/utils"
	"os"
	"path/filepath"
	"slices"
)

func Add() error {
	wd, err := os.Getwd()
	if err != nil {
		return err
	}

	notgitDir := filepath.Join(wd, ".notgit")
	if _, err := os.Stat(notgitDir); os.IsNotExist(err) {
		// TODO: add handling for parent directories (fatal: not a git repository (or any of the parent directories): .git)
		return errors.New("not a git repository")
	}

	var all, force bool

	fs := flag.NewFlagSet("add", flag.ExitOnError)

	fs.BoolVar(&all, "all", false, "add all")
	fs.BoolVar(&all, "a", false, "add all")

	fs.BoolVar(&force, "force", false, "force")
	fs.BoolVar(&force, "f", false, "force")

	// fs.BoolVar(&update, "update", false, "update")
	// fs.BoolVar(&update, "u", false, "update")

	fs.Parse(os.Args[2:])
	args := fs.Args()

	if len(args) == 0 {
		return errors.New("no arguments")
	}

	for _, arg := range args {
		err := add(arg, force)

		if err != nil {
			return err
		}
	}

	return nil
}

var alwaysIgnoredPaths = []string{".git", ".notgit"}

// TODO: split to addFile and addDir
func add(path string, force bool) error {
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		return err
	}

	name := filepath.Base(path)
	if slices.Contains(alwaysIgnoredPaths, name) {
		fmt.Println("ignored", path)
		return nil
	}

	if utils.Ignored(path) && !force {
		fmt.Println("ignored", path)
		return nil
	}

	if info.IsDir() {
		if info.Name() == ".notgit" {
			return nil
		}

		entries, err := os.ReadDir(path)
		if err != nil {
			return err
		}

		for _, entry := range entries {
			entryPath := filepath.Join(path, entry.Name())
			add(entryPath, force)
		}

		return nil
	}

	workingDir, err := os.Getwd()
	if err != nil {
		return err
	}

	indexPath := filepath.Join(workingDir, ".notgit", "index")

	_, err = os.Stat(indexPath)
	if os.IsNotExist(err) {
		err = os.WriteFile(indexPath, []byte(""), 0644)
		if err != nil {
			return err
		}
	}

	b, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	err = utils.CreateBlob(b)
	if err != nil {
		return err
	}

	// dir := filepath.Join(workingDir, ".notgit", "objects", hex[:2])
	// fileName := hex[2:]

	var compressed bytes.Buffer
	w := zlib.NewWriter(&compressed)
	w.Write(b)
	w.Close()

	return nil
}
