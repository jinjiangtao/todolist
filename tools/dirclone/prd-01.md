我需要开发一个Go命令行工具，请帮忙生成完整代码。

【工具名称】目录结构克隆器（dirclone）

【功能描述】
克隆源文件夹的目录结构（仅文件夹，不复制文件），
可选保留文件占位符或空文件。

【命令行参数】
-src string      源目录路径（必填）
-dst string      目标目录路径（必填）
-depth int       最大目录深度（0=无限制，默认0）
-empty           创建空文件而非空文件夹（默认false）
-preserve        保留原目录属性（权限/修改时间，默认false）
-ignore string   忽略的目录名，逗号分隔（如.git,node_modules,cache）
-simulate        模拟运行，不实际创建（默认false）
-v               显示详细日志

【使用示例】
# 基础用法：克隆空文件夹结构
dirclone -src /home/project -dst /backup/project

# 只克隆前3层目录
dirclone -src ./app -dst ./empty-app -depth 3

# 创建空文件代替空文件夹（适合某些构建系统）
dirclone -src ./data -dst ./data-empty -empty

# 排除常见目录，模拟运行查看效果
dirclone -src . -dst ../clone -ignore ".git,node_modules,tmp" -simulate

【输出示例】
$ dirclone -src ./src -dst ./dst -v
[模拟运行] 不会实际创建文件/文件夹
[扫描] 源目录: /home/user/project
[深度] 最大深度: 无限制
[忽略] .git, node_modules
[创建] 📁 /dst/config
[创建] 📁 /dst/config/nginx
[创建] 📁 /dst/src/internal
[创建] 📄 /dst/src/main.go (空文件，-empty模式)
[完成] 共创建 47 个目录, 12 个空文件

【技术要求】
1. 纯Go标准库实现（不使用第三方包）
2. 递归扫描用 filepath.WalkDir（Go 1.16+）
3. 并发安全（不需要并发，保持简单）
4. 错误处理：权限不足时跳过并警告，不中断整个流程
5. 跨平台：Windows/Linux/macOS路径兼容（自动处理反斜杠）
6. 性能：百万级目录结构应在合理时间内完成

【代码结构要求】
- 单一main包，可直接编译
- 使用flag包解析参数
- 核心函数：scanAndClone(src, dst string) error
- 辅助函数：shouldIgnore(name string) bool
- 深度控制：传入当前深度参数

【额外要求】
- 提供编译命令注释
- 添加必要的错误处理和用户友好提示
- 输出彩色日志（可选，使用ANSI颜色码）
- 支持中文路径

【输出格式】
素有的输出不能放到一个文件中， 按照功能模块拆开文件， 方便后面扩展。
