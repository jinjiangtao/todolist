package config

import (
	"flag"
	"fmt"
	"os"
	"regexp"
)

type Config struct {
	Pattern     string
	Path        string
	Recursive   bool
	IgnoreCase  bool
	ShowLineNum bool
	Count       bool
	Total       bool
	ListFiles   bool
	InvertMatch bool
	Regex       *regexp.Regexp
}

func Parse() (*Config, error) {
	var cfg Config

	flag.BoolVar(&cfg.Recursive, "r", false, "递归目录搜索")
	flag.BoolVar(&cfg.IgnoreCase, "i", false, "忽略大小写")
	flag.BoolVar(&cfg.ShowLineNum, "n", false, "显示行号")
	flag.BoolVar(&cfg.Count, "c", false, "统计每个文件中匹配的行数")
	flag.BoolVar(&cfg.Total, "total", false, "输出全局总匹配行数、匹配文件总数")
	flag.BoolVar(&cfg.ListFiles, "l", false, "仅输出包含匹配内容的文件名")
	flag.BoolVar(&cfg.InvertMatch, "v", false, "反向匹配（输出不包含关键词的行）")

	flag.Parse()

	if flag.NArg() < 2 {
		return nil, fmt.Errorf("用法: grep [选项] <关键词> <文件/目录>")
	}

	cfg.Pattern = flag.Arg(0)
	cfg.Path = flag.Arg(1)

	if cfg.IgnoreCase {
		cfg.Pattern = "(?i)" + cfg.Pattern
	}

	re, err := regexp.Compile(cfg.Pattern)
	if err != nil {
		return nil, fmt.Errorf("正则表达式错误: %v", err)
	}
	cfg.Regex = re

	return &cfg, nil
}

func PrintUsage() {
	fmt.Println("用法: grep [选项] <关键词> <文件/目录>")
	fmt.Println("选项:")
	flag.PrintDefaults()
	os.Exit(1)
}