package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/rhysd/notes-cli"
)

func removeDirRec(dirpath, homepath string) {
	for dirpath != homepath {
		os.Remove(dirpath) // Remove directory if empty
		dirpath = filepath.Dir(dirpath)
	}
}

func recategorize(path string, config *notes.Config) error {
	note, err := notes.LoadNote(path, config)
	if err != nil && !errors.Is(err, &notes.MismatchCategoryError{}) {
		return err
	}
	rel, err := filepath.Rel(config.HomePath, filepath.Dir(path))
	if err != nil {
		return err
	}
	if filepath.ToSlash(rel) == note.Category {
		return nil
	}

	catpath := filepath.Join(config.HomePath, filepath.FromSlash(note.Category))
	if err := os.MkdirAll(catpath, 0755); err != nil {
		return err
	}

	newpath := filepath.Join(catpath, note.File)
	if _, err := os.Stat(newpath); err == nil {
		fmt.Printf("File %s already exists.\n", newpath)
		return nil
	}
	if err := os.Rename(path, newpath); err != nil {
		return err
	}
	removeDirRec(filepath.Dir(path), config.HomePath)

	fmt.Printf("Successfully moved %s to %s.\n", path, newpath)

	return nil
}

func main() {
	config, err := notes.NewConfig()
	if err != nil {
		panic(err)
	}

	cats, err := notes.CollectCategories(config, 0)
	if err != nil {
		panic(err)
	}
	for _, cat := range cats {
		for _, path := range cat.NotePaths {
			if err := recategorize(path, config); err != nil {
				panic(err)
			}
		}
	}
}
