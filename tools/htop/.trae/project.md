# Windows HTOP - Trae IDE 项目配置

## 项目概述
Go 语言实现的类 Linux htop 进程监控工具，支持 Windows 系统，纯标准库实现。

## 项目结构
```
htop/
├── system/
│   ├── sysinfo.go
│   └── sysinfo_test.go
├── process/
│   ├── process.go
│   └── process_test.go
├── terminal/
│   ├── render.go
│   └── render_test.go
├── input/
│   ├── input.go
│   └── input_test.go
├── .trae/
│   ├── init-terminal.ps1
│   └── project.md
├── go.mod
├── README.md
├── INSTALL_GUIDE.md
├── QUICKSTART.md
└── htop.exe (编译后)
```

## 终端配置
在 Trae IDE 终端中直接使用 `htop` 命令：

### 方法 1：使用完整路径（推荐用于临时测试）
```powershell
# 在项目根目录中
.\htop.exe
```

### 方法 2：临时添加到当前会话的 PATH
```powershell
# 在项目根目录中运行
$env:PATH += ";$(Get-Location)"
```

### 方法 3：永久添加到用户环境变量（已设置）
- 路径：`D:\trae\todolist\tools\htop`
- 重新打开终端即可生效

### 方法 4：Trae IDE 终端自动配置（推荐）
- Trae IDE 会自动加载 `.trae/init-terminal.ps1`
- 打开终端时自动配置 PATH
- 可以直接使用 `htop` 命令

## 编译命令
```powershell
go build -o htop.exe .
```

## 测试命令
```powershell
# 运行所有模块测试
go test ./system
go test ./process
go test ./terminal
go test ./input

# 或者运行所有测试
go test ./...
```

## 功能特性
- 实时显示 CPU、内存、Swap 使用率（带彩色进度条）
- 实时进程列表展示：PID、名称、CPU、内存、线程数
- 支持按 CPU 排序、按内存排序
- 支持上下键滚动进程列表
- 支持 1 秒刷新频率
- 支持 ANSI 颜色输出
- 支持中文显示

## 快捷键
```
↑ / ↓    - 滚动进程列表
P        - 按 CPU 使用率排序
M        - 按内存使用排序
L / F1   - 显示帮助信息
Q        - 退出程序
```

## 使用示例
```powershell
# 运行进程监控
htop

# 按 CPU 排序
# 在 htop 中按 P 键

# 按内存排序
# 在 htop 中按 M 键

# 退出
# 在 htop 中按 Q 键
```

## 运行要求
- Windows 7/8/10/11
- Go 1.16+
- 支持 CMD、PowerShell、Windows Terminal
- 支持 ANSI 颜色输出
