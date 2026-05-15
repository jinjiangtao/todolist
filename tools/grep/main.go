package main

import (
	"bufio"
	"flag"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

const (
	red    = "\x1b[31m"
	reset  = "\x1b[0m"
)

var (
	recursive    = flag.Bool("r", false, "递归目录搜索")
	ignoreCase   = flag.Bool("i", false, "忽略大小写")
	showLineNum  = flag.Bool("n", false, "显示行号")
	count        = flag.Bool("c", false, "统计每个文件中匹配的行数")
	total        = flag.Bool("total", false, "输出全局总匹配行数、匹配文件总数")
	listFiles    = flag.Bool("l", false, "仅输出包含匹配内容的文件名")
	invertMatch  = flag.Bool("v", false, "反向匹配（输出不包含关键词的行）")
)

type MatchResult struct {
	filePath     string
	lineNum      int
	lineContent  string
	matchCount   int
}

func main() {
	flag.Parse()

	if flag.NArg() < 2 {
		fmt.Println("用法: grep [选项] <关键词> <文件/目录>")
		fmt.Println("选项:")
		flag.PrintDefaults()
		os.Exit(1)
	}

	pattern := flag.Arg(0)
	path := flag.Arg(1)

	if *ignoreCase {
		pattern = "(?i)" + pattern
	}

	re, err := regexp.Compile(pattern)
	if err != nil {
		fmt.Printf("正则表达式错误: %v\n", err)
		os.Exit(1)
	}

	var results []MatchResult
	var totalMatchLines int
	var totalMatchFiles int

	info, err := os.Stat(path)
	if err != nil {
		fmt.Printf("错误: %v\n", err)
		os.Exit(1)
	}

	if info.IsDir() {
		if !*recursive {
			fmt.Println("错误: 需要使用 -r 选项递归搜索目录")
			os.Exit(1)
		}
		err = filepath.WalkDir(path, func(filePath string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if d.IsDir() {
				return nil
			}
			return processFile(filePath, re, &results, &totalMatchLines, &totalMatchFiles)
		})
		if err != nil {
			fmt.Printf("遍历目录时出错: %v\n", err)
			os.Exit(1)
		}
	} else {
		err = processFile(path, re, &results, &totalMatchLines, &totalMatchFiles)
		if err != nil {
			fmt.Printf("处理文件时出错: %v\n", err)
			os.Exit(1)
		}
	}

	if *total {
		fmt.Printf("总匹配行数: %d, 匹配文件数: %d\n", totalMatchLines, totalMatchFiles)
		return
	}

	printResults(results)
}

func processFile(filePath string, re *regexp.Regexp, results *[]MatchResult, totalMatchLines, totalMatchFiles *int) error {
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Printf("无法打开文件: %s, 错误: %v\n", filePath, err)
		return nil
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	lineNum := 0
	matchCount := 0
	var fileResults []MatchResult

	for scanner.Scan() {
		lineNum++
		line := scanner.Text()

		isMatch := re.MatchString(line)
		if *invertMatch {
			isMatch = !isMatch
		}

		if isMatch {
			matchCount++
			if !*count && !*listFiles {
				fileResults = append(fileResults, MatchResult{
					filePath:    filePath,
					lineNum:     lineNum,
					lineContent: line,
				})
			}
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Printf("读取文件时出错: %s, 错误: %v\n", filePath, err)
		return nil
	}

	if matchCount > 0 {
		*totalMatchFiles++
		*totalMatchLines += matchCount
	}

	if *count {
		*results = append(*results, MatchResult{
			filePath:   filePath,
			matchCount: matchCount,
		})
	} else if *listFiles && matchCount > 0 {
		*results = append(*results, MatchResult{
			filePath: filePath,
		})
	} else {
		*results = append(*results, fileResults...)
	}

	return nil
}

func printResults(results []MatchResult) {
	if *count {
		for _, res := range results {
			fmt.Printf("%s: %d\n", res.filePath, res.matchCount)
		}
		return
	}

	if *listFiles {
		seen := make(map[string]bool)
		for _, res := range results {
			if !seen[res.filePath] {
				fmt.Println(res.filePath)
				seen[res.filePath] = true
			}
		}
		return
	}

	fileGroups := make(map[string][]MatchResult)
	for _, res := range results {
		fileGroups[res.filePath] = append(fileGroups[res.filePath], res)
	}

	for filePath, fileResults := range fileGroups {
		for _, res := range fileResults {
			if *showLineNum {
				fmt.Printf("%s:%d:%s\n", filePath, res.lineNum, highlightMatch(res.lineContent))
			} else {
				fmt.Printf("%s:%s\n", filePath, highlightMatch(res.lineContent))
			}
		}
	}
}

func highlightMatch(line string) string {
	pattern := flag.Arg(0)
	if *ignoreCase {
		idx := strings.Index(strings.ToLower(line), strings.ToLower(pattern))
		if idx != -1 {
			return line[:idx] + red + line[idx:idx+len(pattern)] + reset + line[idx+len(pattern):]
		}
	} else {
		idx := strings.Index(line, pattern)
		if idx != -1 {
			return line[:idx] + red + line[idx:idx+len(pattern)] + reset + line[idx+len(pattern):]
		}
	}
	return line
}