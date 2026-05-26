package main

import (
	"packit/cmd/list"
	"packit/cmd/start"
	"github.com/spf13/cobra"
)

func main() {
	var rootCmd = &cobra.Command{
		Use:   "packit",
		Short: "packit - Windows 轻量命令行抓包工具",
		Long: `packit 是一款轻量、高效的 Windows 平台命令行抓包工具。
无需安装 Wireshark，管理员权限运行即可捕获本机网络流量。

支持功能：
- 自动扫描并列出所有可用网卡
- 实时抓包，彩色终端输出
- 支持 TCP/UDP/ICMP/DNS/HTTP 协议
- 支持按协议、端口、IP 过滤
- 支持保存为 pcap 文件`,
	}

	rootCmd.AddCommand(list.NewListCommand())
	rootCmd.AddCommand(start.NewStartCommand())

	if err := rootCmd.Execute(); err != nil {
		return
	}
}