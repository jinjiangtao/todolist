package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"
)

func (t *Tailer) followFile(file *os.File) error {
	scanner := bufio.NewScanner(file)
	scanner.Buffer(make([]byte, 4096), 1024*1024)

	for {
		if scanner.Scan() {
			t.mu.Lock()
			fmt.Printf("%s: %s\n", filepath.Base(t.filePath), scanner.Text())
			t.mu.Unlock()
		} else {
			if err := scanner.Err(); err != nil {
				return err
			}

			time.Sleep(100 * time.Millisecond)
			file.Seek(0, io.SeekCurrent)
			scanner = bufio.NewScanner(file)
			scanner.Buffer(make([]byte, 4096), 1024*1024)
		}
	}
}