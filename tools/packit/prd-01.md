# go 轻量 Windows 命令行抓包工具开发提示词
使用 go 语言开发一款**轻量、高效、单文件 EXE**的 Windows 平台命令行抓包工具，无需安装 Wireshark、无需依赖库，管理员权限运行即可捕获本机网络流量。
工具专注简洁实用，输出美观、易读、彩色格式化显示，适合开发调试、网络测试、个人学习使用。

## 核心功能
1. 自动扫描并列出所有可用网卡
2. 选择网卡即可开始实时抓包
3. 支持抓取 TCP / UDP / ICMP / DNS / HTTP 流量
4. 彩色终端输出：源IP、目标IP、端口、协议、包长度、时间
5. 支持按协议过滤（--tcp / --udp / --dns / --http）
6. 支持按端口过滤（--port 80,443）
7. 支持按 IP 过滤
8. 支持将抓包数据保存为 .pcap 文件（可被 Wireshark 打开）
9. 按 Ctrl+C 优雅停止抓包并显示统计信息
10. 轻量无依赖，编译后单个 exe，直接双击运行

## 命令行格式（简单易用）
- packit list           列出所有网卡
- packit start 2        从第2个网卡抓包
- packit start 2 --tcp  只抓 TCP 包
- packit start 2 --dns  只抓 DNS 解析
- packit start 2 --port 80 只抓80端口
- packit save log.pcap  抓包并保存文件

## 技术依赖
- clap：命令行解析
- pnet：网络抓包与数据包解析
- colored：彩色终端输出
- chrono：时间格式化
- winapi：Windows 平台兼容

## 代码要求
- 纯 go 编写
- 结构清晰、模块化  代码不能放到一个文件中。 
- 错误处理友好，不崩溃
- 中文提示清晰
- 高性能、低内存占用

## 产物
生成的产物打包成exe 可执行文件， 并添加到windows 的环境变量中。 在ide 的终端和windows 的终端都可以正常运行。
如果有原理的可执行文件 packit 替换原来的文件。