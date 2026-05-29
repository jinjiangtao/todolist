# Smart Renamer - 智能文件批量重命名工具

## 项目介绍

Smart Renamer 是一款使用 Go + Wails + Vue3 开发的智能文件批量重命名工具。

## 功能特性

1. **5种重命名规则：
   - 序号填充：自定义起始数字、位数、前缀后缀
   - 查找替换：支持普通字符串和正则表达式
   - 日期标记：从文件修改时间提取日期作为前缀/后缀
   - 大小写转换：全大写/全小写/首字母大写/驼峰式
   - 扩展名修改：批量更改扩展名（含大小写统一）

2. 实时预览：显示原文件名→新文件名对照表
3. 支持多规则组合应用（按顺序执行）
4. 一键撤销：内置撤销功能
5. 冲突处理：遇到同名文件时可选择覆盖/跳过/自动编号
6. 支持拖拽文件夹

## 技术栈

- **后端**: Go 1.21+
- **GUI框架**: Wails v2
- **前端**: Vue 3 + Vite

## 项目结构

```
smart-renamer/
├── main.go                 # 主程序入口
├── app.go                  # 应用主逻辑
├── go.mod                  # Go 模块配置
├── wails.json             # Wails 配置
├── internal/
│   ├── models/            # 数据模型
│   │   └── models.go
│   ├── scanner/           # 文件扫描器
│   │   └── scanner.go
│   ├── rules/             # 规则引擎
│   │   └── engine.go
│   ├── previewer/         # 预览管理器
│   │   └── previewer.go
│   └── renamer/           # 重命名执行器
│       └── renamer.go
└── frontend/              # Vue 前端
    ├── index.html
    ├── package.json
    ├── vite.config.js
    └── src/
        ├── main.js
        ├── App.vue
        ├── style.css
        └── components/
            ├── DirectorySelector.vue
            ├── RuleManager.vue
            ├── PreviewTable.vue
            └── ActionPanel.vue
```

## 开发环境准备

### 1. 安装 Go

下载并安装 Go 1.21 或更高版本：https://go.dev/dl/

### 2. 安装 Wails

```bash
go install github.com/wailsapp/wails/v2/cmd/wails@latest
```

### 3. 安装 Node.js

下载并安装 Node.js 16 或更高版本：https://nodejs.org/

## 开发模式运行

1. 进入项目目录，先安装前端依赖：

```bash
cd frontend
npm install
cd ..
```

2. 在开发模式下运行（热重载）：

```bash
wails dev
```

## 编译打包

### Windows

```bash
wails build
```

编译后会在 `build/bin/` 目录下生成单个 exe 文件。

### 编译选项

如果需要 UPX 压缩（减小体积），先安装 UPX，然后：

```bash
wails build -upx
```

## 使用说明

1. 启动应用
2. 拖放文件夹到输入框，或手动输入目录路径
3. 在左侧勾选需要使用的重命名规则
4. 配置规则参数
5. 在右侧预览重命名结果
6. 确认无误后点击"开始重命名"
7. 如需撤销，点击"撤销上次操作"

## 核心模块说明

### models
定义了所有数据结构，包括 File, Rule, PreviewItem 等。

### scanner
使用并发扫描目录，返回所有文件信息。

### rules
实现5种重命名规则的具体实现，支持链式调用。

### previewer
生成预览结果并检测冲突。

### renamer
执行重命名操作，并保存撤销历史记录。

## 许可证

MIT License
