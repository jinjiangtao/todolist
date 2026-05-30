package clone

import (
	"strings"
)

type Ignorer struct {
	patterns []string
}

func NewIgnorer(patterns []string) *Ignorer {
	normalized := make([]string, len(patterns))
	for i, p := range patterns {
		normalized[i] = strings.ToLower(strings.TrimSpace(p))
	}
	return &Ignorer{patterns: normalized}
}

func (i *Ignorer) ShouldIgnore(name string) bool {
	lowerName := strings.ToLower(strings.TrimSpace(name))

	for _, pattern := range i.patterns {
		if pattern == lowerName {
			return true
		}
	}
	return false
}

func (i *Ignorer) IsEmpty() bool {
	return len(i.patterns) == 0
}
