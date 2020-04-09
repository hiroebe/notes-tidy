# notes-tidy

This is a subcommand for [notes-cli](https://github.com/rhysd/notes-cli),
which changes the directory structure automatically based on the category of each note.

## Installation

```
$ go get -u github.com/hiroebe/notes-tidy
```

## Usage

Change categories of notes, and run

```
$ notes tidy
```

It automatically moves the notes to appropriate places, creates directories if needed,
and removes empty directories.
