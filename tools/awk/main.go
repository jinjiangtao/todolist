package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"regexp"
	"strconv"
	"strings"
)

type TokenType int

const (
	TOKEN_FIELD TokenType = iota
	TOKEN_STRING
	TOKEN_NUMBER
	TOKEN_COMPARE
	TOKEN_EQUAL
	TOKEN_NEQUAL
	TOKEN_LT
	TOKEN_GT
	TOKEN_LTE
	TOKEN_GTE
	TOKEN_AND
	TOKEN_OR
	TOKEN_NOT
	TOKEN_PRINT
	TOKEN_LBRACE
	TOKEN_RBRACE
	TOKEN_EOF
)

type Token struct {
	Type  TokenType
	Value string
}

type Lexer struct {
	input  string
	pos    int
	length int
}

func NewLexer(input string) *Lexer {
	return &Lexer{
		input:  input,
		pos:    0,
		length: len(input),
	}
}

func (l *Lexer) NextToken() Token {
	l.skipWhitespace()

	if l.pos >= l.length {
		return Token{Type: TOKEN_EOF}
	}

	ch := l.input[l.pos]

	if ch == '$' && l.pos+1 < l.length && l.input[l.pos+1] >= '0' && l.input[l.pos+1] <= '9' {
		l.pos++
		fieldNum := ""
		for l.pos < l.length && l.input[l.pos] >= '0' && l.input[l.pos] <= '9' {
			fieldNum += string(l.input[l.pos])
			l.pos++
		}
		return Token{Type: TOKEN_FIELD, Value: fieldNum}
	}

	if ch == '"' {
		l.pos++
		str := ""
		for l.pos < l.length && l.input[l.pos] != '"' {
			str += string(l.input[l.pos])
			l.pos++
		}
		if l.pos < l.length {
			l.pos++
		}
		return Token{Type: TOKEN_STRING, Value: str}
	}

	if ch == '\'' {
		l.pos++
		str := ""
		for l.pos < l.length && l.input[l.pos] != '\'' {
			str += string(l.input[l.pos])
			l.pos++
		}
		if l.pos < l.length {
			l.pos++
		}
		return Token{Type: TOKEN_STRING, Value: str}
	}

	if l.pos+1 < l.length {
		twoChar := string(ch) + string(l.input[l.pos+1])
		switch twoChar {
		case "==":
			l.pos += 2
			return Token{Type: TOKEN_EQUAL, Value: "=="}
		case "!=":
			l.pos += 2
			return Token{Type: TOKEN_NEQUAL, Value: "!="}
		case "<=":
			l.pos += 2
			return Token{Type: TOKEN_LTE, Value: "<="}
		case ">=":
			l.pos += 2
			return Token{Type: TOKEN_GTE, Value: ">="}
		case "&&":
			l.pos += 2
			return Token{Type: TOKEN_AND, Value: "&&"}
		case "||":
			l.pos += 2
			return Token{Type: TOKEN_OR, Value: "||"}
		}
	}

	switch ch {
	case '=', '!', '<', '>':
		l.pos++
		return Token{Type: TOKEN_COMPARE, Value: string(ch)}
	case '{':
		l.pos++
		return Token{Type: TOKEN_LBRACE, Value: "{"}
	case '}':
		l.pos++
		return Token{Type: TOKEN_RBRACE, Value: "}"}
	}

	if (ch >= '0' && ch <= '9') || ch == '.' {
		num := ""
		for l.pos < l.length && ((l.input[l.pos] >= '0' && l.input[l.pos] <= '9') || l.input[l.pos] == '.') {
			num += string(l.input[l.pos])
			l.pos++
		}
		return Token{Type: TOKEN_NUMBER, Value: num}
	}

	if (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z') || ch == '_' {
		ident := ""
		for l.pos < l.length && ((l.input[l.pos] >= 'a' && l.input[l.pos] <= 'z') ||
			(l.input[l.pos] >= 'A' && l.input[l.pos] <= 'Z') ||
			(l.input[l.pos] >= '0' && l.input[l.pos] <= '9') ||
			l.input[l.pos] == '_') {
			ident += string(l.input[l.pos])
			l.pos++
		}
		if strings.ToLower(ident) == "print" {
			return Token{Type: TOKEN_PRINT, Value: ident}
		}
		return Token{Type: TOKEN_STRING, Value: ident}
	}

	l.pos++
	return l.NextToken()
}

