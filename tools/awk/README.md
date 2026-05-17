# GoAWK - Windows AWK 命令行工具

使用 Go 语言开发的类 Linux AWK 命令行工具，适用于 Windows 系统，纯标准库实现，无需第三方依赖。

## 功能特性

- ✅ 支持从标准输入/文件读取文本
- ✅ 支持 `-F` 指定分隔符（默认空格）
- ✅ 支持按列输出：`$1` 第1列、`$2` 第2列、`$0` 整行
- ✅ 支持简单表达式过滤：`$1=="xxx"`、`$2>100`
- ✅ 支持打印整行、指定列、多列组合
- ✅ 支持管道命令：`echo "a b c" | goawk '$1'`
- ✅ 兼容 Windows CMD、PowerShell
- ✅ 支持中文、UTF-8 文本
- ✅ 错误提示友好，不崩溃

## 编译

### Windows 编译命令

```bash
go build -o awk.exe main.go
```

### Linux/macOS 编译命令

```bash
go build -o awk main.go
```

## 使用方法

### 基本语法

```bash
goawk [-F 分隔符] 'pattern { action }' [文件...]
```

### 示例数据文件

创建一个测试数据文件 `test.txt`：

```
John 25 Engineer
Alice 30 Manager
Bob 35 Developer
Charlie 28 Designer
```

创建一个 CSV 测试文件 `test.csv`：

```
Name,Age,Position
John,25,Engineer
Alice,30,Manager
Bob,35,Developer
Charlie,28,Designer
```

## 8 个常用使用示例

### 示例 1: 打印第一列

```bash
goawk '$1' test.txt
```

**输出：**
```
John
Alice
Bob
Charlie
```

**说明：** 打印每行的第一个字段（以空格分隔）。

---

### 示例 2: 使用逗号分隔符处理 CSV

```bash
goawk -F ',' '$2' test.csv
```

**输出：**
```
Age
25
30
35
28
```

**说明：** 使用 `-F` 指定逗号作为字段分隔符，打印第二列。

---

### 示例 3: 条件过滤 - 字符串比较

```bash
goawk '$1=="Alice" { print $2, $3 }' test.txt
```

**输出：**
```
30 Manager
```

**说明：** 当第一列等于 "Alice" 时，打印第二和第三列。

---

### 示例 4: 数值比较

```bash
goawk '$2 > 28 { print $1, "年龄:", $2 }' test.txt
```

**输出：**
```
Alice 年龄: 30
Bob 年龄: 35
```

**说明：** 第二列数值大于 28 时，打印相关信息。

---

### 示例 5: 管道输入

```bash
echo "hello world go programming" | goawk '$1'
```

**输出：**
```
hello
```

**说明：** 从标准输入读取数据，打印第一列。

---

### 示例 6: 多列组合打印

```bash
goawk '{ print $1, "-", $3 }' test.txt
```

**输出：**
```
John - Engineer
Alice - Manager
Bob - Developer
Charlie - Designer
```

**说明：** 打印第一列、分隔符、第三列。

---

### 示例 7: 正则表达式匹配

```bash
goawk '$1 ~ /^A/ { print $1 }' test.txt
```

**输出：**
```
Alice
```

**说明：** 打印第一列以 "A" 开头的行。

---

### 示例 8: 跳过表头处理数据

```bash
goawk 'NR>1 { print $1, $2 }' test.csv
```

**输出：**
```
John 25
Alice 30
Bob 35
Charlie 28
```

**说明：** 跳过第一行（表头），打印剩余行的第一和第二列。

---

## 高级用法

### 组合多个条件（AND）

```bash
goawk '$2 > 25 && $2 < 35 { print $1 }' test.txt
```

**输出：**
```
Alice
Charlie
```

### 组合多个条件（OR）

```bash
goawk '$1 == "John" || $1 == "Bob" { print $1 }' test.txt
```

**输出：**
```
John
Bob
```

### 反向匹配（NOT）

```bash
goawk '!($2 > 30) { print $1 }' test.txt
```

**输出：**
```
John
Alice
Charlie
```

### 打印整行（$0）

```bash
goawk '$2 > 30 { print $0 }' test.txt
```

**输出：**
```
Bob 35 Developer
```

### 统计计算

```bash
goawk '{ sum += $2 } END { print "平均年龄:", sum/(NR-1) }' test.txt
```

**说明：** 计算第二列的总和和平均值。

## 在不同环境中的使用

### CMD (Windows 命令提示符)

```cmd
echo hello world | awk.exe "$1"
```

### PowerShell

```powershell
echo "hello world" | .\awk.exe '$1'
Get-Content test.txt | .\awk.exe '$1'
```

### Git Bash / MSYS2

```bash
echo "hello world" | ./awk.exe '$1'
```

## 与标准 AWK 的兼容性

本工具实现了 AWK 的核心功能，包括：

- ✅ 字段引用（`$0`, `$1`, `$2`, ...）
- ✅ 字段分隔符（`-F` 或 `BEGIN { FS="..." }`）
- ✅ 关系运算符（`==`, `!=`, `<`, `>`, `<=`, `>=`）
- ✅ 逻辑运算符（`&&`, `||`, `!`）
- ✅ 正则表达式匹配（`~`, `!~`）
- ✅ 内置变量概念（`NR`, `NF` - 部分支持）

## 错误处理

程序包含完善的错误处理机制：

- 文件读取错误：友好的错误提示，不会崩溃
- 语法错误：清晰的错误信息，指出问题所在
- 字段越界：自动处理，不会panic

## 性能特点

- 纯 Go 标准库实现，无外部依赖
- 高效的字段分割算法
- 支持大文件处理
- 内存占用低

## 测试

运行内置测试：

```bash
go run main.go
```

测试将验证：
- 词法分析器（Lexer）
- 字段分割功能
- 模式匹配
- 比较操作
- CSV 处理
- 多字段打印
- 中文文本支持

## 许可证

MIT License

## 贡献

欢迎提交 Issue 和 Pull Request！
