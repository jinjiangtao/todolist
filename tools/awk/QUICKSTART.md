# GoAWK 快速入门指南

## 编译

```bash
go build -o awk.exe main.go
```

## 基本用法

### 1. 打印第一列
```bash
awk '$1' test.txt
```

### 2. 使用自定义分隔符
```bash
awk -F ',' '$2' test.csv
```

### 3. 条件过滤
```bash
awk '$1 == "Alice" { print $2 }' test.txt
awk '$2 > 28 { print $1 }' test.txt
```

### 4. 管道输入
```bash
echo "hello world" | awk '$1'
```

### 5. 多列组合
```bash
awk '{ print $1, "-", $3 }' test.txt
```

### 6. 正则匹配
```bash
awk '$1 ~ /^A/ { print }' test.txt
```

### 7. 组合条件
```bash
awk '$2 > 25 && $2 < 35 { print $1 }' test.txt
```

### 8. 跳过表头
```bash
awk 'NR>1 { print $1 }' data.csv
```

## 测试文件

项目包含以下测试文件：
- `test.txt` - 空格分隔的测试数据
- `test.csv` - CSV 格式测试数据
- `test_cn.txt` - 中文测试数据

## 运行测试

```bash
# 完整功能测试
awk '$1' test.txt
awk '$2 > 28 { print $1 }' test.txt
awk -F ',' '$2' test.csv
```

## 环境变量

已添加到用户 PATH：
```
D:\trae\todolist\tools\awk
```

新打开的终端可以直接使用 `awk` 命令。

## PowerShell 使用

```powershell
# 打印第一列
awk '$1' test.txt

# 管道输入
echo "a b c" | awk '$1'

# CSV 处理
awk -F ',' '$2' test.csv
```

## CMD 使用

```cmd
awk "$1" test.txt
echo a b c | awk "$1"
```

## Git Bash / MSYS2 使用

```bash
./awk.exe '$1' test.txt
echo "a b c" | ./awk.exe '$1'
```

## 错误处理

程序包含完善的错误处理：
- 文件不存在：友好提示
- 语法错误：清晰的错误信息
- 字段越界：自动处理

## 性能特点

- 纯 Go 标准库实现
- 高效的字段分割
- 支持大文件
- 低内存占用

## 故障排除

### 命令未找到
关闭并重新打开终端，或手动添加到 PATH：
```powershell
$env:Path += ";D:\trae\todolist\tools\awk"
```

### 中文乱码
确保文件保存为 UTF-8 编码

### 编译错误
```bash
go version  # 确保 Go 已安装
go build -o awk.exe main.go
```

## 更多信息

详见 [README.md](README.md)
