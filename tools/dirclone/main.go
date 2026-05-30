package main

import (
	"fmt"
	"os"
	"path/filepath"

	"dirclone/cli"
	"dirclone/clone"
	"dirclone/config"
	"dirclone/logger"
)

func main() {
	cfg, shouldExit := cli.Parse()
	if shouldExit {
		return
	}

	log := logger.New(cfg.Verbose)

	if cfg.Simulate {
		log.Simulate("不会实际创建文件/文件夹")
	}

	log.Scan(cfg.SourcePath)
	log.Depth(cfg.GetDepthLimit())
	log.Ignore(cfg.IgnorePatterns)

	scanner := clone.NewScanner(cfg.IgnorePatterns, cfg.MaxDepth)

	entries, err := scanner.Scan(cfg.SourcePath)
	if err != nil {
		log.Error(fmt.Sprintf("扫描目录失败: %v", err))
		os.Exit(1)
	}

	if !scanner.GetIgnorer().IsEmpty() && cfg.Verbose {
		log.Info(fmt.Sprintf("共发现 %d 个条目", len(entries)))
	}

	cloner := clone.NewCloner(
		cfg.DestPath,
		cfg.CreateEmpty,
		cfg.PreserveAttrs,
		cfg.Simulate,
		log,
	)

	if !cfg.Simulate {
		if err := os.MkdirAll(cfg.DestPath, 0755); err != nil {
			log.Error(fmt.Sprintf("创建目标目录失败: %v", err))
			os.Exit(1)
		}
	}

	result := cloner.Clone(entries)
	log.Complete(result.Dirs, result.Files)
}

func scanAndClone(cfg *config.Config, log *logger.Logger) error {
	absSrc, err := filepath.Abs(cfg.SourcePath)
	if err != nil {
		return fmt.Errorf("解析源路径失败: %w", err)
	}

	absDst, err := filepath.Abs(cfg.DestPath)
	if err != nil {
		return fmt.Errorf("解析目标路径失败: %w", err)
	}

	scanner := clone.NewScanner(cfg.IgnorePatterns, cfg.MaxDepth)

	entries, err := scanner.Scan(absSrc)
	if err != nil {
		return fmt.Errorf("扫描目录失败: %w", err)
	}

	cloner := clone.NewCloner(
		absDst,
		cfg.CreateEmpty,
		cfg.PreserveAttrs,
		cfg.Simulate,
		log,
	)

	if !cfg.Simulate {
		if err := os.MkdirAll(absDst, 0755); err != nil {
			return fmt.Errorf("创建目标目录失败: %w", err)
		}
	}

	cloner.Clone(entries)

	return nil
}
