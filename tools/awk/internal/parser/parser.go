package parser

import (
	"fmt"
	"goawk/internal/lexer"
	"regexp"
	"strconv"
	"strings"
)

type Expression struct {
	Field1    int
	Field2    int
	Operator  string
	String1   string
	String2   string
	CompareOp string
	LogicOp   string
	Not       bool
}

type PatternAction struct {
	HasPattern  bool
	Pattern     Expression
	HasAction   bool
	PrintAll    bool
	PrintFields []int
	PrintRaw    string
}

type Parser struct {
	program PatternAction
}

func NewParser() *Parser {
	return &Parser{}
}

func (p *Parser) ParseProgram(programStr string) (PatternAction, error) {
	programStr = strings.TrimSpace(programStr)
	if programStr == "" {
		return p.program, fmt.Errorf("程序不能为空")
	}

	hasAction := false
	hasPattern := false

	if regexp.MustCompile(`^\$[0-9]+\s*$`).MatchString(programStr) {
		p.program.PrintFields = []int{}
		fieldStr := strings.TrimSpace(programStr[1:])
		fieldNum, err := strconv.Atoi(fieldStr)
		if err != nil {
			return p.program, fmt.Errorf("无效的字段编号: %s", fieldStr)
		}
		p.program.PrintFields = append(p.program.PrintFields, fieldNum)
		p.program.PrintAll = false
		p.program.HasAction = true
		p.program.HasPattern = true
		return p.program, nil
	}

	lbrace := strings.Index(programStr, "{")
	rbrace := strings.LastIndex(programStr, "}")

	if lbrace >= 0 && rbrace > lbrace {
		hasAction = true
		actionPart := strings.TrimSpace(programStr[lbrace+1 : rbrace])
		if actionPart != "" {
			if err := p.parseAction(actionPart); err != nil {
				return p.program, err
			}
		}
		if lbrace > 0 {
			hasPattern = true
			patternPart := strings.TrimSpace(programStr[:lbrace])
			if err := p.parsePattern(patternPart); err != nil {
				return p.program, err
			}
		}
	} else {
		hasPattern = true
		if err := p.parsePattern(programStr); err != nil {
			return p.program, err
		}
	}

	p.program.HasPattern = hasPattern
	p.program.HasAction = hasAction

	if !hasAction {
		p.program.PrintAll = true
	}

	return p.program, nil
}

func (p *Parser) parseAction(action string) error {
	action = strings.TrimSpace(action)
	l := lexer.NewLexer(action)
	token := l.NextToken()

	if token.Type == lexer.TokenPrint {
		printPart := strings.TrimSpace(action[len("print"):])
		if printPart == "" || printPart == ";" {
			p.program.PrintAll = true
			return nil
		}

		printPart = strings.Trim(printPart, ";")
		fields, rawStr, err := p.parsePrintArgs(printPart)
		if err != nil {
			return err
		}

		if rawStr != "" {
			p.program.PrintRaw = rawStr
		}
		p.program.PrintFields = fields
		if len(fields) > 0 || rawStr != "" {
			p.program.PrintAll = false
		} else {
			p.program.PrintAll = true
		}
	} else {
		p.program.PrintRaw = action
	}

	return nil
}

func (p *Parser) parsePrintArgs(args string) ([]int, string, error) {
	var fields []int
	var rawStr string

	parts := strings.Split(args, ",")
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if strings.HasPrefix(part, "$") && len(part) > 1 {
			fieldStr := part[1:]
			if fieldStr == "0" {
				fields = append(fields, 0)
			} else {
				fieldNum, err := strconv.Atoi(fieldStr)
				if err != nil {
					return nil, "", fmt.Errorf("无效的字段编号: %s", part)
				}
				fields = append(fields, fieldNum)
			}
		} else if strings.HasPrefix(part, "\"") || strings.HasPrefix(part, "'") {
			rawStr += strings.Trim(part, "\"'")
		} else if part != "" {
			rawStr += part
		}
	}

	return fields, rawStr, nil
}

func (p *Parser) parsePattern(pattern string) error {
	pattern = strings.TrimSpace(pattern)
	p.program.Pattern = Expression{}

	if pattern == "" {
		return nil
	}

	orParts := strings.Split(pattern, "||")
	if len(orParts) > 1 {
		p.program.Pattern.LogicOp = "||"
		for i, part := range orParts {
			part = strings.TrimSpace(part)
			if i == 0 {
				if err := p.parseSinglePattern(part, &p.program.Pattern, true); err != nil {
					return err
				}
			}
		}
		return nil
	}

	andParts := strings.Split(pattern, "&&")
	if len(andParts) > 1 {
		p.program.Pattern.LogicOp = "&&"
		for i, part := range andParts {
			part = strings.TrimSpace(part)
			if i == 0 {
				if err := p.parseSinglePattern(part, &p.program.Pattern, true); err != nil {
					return err
				}
			}
		}
		return nil
	}

	return p.parseSinglePattern(pattern, &p.program.Pattern, false)
}

func (p *Parser) parseSinglePattern(pattern string, expr *Expression, hasLogicOp bool) error {
	pattern = strings.TrimSpace(pattern)

	if strings.HasPrefix(pattern, "!") {
		expr.Not = true
		pattern = strings.TrimSpace(pattern[1:])
	}

	if strings.HasPrefix(pattern, "$") {
		eqIdx := -1
		for _, op := range []string{"==", "!=", "<=", ">=", "<", ">"} {
			idx := strings.Index(pattern, op)
			if idx > 0 && (eqIdx == -1 || idx < eqIdx) {
				eqIdx = idx
			}
		}

		if eqIdx > 0 {
			fieldStr := strings.TrimSpace(pattern[1:eqIdx])
			fieldNum, err := strconv.Atoi(fieldStr)
			if err != nil {
				return fmt.Errorf("无效的字段编号: %s", fieldStr)
			}

			op := ""
			for _, candidate := range []string{"==", "!=", "<=", ">=", "<", ">"} {
				if strings.HasPrefix(pattern[eqIdx:], candidate) {
					op = candidate
					break
				}
			}

			expr.CompareOp = op
			expr.Field1 = fieldNum

			value := strings.TrimSpace(pattern[eqIdx+len(op):])
			value = strings.Trim(value, "\"'")

			if strings.HasPrefix(value, "$") {
				fieldNum2, err := strconv.Atoi(value[1:])
				if err == nil {
					expr.Field2 = fieldNum2
					return nil
				}
			}

			if _, err := strconv.ParseFloat(value, 64); err == nil {
				expr.String2 = value
			} else {
				expr.String2 = value
			}

			return nil
		}
	}

	_, err := regexp.MatchString(pattern, "")
	if err == nil {
		expr.String1 = pattern
	}

	return nil
}
