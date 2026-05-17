请使用 Go 语言开发一个类 Linux AWK 的命令行工具，适用于 Windows 系统，纯标准库实现，无需第三方依赖，可编译成 awk.exe。

功能要求：
1. 支持从标准输入/文件读取文本
2. 支持 -F 指定分隔符（默认空格）
3. 支持按列输出：$1 第1列、$2 第2列、$0 整行
4. 支持简单表达式过滤：$1=="xxx"、$2>100
5. 支持打印整行、指定列、多列组合
6. 支持管道命令：echo "a b c" | goawk '$1'
7. 兼容 Windows CMD、PowerShell
8. 支持中文、UTF-8 文本
9. 错误提示友好，不崩溃

命令格式示例：
goawk [-F 分隔符] 'pattern { action }' 文件
goawk '$1' test.txt
goawk -F "," '$2' test.csv
goawk '$1=="error" {print $2}' test.log

输出要求：
1. 完整单文件 main.go 代码， 包含所有功能实现， 代码也分模块组织
2. Windows 编译命令：go build -o awk.exe main.go
3. 提供 8 个常用使用示例
4. 代码结构清晰、注释完整
5. 输出放到 说明文档到README.md
6. go 封装的模块方法提供test 函数， 并运行test 函数。
7. 生成awk 命令自动注册到环境变量， 可直接在命令行和ide 的shell中使用。