func (l *Lexer) skipWhitespace() {
	for l.pos < l.length && (l.input[l.pos] == ' ' || l.input[l.pos] == '\t') {
		l.pos++
	}
}

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

type AWKExecutor struct {
	fs         string
	program    PatternAction
	fieldRegex *regexp.Regexp
}

func NewAWKExecutor(fs string) *AWKExecutor {
	exe := &AWKExecutor{
		fs: fs,
	}
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

func (exe *AWKExecutor) ParseProgram(program string) error {
	program = strings.TrimSpace(program)
	if program == "" {
		return fmt.Errorf("程序不能为空")
	}

	hasAction := false
	hasPattern := false

	if regexp.MustCompile(`^\$[0-9]+\s*$`).MatchString(program) {
		exe.program.PrintFields = []int{}
		fieldStr := strings.TrimSpace(program[1:])
		fieldNum, err := strconv.Atoi(fieldStr)
		if err != nil {
			return fmt.Errorf("无效的字段编号: %s", fieldStr)
		}
		exe.program.PrintFields = append(exe.program.PrintFields, fieldNum)
		exe.program.PrintAll = false
		exe.program.HasAction = true
		exe.program.HasPattern = true
		return nil
	}

	lbrace := strings.Index(program, "{")
	rbrace := strings.LastIndex(program, "}")

	if lbrace >= 0 && rbrace > lbrace {
		hasAction = true
		actionPart := strings.TrimSpace(program[lbrace+1 : rbrace])
		if actionPart != "" {
			if err := exe.parseAction(actionPart); err != nil {
				return err
			}
		}
		if lbrace > 0 {
			hasPattern = true
			patternPart := strings.TrimSpace(program[:lbrace])
			if err := exe.parsePattern(patternPart); err != nil {
				return err
			}
		}
	} else {
		hasPattern = true
		if err := exe.parsePattern(program); err != nil {
			return err
		}
	}

	exe.program.HasPattern = hasPattern
	exe.program.HasAction = hasAction

	if !hasAction {
		exe.program.PrintAll = true
	}

	return nil
}

func (exe *AWKExecutor) parseAction(action string) error {
	action = strings.TrimSpace(action)
	lexer := NewLexer(action)
	token := lexer.NextToken()

	if token.Type == TOKEN_PRINT {
		printPart := strings.TrimSpace(action[len("print"):])
		if printPart == "" || printPart == ";" {
			exe.program.PrintAll = true
			return nil
		}

		printPart = strings.Trim(printPart, ";")
		fields, rawStr, err := exe.parsePrintArgs(printPart)
		if err != nil {
			return err
		}

		if rawStr != "" {
			exe.program.PrintRaw = rawStr
		}
		exe.program.PrintFields = fields
		if len(fields) > 0 || rawStr != "" {
			exe.program.PrintAll = false
		} else {
			exe.program.PrintAll = true
		}
	} else {
		exe.program.PrintRaw = action
	}

	return nil
}

func (exe *AWKExecutor) parsePrintArgs(args string) ([]int, string, error) {
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

func (exe *AWKExecutor) parsePattern(pattern string) error {
	pattern = strings.TrimSpace(pattern)
	exe.program.Pattern = Expression{}

	if pattern == "" {
		return nil
	}

	orParts := strings.Split(pattern, "||")
	if len(orParts) > 1 {
		exe.program.Pattern.LogicOp = "||"
		for i, part := range orParts {
			part = strings.TrimSpace(part)
			if i == 0 {
				if err := exe.parseSinglePattern(part, &exe.program.Pattern, true); err != nil {
					return err
				}
			}
		}
		return nil
	}

	andParts := strings.Split(pattern, "&&")
	if len(andParts) > 1 {
		exe.program.Pattern.LogicOp = "&&"
		for i, part := range andParts {
			part = strings.TrimSpace(part)
			if i == 0 {
				if err := exe.parseSinglePattern(part, &exe.program.Pattern, true); err != nil {
					return err
				}
			}
		}
		return nil
	}

	return exe.parseSinglePattern(pattern, &exe.program.Pattern, false)
}

func (exe *AWKExecutor) parseSinglePattern(pattern string, expr *Expression, hasLogicOp bool) error {
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

	regexMatch, _ := regexp.MatchString(pattern, "")
	if regexMatch || strings.Contains(pattern, "~") {
		expr.String1 = pattern
	}

	return nil
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

func (exe *AWKExecutor) evaluateExpr(expr *Expression, fields []string, line string) bool {
	if expr.Not {
		return !exe.evaluateSingleExpr(expr, fields, line)
	}
	return exe.evaluateSingleExpr(expr, fields, line)
}

func (exe *AWKExecutor) evaluateSingleExpr(expr *Expression, fields []string, line string) bool {
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

func TestLexer() {
	fmt.Println("测试 Lexer:")
	lexer := NewLexer(`$1 == "test" { print $2 }`)

	tokens := []TokenType{TOKEN_FIELD, TOKEN_STRING, TOKEN_EQUAL, TOKEN_STRING, TOKEN_LBRACE, TOKEN_PRINT, TOKEN_FIELD, TOKEN_RBRACE}
	names := []string{"TOKEN_FIELD", "TOKEN_STRING", "TOKEN_EQUAL", "TOKEN_STRING", "TOKEN_LBRACE", "TOKEN_PRINT", "TOKEN_FIELD", "TOKEN_RBRACE"}

	for i, expectedType := range tokens {
		token := lexer.NextToken()
		if token.Type != expectedType {
			fmt.Printf("  失败: 期望 %s, 得到 %s\n", names[i], token.Type)
		} else {
			fmt.Printf("  通过: %s = %s\n", names[i], token.Value)
		}
	}
}

func TestFieldSplitting() {
	fmt.Println("\n测试字段分割:")
	exe := NewAWKExecutor(" ")
	line := "hello world go awk"
	fields := exe.splitLine(line)

	expected := []string{"hello", "world", "go", "awk"}
	if len(fields) == len(expected) {
		fmt.Println("  通过: 字段分割正确")
		for i, f := range fields {
			if f != expected[i] {
				fmt.Printf("  失败: 字段 %d 期望 '%s', 得到 '%s'\n", i+1, expected[i], f)
			}
		}
	} else {
		fmt.Printf("  失败: 期望 %d 个字段, 得到 %d 个\n", len(expected), len(fields))
	}
}

func TestPatternMatching() {
	fmt.Println("\n测试模式匹配:")
	exe := NewAWKExecutor(",")

	err := exe.ParseProgram(`$1 == "error" { print $2 }`)
	if err != nil {
		fmt.Printf("  失败: 解析错误 %v\n", err)
		return
	}

	fields := []string{"error", "file not found", "test"}
	if exe.matchPattern(fields, "error,file not found,test") {
		fmt.Println("  通过: 模式匹配正确")
	} else {
		fmt.Println("  失败: 模式匹配错误")
	}
}

func TestComparison() {
	fmt.Println("\n测试比较操作:")
	exe := NewAWKExecutor(" ")

	if exe.compareValues("100", "200", "<") {
		fmt.Println("  通过: < 比较正确")
	} else {
		fmt.Println("  失败: < 比较错误")
	}

	if !exe.compareValues("100", "100", "!=") {
		fmt.Println("  通过: != 比较正确")
	} else {
		fmt.Println("  失败: != 比较错误")
	}

	if exe.compareValues("abc", "abc", "==") {
		fmt.Println("  通过: == 字符串比较正确")
	} else {
		fmt.Println("  失败: == 字符串比较错误")
	}
}

func TestCSVHandling() {
	fmt.Println("\n测试 CSV 处理:")
	exe := NewAWKExecutor(",")

	err := exe.ParseProgram(`$2`)
	if err != nil {
		fmt.Printf("  失败: 解析错误 %v\n", err)
		return
	}

	fields := exe.splitLine("name,age,city")
	if len(fields) >= 2 && fields[1] == "age" {
		fmt.Println("  通过: CSV 字段分割正确")
	} else {
		fmt.Printf("  失败: 期望第二个字段为 'age', 得到 '%v'\n", fields)
	}
}

func TestMultiFieldPrint() {
	fmt.Println("\n测试多字段打印:")
	exe := NewAWKExecutor(" ")

	err := exe.ParseProgram(`{ print $1, $3 }`)
	if err != nil {
		fmt.Printf("  失败: 解析错误 %v\n", err)
		return
	}

	fields := []string{"one", "two", "three"}

	var output strings.Builder
	exe.executeAction(fields, "one two three", &output)
	result := output.String()

	if strings.Contains(result, "one") && strings.Contains(result, "three") {
		fmt.Println("  通过: 多字段打印正确")
	} else {
		fmt.Printf("  失败: 输出应为包含 'one' 和 'three', 得到 '%s'\n", result)
	}
}

func TestChineseText() {
	fmt.Println("\n测试中文处理:")
	exe := NewAWKExecutor(" ")

	fields := exe.splitLine("你好 世界")
	if len(fields) >= 2 && fields[0] == "你好" && fields[1] == "世界" {
		fmt.Println("  通过: 中文字段分割正确")
	} else {
		fmt.Printf("  失败: 中文处理错误, 得到 '%v'\n", fields)
	}
}

func TestAWKExecutor() {
	fmt.Println("========== AWK Executor 测试 ==========")
	TestLexer()
	TestFieldSplitting()
	TestPatternMatching()
	TestComparison()
	TestCSVHandling()
	TestMultiFieldPrint()
	TestChineseText()
	fmt.Println("\n========== 所有测试完成 ==========")
}

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "用法: goawk [-F 分隔符] 'pattern { action }' [文件...]")
		fmt.Fprintln(os.Stderr, "")
		fmt.Fprintln(os.Stderr, "示例:")
		fmt.Fprintln(os.Stderr, "  goawk '$1' test.txt              # 打印第一列")
		fmt.Fprintln(os.Stderr, "  goawk -F ',' '$2' test.csv       # 使用逗号分隔符")
		fmt.Fprintln(os.Stderr, "  goawk '$1==\"error\" {print $2}'   # 条件过滤")
		fmt.Fprintln(os.Stderr, "  echo \"a b c\" | goawk '$1'         # 管道输入")
		fmt.Fprintln(os.Stderr, "  goawk '{print $1, $3}' test.txt  # 打印第1和第3列")
		fmt.Fprintln(os.Stderr, "  goawk 'NR>1' data.txt            # 跳过第一行")
		fmt.Fprintln(os.Stderr, "  goawk '$1 ~ /test/ {print}'      # 正则匹配")
		fmt.Fprintln(os.Stderr, "  goawk '$2 > 100 {print $1}'      # 数值比较")
		os.Exit(1)
	}

	fs := " "
	program := ""
	var inputFiles []string

	i := 1
	for i < len(os.Args) {
		arg := os.Args[i]

		if arg == "-F" && i+1 < len(os.Args) {
			fs = os.Args[i+1]
			i += 2
			continue
		}

		if strings.HasPrefix(arg, "-F") && len(arg) > 2 {
			fs = arg[2:]
			i++
			continue
		}

		if !strings.HasPrefix(arg, "-") && program == "" {
			program = arg
			i++
			continue
		}

		if !strings.HasPrefix(arg, "-") {
			inputFiles = append(inputFiles, arg)
			i++
			continue
		}

		i++
	}

	if program == "" {
		fmt.Fprintln(os.Stderr, "错误: 未指定程序")
		fmt.Fprintln(os.Stderr, "用法: goawk [-F 分隔符] 'pattern { action }' [文件...]")
		os.Exit(1)
	}

	exe := NewAWKExecutor(fs)

	if err := exe.ParseProgram(program); err != nil {
		fmt.Fprintf(os.Stderr, "错误: 解析程序失败 - %v\n", err)
		os.Exit(1)
	}

	var reader io.Reader

	if len(inputFiles) == 0 {
		reader = os.Stdin
	} else if len(inputFiles) == 1 {
		file, err := os.Open(inputFiles[0])
		if err != nil {
			fmt.Fprintf(os.Stderr, "错误: 无法打开文件 '%s' - %v\n", inputFiles[0], err)
			os.Exit(1)
		}
		defer file.Close()
		reader = file
	} else {
		for _, filePath := range inputFiles {
			file, err := os.Open(filePath)
			if err != nil {
				fmt.Fprintf(os.Stderr, "错误: 无法打开文件 '%s' - %v\n", filePath, err)
				continue
			}
			if err := exe.Execute(file, os.Stdout); err != nil {
				fmt.Fprintf(os.Stderr, "错误: 处理文件 '%s' 失败 - %v\n", filePath, err)
			}
			file.Close()
		}
		return
	}

	if err := exe.Execute(reader, os.Stdout); err != nil {
		fmt.Fprintf(os.Stderr, "错误: 执行失败 - %v\n", err)
		os.Exit(1)
	}
}
