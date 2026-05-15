# Windows Grep Tool

一个基于Go语言开发的grep风格命令行工具，专为Windows平台设计，支持中文路径和中文关键词搜索。

## 功能特性

### 基础搜索功能
- 支持关键词搜索文件
- 支持正则表达式匹配
- 递归目录搜索 (`-r`)
- 忽略大小写 (`-i`)
- 显示行号 (`-n`)
- 匹配文本红色高亮输出（Windows终端兼容ANSI颜色）

### 统计功能
- `-c`: 统计每个文件中匹配的行数
- `--total`: 输出全局总匹配行数、匹配文件总数
- `-l`: 仅输出包含匹配内容的文件名
- `-v`: 反向匹配（输出不包含关键词的行）

## 编译命令

```bash
go build -o grep.exe main.go
```

## 使用示例

### 基础搜索
```bash
# 在单个文件中搜索关键词
grep "Hello" test.txt

# 显示行号
grep -n "Hello" test.txt

# 忽略大小写
grep -i "hello" test.txt

# 使用正则表达式
grep "Hello.*World" test.txt
```

### 递归搜索
```bash
# 递归搜索目录
grep -r "Hello" ./

# 递归搜索并显示行号
grep -rn "Hello" ./
```

### 统计功能
```bash
# 统计每个文件中匹配的行数
grep -c "Hello" test.txt

# 递归搜索并统计总数
grep -r --total "Hello" ./

# 仅输出包含匹配内容的文件名
grep -l "Hello" ./

# 反向匹配（输出不包含关键词的行）
grep -v "Hello" test.txt
```

### 组合使用
```bash
# 递归搜索、忽略大小写、显示行号
grep -rin "hello" ./

# 递归搜索并统计总数，同时显示文件名列表
grep -rl --total "Hello" ./
```

## 安装到全局

将 `grep.exe` 复制到系统路径目录（如 `C:\Windows\System32`），或添加工具所在目录到系统环境变量 `PATH` 中。

```bash
# 方法1：复制到System32
copy grep.exe C:\Windows\System32\grep.exe

# 方法2：添加到PATH环境变量
setx PATH "%PATH%;C:\path\to\your\grep\directory"
```

## 技术特点

- 纯Go标准库实现，无需第三方依赖
- 兼容CMD和PowerShell
- 支持中文路径和中文关键词
- 友好的错误提示

## 命令行参数

| 参数 | 说明 |
|------|------|
| `-r` | 递归目录搜索 |
| `-i` | 忽略大小写 |
| `-n` | 显示行号 |
| `-c` | 统计每个文件中匹配的行数 |
| `--total` | 输出全局总匹配行数、匹配文件总数 |
| `-l` | 仅输出包含匹配内容的文件名 |
| `-v` | 反向匹配（输出不包含关键词的行） |