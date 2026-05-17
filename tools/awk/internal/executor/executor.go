package executor

import (
	"bufio"
	"fmt"
	"goawk/internal/parser"
	"io"
	"regexp"
	"strconv"
	"strings"
)

type AWKExecutor struct {
	fs         string
	program    parser.PatternAction
	fieldRegex *regexp.Regexp
}

func NewAWKExecutor(fs string) *AWKExecutor {
	exe := &AWKExecutor{fs: fs}
	if fs == " " || fs == "" {
		exe.fs = "\\s+"
		exe.fieldRegex = regexp.MustCompile(`\s+`)
	} else if len(fs) == 1 {
		exe.fs = regexp.QuoteMeta(fs)
		exe.fieldRegex = regexp.MustCompile(`\` + fs)
	} else {
		exe.fieldRegex = regexp.MustCompile(fs)
	}
	return exe
}

func (exe *AWKExecutor) SetProgram(program parser.PatternAction) {
	exe.program = program
}

func (exe *AWKExecutor) Execute(reader io.Reader, writer io.Writer) error {
	scanner := bufio.NewScanner(reader)
	lineNum := 0

	for scanner.Scan() {
		lineNum++
		line := scanner.Text()
		fields := exe.splitLine(line)

		if exe.matchPattern(fields, line) {
			exe.executeAction(fields, line, writer)
		}
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("读取输入错误: %v", err)
	}

	return nil
}

func (exe *AWKExecutor) splitLine(line string) []string {
	if exe.fs == "\\s+" {
		return exe.fieldRegex.Split(line, -1)
	}
	parts := exe.fieldRegex.Split(line, -1)
	var result []string
	for _, p := range parts {
		if p != "" || len(result) == 0 {
			result = append(result, p)
		}
	}
	return result
}

func (exe *AWKExecutor) matchPattern(fields []string, line string) bool {
	if !exe.program.HasPattern {
		return true
	}

	expr := &exe.program.Pattern

	if expr.LogicOp == "||" {
		return true
	}

	if expr.LogicOp == "&&" {
		return exe.evaluateExpr(expr, fields, line)
	}

	return exe.evaluateExpr(expr, fields, line)
}

func (exe *AWKExecutor) evaluateExpr(expr *parser.Expression, fields []string, line string) bool {
	if expr.Not {
		return !exe.evaluateSingleExpr(expr, fields, line)
	}
	return exe.evaluateSingleExpr(expr, fields, line)
}

func (exe *AWKExecutor) evaluateSingleExpr(expr *parser.Expression, fields []string, line string) bool {
	if expr.CompareOp == "" {
		if expr.String1 != "" {
			matched, _ := regexp.MatchString(expr.String1, line)
			return matched
		}
		return true
	}

	field1Val := ""
	if expr.Field1 == 0 {
		field1Val = line
	} else if expr.Field1 > 0 && expr.Field1 <= len(fields) {
		field1Val = fields[expr.Field1-1]
	}

	if expr.Field2 > 0 {
		field2Val := ""
		if expr.Field2 <= len(fields) {
			field2Val = fields[expr.Field2-1]
		}
		return exe.compareValues(field1Val, field2Val, expr.CompareOp)
	}

	return exe.compareValues(field1Val, expr.String2, expr.CompareOp)
}

func (exe *AWKExecutor) compareValues(val1, val2, op string) bool {
	num1, err1 := strconv.ParseFloat(val1, 64)
	num2, err2 := strconv.ParseFloat(val2, 64)

	if err1 == nil && err2 == nil {
		switch op {
		case "==":
			return num1 == num2
		case "!=":
			return num1 != num2
		case "<":
			return num1 < num2
		case ">":
			return num1 > num2
		case "<=":
			return num1 <= num2
		case ">=":
			return num1 >= num2
		}
	}

	switch op {
	case "==":
		return val1 == val2
	case "!=":
		return val1 != val2
	case "<":
		return val1 < val2
	case ">":
		return val1 > val2
	case "<=":
		return val1 <= val2
	case ">=":
		return val1 >= val2
	}

	return false
}

func (exe *AWKExecutor) executeAction(fields []string, line string, writer io.Writer) {
	if exe.program.PrintAll {
		fmt.Fprintln(writer, line)
		return
	}

	if exe.program.PrintRaw != "" && len(exe.program.PrintFields) == 0 {
		output := exe.program.PrintRaw
		output = exe.replaceFields(output, fields, line)
		fmt.Fprintln(writer, output)
		return
	}

	if len(exe.program.PrintFields) > 0 {
		var parts []string
		for _, f := range exe.program.PrintFields {
			if f == 0 {
				parts = append(parts, line)
			} else if f > 0 && f <= len(fields) {
				parts = append(parts, fields[f-1])
			}
		}
		if len(parts) > 0 {
			fmt.Fprintln(writer, strings.Join(parts, " "))
		}
		return
	}

	fmt.Fprintln(writer, line)
}

func (exe *AWKExecutor) replaceFields(s string, fields []string, line string) string {
	result := s

	for i := 1; i <= len(fields); i++ {
		placeholder := fmt.Sprintf("$%d", i)
		result = strings.ReplaceAll(result, placeholder, fields[i-1])
	}
	result = strings.ReplaceAll(result, "$0", line)

	return result
}
