package renamer

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/google/uuid"
	"smart-renamer/internal/models"
)

type Renamer struct {
	undoHistory map[string]models.UndoRecord
	mu          sync.Mutex
}

func New() *Renamer {
	return &Renamer{
		undoHistory: make(map[string]models.UndoRecord),
	}
}

func (r *Renamer) Execute(items []models.PreviewItem, conflictRes string) (*models.RenameResult, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	undoID := uuid.New().String()
	undoItems := []models.UndoItem{}
	renamed := 0
	skipped := 0

	processedPaths := make(map[string]bool)

	for _, item := range items {
		oldPath := item.OriginalFile.Path
		newPath := item.NewPath

		if oldPath == newPath {
			continue
		}

		finalPath := newPath
		shouldSkip := false

		if item.Conflict || fileExists(newPath) {
			switch models.ConflictResolution(conflictRes) {
			case models.ConflictSkip:
				skipped++
				shouldSkip = true
			case models.ConflictOverwrite:
				if err := os.Remove(newPath); err != nil {
					skipped++
					shouldSkip = true
				}
			case models.ConflictAutoNum:
				finalPath = r.generateUniquePath(newPath, processedPaths)
			}
		}

		if shouldSkip {
			continue
		}

		if err := os.Rename(oldPath, finalPath); err != nil {
			skipped++
			continue
		}

		undoItems = append(undoItems, models.UndoItem{
			OldPath: finalPath,
			NewPath: oldPath,
		})
		processedPaths[finalPath] = true
		renamed++
	}

	if len(undoItems) > 0 {
		r.undoHistory[undoID] = models.UndoRecord{
			ID:        undoID,
			Timestamp: time.Now(),
			Items:     undoItems,
		}
	}

	return &models.RenameResult{
		Renamed: renamed,
		Skipped: skipped,
		UndoID:  undoID,
	}, nil
}

func (r *Renamer) Undo(undoID string) (int, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	record, exists := r.undoHistory[undoID]
	if !exists {
		return 0, fmt.Errorf("undo record not found")
	}

	restored := 0
	for _, item := range record.Items {
		if err := os.Rename(item.OldPath, item.NewPath); err == nil {
			restored++
		}
	}

	delete(r.undoHistory, undoID)

	return restored, nil
}

func (r *Renamer) generateUniquePath(path string, processed map[string]bool) string {
	dir := filepath.Dir(path)
	ext := filepath.Ext(path)
	base := path[:len(path)-len(ext)]
	counter := 1

	for {
		newPath := fmt.Sprintf("%s_%d%s", base, counter, ext)
		if !fileExists(newPath) && !processed[newPath] {
			return newPath
		}
		counter++
	}
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}
