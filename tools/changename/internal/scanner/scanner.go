package scanner

import (
	"os"
	"path/filepath"
	"sync"

	"smart-renamer/internal/models"
)

type Scanner struct {
}

func New() *Scanner {
	return &Scanner{}
}

func (s *Scanner) Scan(directory string) ([]models.File, error) {
	var files []models.File
	var mu sync.Mutex
	var wg sync.WaitGroup

	err := filepath.Walk(directory, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if path == directory {
			return nil
		}

		wg.Add(1)
		go func(p string, fi os.FileInfo) {
			defer wg.Done()

			ext := filepath.Ext(p)
			name := fi.Name()
			if ext != "" {
				name = name[:len(name)-len(ext)]
			}

			file := models.File{
				Path:         p,
				Name:         name,
				Ext:          ext,
				Size:         fi.Size(),
				ModifiedTime: fi.ModTime(),
				IsDir:        fi.IsDir(),
			}

			mu.Lock()
			files = append(files, file)
			mu.Unlock()
		}(path, info)

		return nil
	})

	wg.Wait()

	if err != nil {
		return nil, err
	}

	return files, nil
}
