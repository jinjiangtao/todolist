package clone

import (
	"fmt"
	"os"
	"path/filepath"

	"dirclone/logger"
)

type Result struct {
	Dirs  int
	Files int
}

type Cloner struct {
	dstPath       string
	createEmpty   bool
	preserveAttrs bool
	simulate      bool
	log           *logger.Logger
}

func NewCloner(dstPath string, createEmpty, preserveAttrs, simulate bool, log *logger.Logger) *Cloner {
	return &Cloner{
		dstPath:       dstPath,
		createEmpty:   createEmpty,
		preserveAttrs: preserveAttrs,
		simulate:      simulate,
		log:           log,
	}
}

func (c *Cloner) Clone(entries []EntryInfo) Result {
	result := Result{}

	for _, entry := range entries {
		c.processEntry(entry, &result)
	}

	return result
}

func (c *Cloner) processEntry(entry EntryInfo, result *Result) {
	targetPath := filepath.Join(c.dstPath, entry.RelPath)

	if entry.IsDir {
		c.createDir(entry, targetPath, result)
	} else if c.createEmpty {
		c.createFile(entry, targetPath, result)
	}
}

func (c *Cloner) createDir(entry EntryInfo, targetPath string, result *Result) {
	if c.simulate {
		c.log.CreateDir(targetPath)
		result.Dirs++
		return
	}

	err := os.MkdirAll(targetPath, 0755)
	if err != nil {
		if os.IsPermission(err) {
			c.log.Skip(targetPath, "权限不足")
			return
		}
		c.log.Warn(fmt.Sprintf("创建目录失败 %s: %v", targetPath, err))
		return
	}

	c.log.CreateDir(targetPath)
	result.Dirs++

	if c.preserveAttrs {
		c.preserveDirAttrs(entry.Path, targetPath)
	}
}

func (c *Cloner) createFile(entry EntryInfo, targetPath string, result *Result) {
	if c.simulate {
		c.log.CreateFile(targetPath)
		result.Files++
		return
	}

	file, err := os.Create(targetPath)
	if err != nil {
		if os.IsPermission(err) {
			c.log.Skip(targetPath, "权限不足")
			return
		}
		c.log.Warn(fmt.Sprintf("创建文件失败 %s: %v", targetPath, err))
		return
	}
	file.Close()

	c.log.CreateFile(targetPath)
	result.Files++

	if c.preserveAttrs {
		c.preserveFileAttrs(entry.Path, targetPath)
	}
}

func (c *Cloner) preserveDirAttrs(srcPath, dstPath string) {
	srcInfo, err := os.Stat(srcPath)
	if err != nil {
		return
	}

	os.Chmod(dstPath, srcInfo.Mode())
	os.Chtimes(dstPath, srcInfo.ModTime(), srcInfo.ModTime())
	c.log.PreservedAttrs(dstPath)
}

func (c *Cloner) preserveFileAttrs(srcPath, dstPath string) {
	srcInfo, err := os.Stat(srcPath)
	if err != nil {
		return
	}

	os.Chmod(dstPath, srcInfo.Mode())
	os.Chtimes(dstPath, srcInfo.ModTime(), srcInfo.ModTime())
	c.log.PreservedAttrs(dstPath)
}

func ScanAndClone(src, dst string, entries []EntryInfo, log *logger.Logger) (Result, error) {
	if err := os.MkdirAll(dst, 0755); err != nil {
		return Result{}, fmt.Errorf("无法创建目标目录: %w", err)
	}

	cloner := NewCloner(dst, false, false, false, log)
	result := cloner.Clone(entries)

	return result, nil
}
