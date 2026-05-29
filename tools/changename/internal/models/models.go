package models

import "time"

type File struct {
	Path         string    `json:"path"`
	Name         string    `json:"name"`
	Ext          string    `json:"ext"`
	Size         int64     `json:"size"`
	ModifiedTime time.Time `json:"modifiedTime"`
	IsDir        bool      `json:"isDir"`
}

type RuleType string

const (
	RuleTypeSequence   RuleType = "sequence"
	RuleTypeReplace    RuleType = "replace"
	RuleTypeDate       RuleType = "date"
	RuleTypeCase       RuleType = "case"
	RuleTypeExtension  RuleType = "extension"
)

type CaseType string

const (
	CaseUpper      CaseType = "upper"
	CaseLower      CaseType = "lower"
	CaseTitle      CaseType = "title"
	CaseCamel      CaseType = "camel"
)

type DateSource string

const (
	DateSourceModified DateSource = "modified"
	DateSourceEXIF     DateSource = "exif"
)

type ConflictResolution string

const (
	ConflictOverwrite ConflictResolution = "overwrite"
	ConflictSkip      ConflictResolution = "skip"
	ConflictAutoNum   ConflictResolution = "autoNum"
)

type Rule struct {
	Type         RuleType         `json:"type"`
	Enabled      bool             `json:"enabled"`
	Sequence     SequenceRule     `json:"sequence"`
	Replace      ReplaceRule      `json:"replace"`
	Date         DateRule         `json:"date"`
	Case         CaseRule         `json:"case"`
	Extension    ExtensionRule    `json:"extension"`
}

type SequenceRule struct {
	Start    int    `json:"start"`
	Digits   int    `json:"digits"`
	Prefix   string `json:"prefix"`
	Suffix   string `json:"suffix"`
}

type ReplaceRule struct {
	Search     string `json:"search"`
	Replace    string `json:"replace"`
	UseRegex   bool   `json:"useRegex"`
}

type DateRule struct {
	Source   DateSource `json:"source"`
	Format   string     `json:"format"`
	Position string     `json:"position"`
	Separator string    `json:"separator"`
}

type CaseRule struct {
	Mode CaseType `json:"mode"`
}

type ExtensionRule struct {
	NewExt      string `json:"newExt"`
	UnifyCase   bool   `json:"unifyCase"`
	TargetCase  string `json:"targetCase"`
}

type PreviewItem struct {
	OriginalFile File   `json:"originalFile"`
	NewName      string `json:"newName"`
	NewPath      string `json:"newPath"`
	Conflict     bool   `json:"conflict"`
}

type Conflict struct {
	Filename string   `json:"filename"`
	Items    []string `json:"items"`
}

type RenameResult struct {
	Renamed int
	Skipped int
	UndoID  string
}

type UndoRecord struct {
	ID        string
	Timestamp time.Time
	Items     []UndoItem
}

type UndoItem struct {
	OldPath string
	NewPath string
}
