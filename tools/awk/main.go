package main

import (
	"goawk/internal/executor"
	"goawk/internal/parser"
	"os"
	"strings"
)

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	fs := " "
	programStr := ""
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

		if !strings.HasPrefix(arg, "-") && programStr == "" {
			programStr = arg
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

	if programStr == "" {
		printUsage()
		os.Exit(1)
	}

	p := parser.NewParser()
	program, err := p.ParseProgram(programStr)
	if err != nil {
		os.Stderr.WriteString("错误: 解析程序失败 - " + err.Error() + "\n")
		os.Exit(1)
	}

	exe := executor.NewAWKExecutor(fs)
	exe.SetProgram(program)

	if len(inputFiles) == 0 {
		if err := exe.Execute(os.Stdin, os.Stdout); err != nil {
			os.Stderr.WriteString("错误: 执行失败 - " + err.Error() + "\n")
			os.Exit(1)
		}
	} else {
		for _, filePath := range inputFiles {
			file, err := os.Open(filePath)
			if err != nil {
				os.Stderr.WriteString("错误: 无法打开文件 '" + filePath + " - " + err.Error() + "\n")
				continue
			}
			if err := exe.Execute(file, os.Stdout); err != nil {
				os.Stderr.WriteString("错误: 处理文件 '" + filePath + " - " + err.Error() + "\n")
			}
			file.Close()
		}
	}
}

func printUsage() {
	usage := `用法: awk [-F 分隔符] 'pattern { action }' [文件...]

示例:
  awk '$1' test.txt              # 打印第一列
  awk -F ',' '$2' test.csv       # 使用逗号分隔符
  awk '$1 == "error" {print $2}'   # 条件过滤
  echo "a b c" | awk '$1'         # 管道输入
`
	os.Stderr.WriteString(usage)
}
