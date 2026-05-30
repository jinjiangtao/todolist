package config

import (
	"fmt"
	"strings"
)

type Config struct {
	SourcePath     string
	DestPath       string
	MaxDepth       int
	CreateEmpty    bool
	PreserveAttrs  bool
	IgnorePatterns []string
	Simulate       bool
	Verbose        bool
}

func New(src, dst string) *Config {
	return &Config{
		SourcePath:     src,
		DestPath:       dst,
		MaxDepth:       0,
		CreateEmpty:    false,
		PreserveAttrs:  false,
		IgnorePatterns: []string{},
		Simulate:       false,
		Verbose:        false,
	}
}

func (c *Config) SetIgnorePatterns(patterns string) {
	if patterns == "" {
		c.IgnorePatterns = []string{}
		return
	}
	c.IgnorePatterns = strings.Split(patterns, ",")
	for i, p := range c.IgnorePatterns {
		c.IgnorePatterns[i] = strings.TrimSpace(p)
	}
}

func (c *Config) IsUnlimitedDepth() bool {
	return c.MaxDepth == 0
}

func (c *Config) GetDepthLimit() string {
	if c.IsUnlimitedDepth() {
		return "无限制"
	}
	return fmt.Sprintf("%d层", c.MaxDepth)
}
