package main

import (
	"os"
	"sync"
)

type Tailer struct {
	filePath string
	lines    int
	follow   bool
	mu       sync.Mutex
}

func NewTailer(filePath string, lines int, follow bool) *Tailer {
	return &Tailer{
		filePath: filePath,
		lines:    lines,
		follow:   follow,
	}
}

func (t *Tailer) Tail() error {
	file, err := os.Open(t.filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	if err := t.readLastLines(file); err != nil {
		return err
	}

	if t.follow {
		if err := t.followFile(file); err != nil {
			return err
		}
	}

	return nil
}