package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"smart-renamer/internal/models"
	"smart-renamer/internal/scanner"
	"smart-renamer/internal/rules"
	"smart-renamer/internal/previewer"
	"smart-renamer/internal/renamer"
)

type App struct {
	ctx         context.Context
	scanner     *scanner.Scanner
	rulesEngine *rules.Engine
	previewer   *previewer.Previewer
	renamer     *renamer.Renamer
	mu          sync.Mutex
}

func NewApp() *App {
	return &App{
		scanner:     scanner.New(),
		rulesEngine: rules.NewEngine(),
		previewer:   previewer.New(),
		renamer:     renamer.New(),
	}
}

func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

func (a *App) shutdown(ctx context.Context) {
}

type ScanRequest struct {
	Directory string `json:"directory"`
}

type ScanResponse struct {
	Success bool           `json:"success"`
	Files   []models.File  `json:"files"`
	Error   string         `json:"error,omitempty"`
}

func (a *App) ScanDirectory(requestJSON string) string {
	var req ScanRequest
	if err := json.Unmarshal([]byte(requestJSON), &req); err != nil {
		return toJSON(ScanResponse{Success: false, Error: "Invalid request format"})
	}

	if req.Directory == "" {
		return toJSON(ScanResponse{Success: false, Error: "Directory is required"})
	}

	if _, err := os.Stat(req.Directory); os.IsNotExist(err) {
		return toJSON(ScanResponse{Success: false, Error: "Directory does not exist"})
	}

	files, err := a.scanner.Scan(req.Directory)
	if err != nil {
		return toJSON(ScanResponse{Success: false, Error: err.Error()})
	}

	return toJSON(ScanResponse{Success: true, Files: files})
}

type PreviewRequest struct {
	Files  []models.File   `json:"files"`
	Rules  []models.Rule   `json:"rules"`
}

type PreviewResponse struct {
	Success bool                  `json:"success"`
	Results []models.PreviewItem  `json:"results"`
	Conflicts []models.Conflict   `json:"conflicts"`
	Error   string                `json:"error,omitempty"`
}

func (a *App) GeneratePreview(requestJSON string) string {
	var req PreviewRequest
	if err := json.Unmarshal([]byte(requestJSON), &req); err != nil {
		return toJSON(PreviewResponse{Success: false, Error: "Invalid request format"})
	}

	if len(req.Files) == 0 {
		return toJSON(PreviewResponse{Success: false, Error: "No files to process"})
	}

	results, conflicts, err := a.previewer.Generate(req.Files, req.Rules)
	if err != nil {
		return toJSON(PreviewResponse{Success: false, Error: err.Error()})
	}

	return toJSON(PreviewResponse{
		Success:   true,
		Results:   results,
		Conflicts: conflicts,
	})
}

type RenameRequest struct {
	Items       []models.PreviewItem `json:"items"`
	ConflictRes string               `json:"conflictResolution"`
}

type RenameResponse struct {
	Success bool               `json:"success"`
	Renamed int                `json:"renamed"`
	Skipped int                `json:"skipped"`
	UndoID  string             `json:"undoId,omitempty"`
	Error   string             `json:"error,omitempty"`
}

func (a *App) ExecuteRename(requestJSON string) string {
	var req RenameRequest
	if err := json.Unmarshal([]byte(requestJSON), &req); err != nil {
		return toJSON(RenameResponse{Success: false, Error: "Invalid request format"})
	}

	if len(req.Items) == 0 {
		return toJSON(RenameResponse{Success: false, Error: "No items to rename"})
	}

	result, err := a.renamer.Execute(req.Items, req.ConflictRes)
	if err != nil {
		return toJSON(RenameResponse{Success: false, Error: err.Error()})
	}

	return toJSON(RenameResponse{
		Success: true,
		Renamed: result.Renamed,
		Skipped: result.Skipped,
		UndoID:  result.UndoID,
	})
}

type UndoRequest struct {
	UndoID string `json:"undoId"`
}

type UndoResponse struct {
	Success bool   `json:"success"`
	Restored int   `json:"restored"`
	Error   string `json:"error,omitempty"`
}

func (a *App) UndoRename(requestJSON string) string {
	var req UndoRequest
	if err := json.Unmarshal([]byte(requestJSON), &req); err != nil {
		return toJSON(UndoResponse{Success: false, Error: "Invalid request format"})
	}

	if req.UndoID == "" {
		return toJSON(UndoResponse{Success: false, Error: "Undo ID is required"})
	}

	restored, err := a.renamer.Undo(req.UndoID)
	if err != nil {
		return toJSON(UndoResponse{Success: false, Error: err.Error()})
	}

	return toJSON(UndoResponse{
		Success:  true,
		Restored: restored,
	})
}

func (a *App) SelectDirectory() string {
	return ""
}

func toJSON(v interface{}) string {
	data, err := json.Marshal(v)
	if err != nil {
		return fmt.Sprintf(`{"success":false,"error":"%s"}`, err.Error())
	}
	return string(data)
}
