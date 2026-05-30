# 目录结构克隆器 (dirclone)

## 编译方法

### Windows
```bash
go build -o dirclone.exe .
```

### Linux/macOS
```bash
go build -o dirclone .
```

### 交叉编译
```bash
# 编译为Linux二进制
GOOS=linux GOARCH=amd64 go build -o dirclone-linux .

# 编译为macOS二进制
GOOS=darwin GOARCH=amd64 go build -o dirclone-macos .

# 编译为Windows二进制(在Linux/macOS上)
GOOS=windows GOARCH=amd64 go build -o dirclone.exe .
```

## 使用方法

### 基础用法
```bash
# 克隆空文件夹结构
dirclone -src /home/project -dst /backup/project

# 只克隆前3层目录
dirclone -src ./app -dst ./empty-app -depth 3

# 创建空文件代替空文件夹
dirclone -src ./data -dst ./data-empty -empty

# 排除常见目录，模拟运行查看效果
dirclone -src . -dst ../clone -ignore ".git,node_modules,tmp" -simulate

# 保留原目录属性
dirclone -src ./src -dst ./dst -preserve -v
```

### 命令行参数

| 参数 | 描述 | 默认值 |
|------|------|--------|
| `-src` | 源目录路径 (必填) | - |
| `-dst` | 目标目录路径 (必填) | - |
| `-depth` | 最大目录深度 (0=无限制) | 0 |
| `-empty` | 创建空文件而非空文件夹 | false |
| `-preserve` | 保留原目录属性 (权限/修改时间) | false |
| `-ignore` | 忽略的目录名，逗号分隔 | "" |
| `-simulate` | 模拟运行，不实际创建 | false |
| `-v` | 显示详细日志 | false |

## 功能特性

- ✅ 纯Go标准库实现（无第三方依赖）
- ✅ 使用 `filepath.WalkDir` 递归扫描（Go 1.16+）
- ✅ 权限不足时跳过并警告，不中断整个流程
- ✅ 跨平台支持（Windows/Linux/macOS）
- ✅ 支持中文路径
- ✅ 彩色输出日志
- ✅ 模拟运行模式
- ✅ 深度限制
- ✅ 忽略指定目录
- ✅ 创建空文件代替空文件夹
- ✅ 保留原目录属性

## 项目结构

```
dirclone/
├── main.go              # 程序入口
├── config/
│   └── config.go       # 配置结构体
├── logger/
│   └── logger.go       # 彩色日志输出
├── clone/
│   ├── ignore.go       # 忽略规则匹配
│   ├── scanner.go      # 目录扫描
│   └── cloner.go       # 克隆执行
├── cli/
│   └── flags.go        # 命令行参数解析
└── go.mod              # Go模块定义
```

## 许可证

MIT License
