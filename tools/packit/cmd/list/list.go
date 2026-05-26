package list

import (
	"fmt"
	"strings"

	"github.com/fatih/color"
	"packit/pkg/iface"
	"github.com/spf13/cobra"
)

func NewListCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "列出所有可用网卡",
		Long:  "列出系统中所有可用的网络接口，供抓包选择使用",
		Run: func(cmd *cobra.Command, args []string) {
			ifaces, err := iface.GetAllInterfaces()
			if err != nil {
				color.Red("获取网卡列表失败: %v", err)
				return
			}

			if len(ifaces) == 0 {
				color.Yellow("未找到可用网卡，请检查网络连接")
				return
			}

			color.Green("可用网卡列表:")
			fmt.Println("-" + strings.Repeat("-", 80))
			fmt.Printf(" %-3s | %-15s | %-18s | %s\n", "序号", "名称", "IP地址", "描述")
			fmt.Println("-" + strings.Repeat("-", 80))

			for _, iface := range ifaces {
				fmt.Printf(" %-3d | %-15s | %-18s | %s\n",
					iface.Index,
					iface.Name,
					iface.IP,
					iface.Desc)
			}

			fmt.Println("-" + strings.Repeat("-", 80))
			color.Cyan("使用命令: packit start <序号> 开始抓包")
		},
	}
	return cmd
}