package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"
	"time"
)

var (
	linesFlag  = flag.Int("n", 10, "显示文件末尾的行数")
	followFlag = flag.Bool("f", false, "实时监控文件新增内容")
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
		return fmt.Errorf("无法打开文件 %s: %w", t.filePath, err)
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

func main() {
	flag.Parse()

	files := flag.Args()
	if len(files) == 0 {
		fmt.Println("用法: tail [-n 行数] [-f] 文件1 [文件2 ...]")
		fmt.Println("选项:")
		fmt.Println("  -n N    显示文件末尾N行，默认10行")
		fmt.Println("  -f      实时监控文件新增内容")
		os.Exit(1)
	}

	var wg sync.WaitGroup

	for _, filePath := range files {
		wg.Add(1)
		go func(fp string) {
			defer wg.Done()

			tailer := NewTailer(fp, *linesFlag, *followFlag)
			if err := tailer.Tail(); err != nil {
				fmt.Fprintf(os.Stderr, "错误: %v\n", err)
			}
		}(filePath)
	}

	wg.Wait()
}
