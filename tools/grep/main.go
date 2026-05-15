package main

import (
	"grep/config"
	"grep/output"
	"grep/search"
)

func main() {
	cfg, err := config.Parse()
	if err != nil {
		config.PrintUsage()
		return
	}

	results, totalLines, totalFiles := search.Search(cfg)

	if cfg.Total {
		output.PrintTotal(totalLines, totalFiles)
		return
	}

	output.PrintResults(results, cfg)
}