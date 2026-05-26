package start

import (
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"

	"github.com/fatih/color"
	"packit/pkg/capture"
	"packit/pkg/iface"
	"packit/pkg/packet"
	"github.com/spf13/cobra"
)

func NewStartCommand() *cobra.Command {
	var tcpFlag bool
	var udpFlag bool
	var dnsFlag bool
	var httpFlag bool
	var portFlag string
	var ipFlag string
	var saveFlag string

	cmd := &cobra.Command{
		Use:   "start <网卡序号>",
		Short: "开始抓包",
		Long:  "在指定网卡上开始抓包，支持多种过滤选项",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			index, err := strconv.Atoi(args[0])
			if err != nil {
				color.Red("无效的网卡序号: %s", args[0])
				return
			}

			selectedIface, err := iface.GetInterfaceByIndex(index)
			if err != nil {
				color.Red("获取网卡信息失败: %v", err)
				return
			}
			if selectedIface == nil {
				color.Red("网卡序号 %d 不存在", index)
				return
			}

			ports := parsePorts(portFlag)

			options := capture.FilterOptions{
				TCP:   tcpFlag,
				UDP:   udpFlag,
				DNS:   dnsFlag,
				HTTP:  httpFlag,
				Ports: ports,
				IP:    ipFlag,
			}

			capturer, err := capture.NewCapturer(selectedIface.Name, options, saveFlag)
			if err != nil {
				color.Red("启动抓包失败: %v", err)
				return
			}
			defer capturer.Close()

			if err := capturer.ApplyBPFFilter(); err != nil {
				color.Red("应用过滤规则失败: %v", err)
				return
			}

			color.Green("开始在网卡 %s (%s) 上抓包...", selectedIface.Name, selectedIface.Desc)
			color.Yellow("按 Ctrl+C 停止抓包")
			fmt.Println("-" + strings.Repeat("-", 100))
			fmt.Printf(" %-12s | %-15s | %-15s | %-6s | %-6s | %-4s | %s\n",
				"时间", "源IP", "目标IP", "源端口", "目标端口", "协议", "内容")
			fmt.Println("-" + strings.Repeat("-", 100))

			sigChan := make(chan os.Signal, 1)
			signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

			go func() {
				<-sigChan
				color.Yellow("\n正在停止抓包...")
				capturer.Stop()
			}()

			err = capturer.Capture(func(info *packet.PacketInfo) {
				printPacket(info)
			})

			if err != nil {
				color.Red("抓包过程出错: %v", err)
				return
			}

			printStatistics(capturer.GetStats())
		},
	}

	cmd.Flags().BoolVar(&tcpFlag, "tcp", false, "只抓取 TCP 包")
	cmd.Flags().BoolVar(&udpFlag, "udp", false, "只抓取 UDP 包")
	cmd.Flags().BoolVar(&dnsFlag, "dns", false, "只抓取 DNS 包")
	cmd.Flags().BoolVar(&httpFlag, "http", false, "只抓取 HTTP 包")
	cmd.Flags().StringVar(&portFlag, "port", "", "按端口过滤，支持多个端口，用逗号分隔")
	cmd.Flags().StringVar(&ipFlag, "ip", "", "按 IP 地址过滤")
	cmd.Flags().StringVar(&saveFlag, "save", "", "将抓包数据保存为 pcap 文件")

	return cmd
}

func parsePorts(portStr string) []uint16 {
	if portStr == "" {
		return nil
	}
	parts := strings.Split(portStr, ",")
	var ports []uint16
	for _, part := range parts {
		p, err := strconv.Atoi(strings.TrimSpace(part))
		if err == nil && p > 0 && p <= 65535 {
			ports = append(ports, uint16(p))
		}
	}
	return ports
}

func printPacket(info *packet.PacketInfo) {
	var protocolColor *color.Color
	switch info.Protocol {
	case "TCP":
		protocolColor = color.New(color.FgBlue)
	case "UDP":
		protocolColor = color.New(color.FgGreen)
	case "DNS":
		protocolColor = color.New(color.FgYellow)
	case "HTTP":
		protocolColor = color.New(color.FgCyan)
	case "ICMP":
		protocolColor = color.New(color.FgRed)
	default:
		protocolColor = color.New(color.FgWhite)
	}

	srcIPColor := color.New(color.FgWhite)
	if packet.IsLocalIP(info.SrcIP) {
		srcIPColor = color.New(color.FgMagenta)
	}

	dstIPColor := color.New(color.FgWhite)
	if packet.IsLocalIP(info.DstIP) {
		dstIPColor = color.New(color.FgMagenta)
	}

	fmt.Printf(" %-12s | %-15s | %-15s | %-6d | %-6d | %-4s | %s\n",
		info.Timestamp,
		srcIPColor.Sprint(info.SrcIP),
		dstIPColor.Sprint(info.DstIP),
		info.SrcPort,
		info.DstPort,
		protocolColor.Sprint(info.Protocol),
		info.Payload)
}

func printStatistics(stats capture.Statistics) {
	fmt.Println("\n" + "-" + strings.Repeat("-", 40))
	color.Green("抓包统计:")
	fmt.Printf("  总数据包: %d\n", stats.TotalPackets)
	fmt.Printf("  TCP: %d\n", stats.TCP)
	fmt.Printf("  UDP: %d\n", stats.UDP)
	fmt.Printf("  DNS: %d\n", stats.DNS)
	fmt.Printf("  HTTP: %d\n", stats.HTTP)
	fmt.Printf("  ICMP: %d\n", stats.ICMP)
	fmt.Printf("  其他: %d\n", stats.Other)
	fmt.Println("-" + strings.Repeat("-", 40))
}