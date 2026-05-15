package main

import (
	"flag"
	"fmt"
	"os"
	"sync"
)

var (
	linesFlag  = flag.Int("n", 10, "显示文件末尾的行数")
	followFlag = flag.Bool("f", false, "实时监控文件新增内容")
)

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