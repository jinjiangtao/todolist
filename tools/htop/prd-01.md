请使用 Go 语言开发一个 Windows 平台可用的 htop 命令行工具，功能完全模仿 Linux htop，要求如下：

1. 运行环境
- 纯 Go 标准库，无第三方依赖
- 支持 Windows CMD、PowerShell、Windows Terminal
- 支持 ANSI 颜色输出（彩色 CPU/内存条）
- 可编译为 exe 文件：htop.exe

2. 核心功能
- 实时显示系统总览：CPU 使用率、内存使用率、Swap
- 实时进程列表展示：PID、名称、CPU 占用、内存占用
- 支持按 CPU 排序、按内存排序
- 支持上下键滚动进程列表
- 支持刷新频率 1 秒
- 支持中文不乱码
- 界面干净、类似 htop 布局

3. 操作快捷键
- ↑ / ↓ 滚动进程
- q 退出
- F1 帮助
- P 按 CPU 排序
- M 按内存排序

4. 代码要求
- 使用 terminal 绘图、清屏、刷新
- 调用 Windows API 获取进程、CPU、内存数据
- 代码结构清晰、注释完整
- 按照模块组织代码， 要可扩展，不能只写一个main.go 文件
- 每个模块都有对应的测试用例

5. 输出内容
- Windows 编译命令
- 使用方法
- 输出的命令可以在 Windows CMD、PowerShell、Windows Terminal 中运行
- 可以在trae ide 的终端运行