package output

import (
	"fmt"
	"strings"

	"grep/config"
	"grep/search"
)

const (
	red   = "\x1b[31m"
	reset = "\x1b[0m"
)

func PrintTotal(totalLines, totalFiles int) {
	fmt.Printf("总匹配行数: %d, 匹配文件数: %d\n", totalLines, totalFiles)
}

func PrintResults(results []search.MatchResult, cfg *config.Config) {
	if cfg.Count {
		printCountResults(results)
		return
	}

	if cfg.ListFiles {
		printListFilesResults(results)
		return
	}

	printMatchResults(results, cfg)
}

func printCountResults(results []search.MatchResult) {
	for _, res := range results {
		fmt.Printf("%s: %d\n", res.FilePath, res.MatchCount)
	}
}

func printListFilesResults(results []search.MatchResult) {
	seen := make(map[string]bool)
	for _, res := range results {
		if !seen[res.FilePath] {
			fmt.Println(res.FilePath)
			seen[res.FilePath] = true
		}
	}
}

func printMatchResults(results []search.MatchResult, cfg *config.Config) {
	fileGroups := make(map[string][]search.MatchResult)
	for _, res := range results {
		fileGroups[res.FilePath] = append(fileGroups[res.FilePath], res)
	}

	pattern := search.GetOriginalPattern(cfg)

	for filePath, fileResults := range fileGroups {
		for _, res := range fileResults {
			highlighted := highlightMatch(res.LineContent, pattern, cfg.IgnoreCase)
			if cfg.ShowLineNum {
				fmt.Printf("%s:%d:%s\n", filePath, res.LineNum, highlighted)
			} else {
				fmt.Printf("%s:%s\n", filePath, highlighted)
			}
		}
	}
}

func highlightMatch(line, pattern string, ignoreCase bool) string {
	if ignoreCase {
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