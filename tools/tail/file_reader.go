package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

func (t *Tailer) readLastLines(file *os.File) error {
	stat, err := file.Stat()
	if err != nil {
		return err
	}

	fileSize := stat.Size()
	if fileSize == 0 {
		return nil
	}

	buf := make([]byte, 4096)
	var lines []string
	bytesRead := int64(0)
	lineCount := 0

	for {
		offset := int64(-1)
		if bytesRead+int64(len(buf)) < fileSize {
			offset = -int64(len(buf))
		} else {
			offset = -fileSize
		}

		if _, err := file.Seek(offset, io.SeekEnd); err != nil {
			return err
		}

		n, err := file.Read(buf)
		if err != nil && err != io.EOF {
			return err
		}

		bytesRead += int64(n)
		content := string(buf[:n])

		for i := len(content) - 1; i >= 0; i-- {
			if content[i] == '\n' {
				if i+1 < len(content) {
					lines = append(lines, content[i+1:])
					content = content[:i]
					lineCount++
					if lineCount >= t.lines {
						break
					}
				} else {
					content = content[:i]
				}
			}
		}

		if lineCount >= t.lines || bytesRead >= fileSize {
			if len(content) > 0 {
				lines = append(lines, content)
			}
			break
		}
	}

	if len(lines) > t.lines {
		lines = lines[:t.lines]
	}

	for i := len(lines) - 1; i >= 0; i-- {
		t.mu.Lock()
		fmt.Printf("%s: %s\n", filepath.Base(t.filePath), lines[i])
		t.mu.Unlock()
	}

	return nil
}