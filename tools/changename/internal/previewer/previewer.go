package previewer

import (
	"path/filepath"
	"sync"

	"smart-renamer/internal/models"
	"smart-renamer/internal/rules"
)

type Previewer struct {
	rulesEngine *rules.Engine
}

func New() *Previewer {
	return &Previewer{
		rulesEngine: rules.NewEngine(),
	}
}

func (p *Previewer) Generate(files []models.File, ruleList []models.Rule) ([]models.PreviewItem, []models.Conflict, error) {
	var items []models.PreviewItem
	var mu sync.Mutex
	var wg sync.WaitGroup

	for i, file := range files {
		wg.Add(1)
		go func(idx int, f models.File) {
			defer wg.Done()

			newName := p.rulesEngine.Apply(f, ruleList, idx)
			newPath := filepath.Join(filepath.Dir(f.Path), newName)

			item := models.PreviewItem{
				OriginalFile: f,
				NewName:      newName,
				NewPath:      newPath,
				Conflict:     false,
			}

			mu.Lock()
			items = append(items, item)
			mu.Unlock()
		}(i, file)
	}

	wg.Wait()

	conflicts := p.detectConflicts(items)

	return items, conflicts, nil
}

func (p *Previewer) detectConflicts(items []models.PreviewItem) []models.Conflict {
	pathMap := make(map[string][]string)

	for _, item := range items {
		pathMap[item.NewPath] = append(pathMap[item.NewPath], item.OriginalFile.Name+item.OriginalFile.Ext)
	}

	var conflicts []models.Conflict
	for path, fileList := range pathMap {
		if len(fileList) > 1 {
			conflicts = append(conflicts, models.Conflict{
				Filename: filepath.Base(path),
				Items:    fileList,
			})

			for i := range items {
				if items[i].NewPath == path {
					items[i].Conflict = true
				}
			}
		}
	}

	return conflicts
}
