# GoAWK - Trae IDE 项目配置

## 项目概述
Go 语言实现的类 Linux AWK 命令行工具，支持 Windows 系统，纯标准库实现。

## 项目结构
```
goawk/
├── cmd/
│   └── awk/
│       └── main.go
├── internal/
│   ├── lexer/
│   │   ├── lexer.go
│   │   └── lexer_test.go
│   ├── parser/
│   │   ├── parser.go
│   │   └── parser_test.go
│   └── executor/
│       ├── executor.go
│       └── executor_test.go
├── pkg/
├── test/
├── go.mod
├── README.md
├── QUICKSTART.md
├── awk.exe (编译后)
└── test.txt (测试数据)
```

## 终端配置
在 Trae IDE 终端中直接使用 `awk` 命令：

### 方法 1：使用完整路径（推荐用于临时测试）
```powershell
# 在项目根目录中
.\awk.exe '$1' test.txt
```

### 方法 2：临时添加到当前会话的 PATH
```powershell
# 在项目根目录中运行
$env:PATH += ";$(Get-Location)"
```

### 方法 3：永久添加到用户环境变量（已设置）
- 路径：`D:\trae\todolist\tools\awk`
- 重新打开终端即可生效

## 编译命令
```powershell
go build -o awk.exe ./cmd/awk
```

## 测试命令
```powershell
# 运行所有模块测试
go test ./internal/lexer
go test ./internal/parser
go test ./internal/executor

# 或者运行所有测试
go test ./...
```

## 功能特性
- 支持从标准输入/文件读取文本
- 支持 `-F` 指定分隔符（默认空格）
- 支持按列输出：`$1` 第 1 列、`$2` 第 2 列、`$0` 整行
- 支持简单表达式过滤：`$1=="xxx"`、`$2>100`
- 支持管道命令
- 兼容 Windows CMD、PowerShell
- 支持中文、UTF-8 文本

## 使用示例
```powershell
# 打印第一列
awk '$1' test.txt

# 使用逗号分隔符
awk -F ',' '$2' test.csv

# 条件过滤
awk '$2 > 25 { print $1 }' test.txt

# 管道输入
echo "hello world" | awk '$1'
```
