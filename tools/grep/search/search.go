package search

import (
	"bufio"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"grep/config"
)

type MatchResult struct {
	FilePath    string
	LineNum     int
	LineContent string
	MatchCount  int
}

func Search(cfg *config.Config) ([]MatchResult, int, int) {
	var results []MatchResult
	var totalMatchLines int
	var totalMatchFiles int

	info, err := os.Stat(cfg.Path)
	if err != nil {
		fmt.Printf("错误: %v\n", err)
		os.Exit(1)
	}

	if info.IsDir() {
		if !cfg.Recursive {
			fmt.Println("错误: 需要使用 -r 选项递归搜索目录")
			os.Exit(1)
		}
		err = filepath.WalkDir(cfg.Path, func(filePath string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if d.IsDir() {
				return nil
			}
			return processFile(filePath, cfg, &results, &totalMatchLines, &totalMatchFiles)
		})
		if err != nil {
			fmt.Printf("遍历目录时出错: %v\n", err)
			os.Exit(1)
		}
	} else {
		err = processFile(cfg.Path, cfg, &results, &totalMatchLines, &totalMatchFiles)
		if err != nil {
			fmt.Printf("处理文件时出错: %v\n", err)
			os.Exit(1)
		}
	}

	return results, totalMatchLines, totalMatchFiles
}

func processFile(filePath string, cfg *config.Config, results *[]MatchResult, totalMatchLines, totalMatchFiles *int) error {
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

		isMatch := cfg.Regex.MatchString(line)
		if cfg.InvertMatch {
			isMatch = !isMatch
		}

		if isMatch {
			matchCount++
			if !cfg.Count && !cfg.ListFiles {
				fileResults = append(fileResults, MatchResult{
					FilePath:    filePath,
					LineNum:     lineNum,
					LineContent: line,
				})
			}
		}
	}

	if err := scanner.Err(); err != nil {
		if err.Error() == "bufio.Scanner: token too long" {
			return nil
		}
		fmt.Printf("读取文件时出错: %s, 错误: %v\n", filePath, err)
		return nil
	}

	if matchCount > 0 {
		*totalMatchFiles++
		*totalMatchLines += matchCount
	}

	if cfg.Count {
		*results = append(*results, MatchResult{
			FilePath:   filePath,
			MatchCount: matchCount,
		})
	} else if cfg.ListFiles && matchCount > 0 {
		*results = append(*results, MatchResult{
			FilePath: filePath,
		})
	} else {
		*results = append(*results, fileResults...)
	}

	return nil
}

func GetOriginalPattern(cfg *config.Config) string {
	pattern := cfg.Pattern
	if cfg.IgnoreCase && len(pattern) > 4 && pattern[:4] == "(?i)" {
		return pattern[4:]
	}
	return pattern
}