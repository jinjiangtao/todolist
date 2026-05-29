package rules

import (
	"fmt"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"smart-renamer/internal/models"
)

type Engine struct {
}

func NewEngine() *Engine {
	return &Engine{}
}

func (e *Engine) Apply(file models.File, rules []models.Rule, index int) string {
	name := file.Name
	ext := file.Ext

	for _, rule := range rules {
		if !rule.Enabled {
			continue
		}

		switch rule.Type {
		case models.RuleTypeSequence:
			name = applySequence(name, rule.Sequence, index)
		case models.RuleTypeReplace:
			name = applyReplace(name, rule.Replace)
		case models.RuleTypeDate:
			name = applyDate(name, ext, rule.Date, file.ModifiedTime)
		case models.RuleTypeCase:
			name = applyCase(name, rule.Case)
		case models.RuleTypeExtension:
			ext = applyExtension(ext, rule.Extension)
		}
	}

	if ext != "" {
		return name + ext
	}
	return name
}

func applySequence(name string, rule models.SequenceRule, index int) string {
	num := rule.Start + index
	format := fmt.Sprintf("%%0%dd", rule.Digits)
	seq := fmt.Sprintf(format, num)
	return rule.Prefix + seq + rule.Suffix
}

func applyReplace(name string, rule models.ReplaceRule) string {
	if rule.Search == "" {
		return name
	}

	if rule.UseRegex {
		re, err := regexp.Compile(rule.Search)
		if err == nil {
			return re.ReplaceAllString(name, rule.Replace)
		}
	}

	return strings.ReplaceAll(name, rule.Search, rule.Replace)
}

func applyDate(name, ext string, rule models.DateRule, modTime time.Time) string {
	dateStr := modTime.Format(convertDateFormat(rule.Format))
	separator := rule.Separator

	if rule.Position == "prefix" {
		return dateStr + separator + name
	}
	return name + separator + dateStr
}

func convertDateFormat(format string) string {
	layouts := map[string]string{
		"YYYYMMDD": "20060102",
		"YYYY-MM-DD": "2006-01-02",
		"YYYY_MM_DD": "2006_01_02",
		"YYYYMMDD_HHMMSS": "20060102_150405",
		"YYYY-MM-DD HH:mm:ss": "2006-01-02 15:04:05",
	}
	if layout, ok := layouts[format]; ok {
		return layout
	}
	return "20060102"
}

func applyCase(name string, rule models.CaseRule) string {
	switch rule.Mode {
	case models.CaseUpper:
		return strings.ToUpper(name)
	case models.CaseLower:
		return strings.ToLower(name)
	case models.CaseTitle:
		return toTitleCase(name)
	case models.CaseCamel:
		return toCamelCase(name)
	default:
		return name
	}
}

func toTitleCase(s string) string {
	if len(s) == 0 {
		return s
	}
	return strings.ToUpper(string(s[0])) + strings.ToLower(s[1:])
}

func toCamelCase(s string) string {
	words := strings.FieldsFunc(s, func(r rune) bool {
		return r == ' ' || r == '_' || r == '-'
	})
	if len(words) == 0 {
		return s
	}
	result := strings.ToLower(words[0])
	for i := 1; i < len(words); i++ {
		result += strings.Title(words[i])
	}
	return result
}

func applyExtension(ext string, rule models.ExtensionRule) string {
	if rule.NewExt != "" {
		newExt := rule.NewExt
		if !strings.HasPrefix(newExt, ".") {
			newExt = "." + newExt
		}
		ext = newExt
	}

	if rule.UnifyCase {
		if rule.TargetCase == "lower" {
			ext = strings.ToLower(ext)
		} else {
			ext = strings.ToUpper(ext)
		}
	}

	return ext
}
