package clone

import (
	"os"
	"path/filepath"
)

type EntryInfo struct {
	Path    string
	IsDir   bool
	RelPath string
	Depth   int
}

type Scanner struct {
	ignorer  *Ignorer
	maxDepth int
}

func NewScanner(patterns []string, maxDepth int) *Scanner {
	return &Scanner{
		ignorer:  NewIgnorer(patterns),
		maxDepth: maxDepth,
	}
}

func (s *Scanner) Scan(srcPath string) ([]EntryInfo, error) {
	var entries []EntryInfo
	var scanErr error

	err := filepath.WalkDir(srcPath, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			if os.IsPermission(err) {
				return nil
			}
			return err
		}

		name := d.Name()

		if s.ignorer.ShouldIgnore(name) {
			if d.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}

		if path == srcPath {
			return nil
		}

		relPath, err := filepath.Rel(srcPath, path)
		if err != nil {
			return nil
		}

		depth := len(filepath.SplitList(relPath))

		if s.maxDepth > 0 && depth > s.maxDepth {
			if d.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}

		entries = append(entries, EntryInfo{
			Path:    path,
			IsDir:   d.IsDir(),
			RelPath: relPath,
			Depth:   depth,
		})

		return nil
	})

	if err != nil {
		scanErr = err
	}

	return entries, scanErr
}

func (s *Scanner) GetIgnorer() *Ignorer {
	return s.ignorer
}
