package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/jessevdk/go-flags"
	"github.com/rhysd/notes-cli"
)

var opts struct {
	FixFilename bool `short:"f" long:"fix-filename" description:"Fix filename by the title"`
}

func removeDirRec(dirpath, homepath string) {
	for dirpath != homepath {
		// Remove directory if empty
		if err := os.Remove(dirpath); err != nil {
			break
		}
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
	filename := note.File
	if opts.FixFilename {
		filename = strings.ToLower(note.Title)
		filename = strings.ReplaceAll(filename, " ", "_")
		filename = strings.ReplaceAll(filename, "/", "_")
		filename += ".md"
	}

	if filepath.ToSlash(rel) == note.Category && filename == note.File {
		return nil
	}

	catpath := filepath.Join(config.HomePath, filepath.FromSlash(note.Category))
	if err := os.MkdirAll(catpath, 0755); err != nil {
		return err
	}

	newpath := filepath.Join(catpath, filename)
	if _, err := os.Stat(newpath); err == nil {
		fmt.Printf("File %s already exists\n", newpath)
		return nil
	}
	if err := os.Rename(path, newpath); err != nil {
		return err
	}
	removeDirRec(filepath.Dir(path), config.HomePath)

	fmt.Printf("Successfully moved %s to %s\n", path, newpath)

	return nil
}

func main() {
	flags.NewParser(&opts, flags.IgnoreUnknown).Parse()

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
				fmt.Fprintf(os.Stderr, "failed to re-categorize %s\nerror: %s\n", path, err.Error())
				continue
			}
		}
	}
}
