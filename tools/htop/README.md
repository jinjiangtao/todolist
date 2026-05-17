# Windows HTOP - 进程监控工具

## 项目概述

这是一个使用纯 Go 标准库开发的 Windows 平台 htop 风格进程监控工具，功能模仿 Linux htop。

## 功能特性

- 实时显示系统 CPU 使用率、内存使用率、Swap
- 实时进程列表展示：PID、名称、CPU 占用、内存占用
- 支持按 CPU 排序、按内存排序
- 支持上下键滚动进程列表
- 支持刷新频率 1 秒
- 支持 ANSI 颜色输出（彩色 CPU/内存条）
- 支持中文显示
- 界面干净，类似 htop 布局

## 项目结构

```
htop/
├── main.go                    # 主程序入口
├── system/
│   ├── sysinfo.go            # 系统信息获取模块
│   └── sysinfo_test.go       # 系统信息模块测试
├── process/
│   ├── process.go            # 进程管理模块
│   └── process_test.go       # 进程管理模块测试
├── terminal/
│   ├── render.go             # 终端渲染模块
│   └── render_test.go        # 终端渲染模块测试
└── input/
    ├── input.go              # 用户输入处理模块
    └── input_test.go         # 用户输入模块测试
```

## 编译方法

### Windows 编译命令

```powershell
# 初始化 Go 模块（如果还没有）
go mod init htop

# 编译为可执行文件
go build -o htop.exe .
```

### 编译为单个可执行文件

```powershell
# Windows
go build -ldflags="-s -w" -o htop.exe .

# 或者使用 go install
go install
```

## 使用方法

### 基本使用

1. **运行程序**
   ```powershell
   # 直接运行编译好的程序
   .\htop.exe

   # 或者如果已经安装到 PATH
   htop
   ```

2. **操作快捷键**
   - `↑` / `↓` - 上下滚动进程列表
   - `P` - 按 CPU 使用率排序
   - `M` - 按内存使用排序
   - `L` 或 `F1` - 显示帮助信息
   - `Q` - 退出程序

### 界面说明

- **顶部区域**：显示系统总览信息
  - CPU 使用率（带颜色条）
  - 内存使用率（带颜色条）
  - Swap 使用率（带颜色条）
  - CPU 核心数

- **中间区域**：进程列表
  - PID：进程 ID
  - 进程名称：进程可执行文件名
  - CPU：CPU 使用百分比
  - 内存：内存使用量
  - 线程数：进程线程数

- **底部区域**：快捷键提示

### 颜色说明

- 绿色：CPU/内存使用率低 (0-50%)
- 黄色：CPU/内存使用率中 (50-80%)
- 红色：CPU/内存使用率高 (80-100%)

## 运行要求

- Windows 7/8/10/11
- Go 1.16+
- 支持 CMD、PowerShell、Windows Terminal
- 支持 ANSI 颜色输出

## 编译说明

### 依赖

本项目使用纯 Go 标准库，无第三方依赖。需要以下 Windows API：

- kernel32.dll
- psapi.dll

### 编译选项

```powershell
# 标准编译
go build -o htop.exe .

# 优化编译（去除调试信息）
go build -ldflags="-s -w" -o htop.exe .

# 交叉编译（从 Linux 编译 Windows 版本）
GOOS=windows GOARCH=amd64 go build -o htop.exe .

# 编译为 32 位版本
GOARCH=386 go build -o htop.exe .
```

## 测试

运行单元测试：

```powershell
# 测试所有模块
go test ./...

# 测试特定模块
go test ./system/...
go test ./process/...
go test ./terminal/...
go test ./input/...
```

## 常见问题

### Q: 程序无法运行，提示缺少 DLL
A: 确保在 Windows 7 或更高版本上运行，并确保系统目录包含 kernel32.dll 和 psapi.dll（通常系统自带）。

### Q: 颜色不显示
A: 确保使用的是支持 ANSI 颜色的终端（Windows 10 1903+ 或 Windows Terminal）。

### Q: 快捷键没有反应
A: 确保焦点在程序窗口上，并且没有其他程序在捕获键盘输入。

## 扩展开发

### 添加新的系统指标

1. 在 `system/sysinfo.go` 中添加新的数据获取函数
2. 在 `terminal/render.go` 中添加新的显示逻辑
3. 在 `main.go` 中集成新的数据

### 添加新的快捷键

1. 在 `input/input.go` 中定义新的 Key 常量
2. 在 `main.go` 的 `handleKeyPress` 函数中添加新的 case

## 许可证

MIT License

## 作者

使用纯 Go 标准库开发，调用 Windows API 获取系统信息。

## 版本历史

- v1.0.0 - 初始版本，包含基本功能
