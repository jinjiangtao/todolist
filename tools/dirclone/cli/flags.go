package cli

import (
	"flag"
	"fmt"
	"os"

	"dirclone/config"
)

type Flags struct {
	Src      *string
	Dst      *string
	Depth    *int
	Empty    *bool
	Preserve *bool
	Ignore   *string
	Simulate *bool
	Verbose  *bool
	Help     *bool
}

func Parse() (*config.Config, bool) {
	flags := Flags{
		Src:      flag.String("src", "", "源目录路径 (必填)"),
		Dst:      flag.String("dst", "", "目标目录路径 (必填)"),
		Depth:    flag.Int("depth", 0, "最大目录深度 (0=无限制，默认0)"),
		Empty:    flag.Bool("empty", false, "创建空文件而非空文件夹 (默认false)"),
		Preserve: flag.Bool("preserve", false, "保留原目录属性 (权限/修改时间，默认false)"),
		Ignore:   flag.String("ignore", "", "忽略的目录名，逗号分隔 (如.git,node_modules,cache)"),
		Simulate: flag.Bool("simulate", false, "模拟运行，不实际创建 (默认false)"),
		Verbose:  flag.Bool("v", false, "显示详细日志"),
		Help:     flag.Bool("h", false, "显示帮助信息"),
	}

	flag.Usage = usage

	flag.Parse()

	if *flags.Help {
		usage()
		return nil, true
	}

	if *flags.Src == "" || *flags.Dst == "" {
		fmt.Println("错误: -src 和 -dst 参数都是必填的")
		fmt.Println()
		usage()
		return nil, true
	}

	if _, err := os.Stat(*flags.Src); os.IsNotExist(err) {
		fmt.Printf("错误: 源目录不存在: %s\n", *flags.Src)
		return nil, true
	}

	cfg := config.New(*flags.Src, *flags.Dst)
	cfg.MaxDepth = *flags.Depth
	cfg.CreateEmpty = *flags.Empty
	cfg.PreserveAttrs = *flags.Preserve
	cfg.Simulate = *flags.Simulate
	cfg.Verbose = *flags.Verbose
	cfg.SetIgnorePatterns(*flags.Ignore)

	return cfg, false
}

func usage() {
	fmt.Println("目录结构克隆器 (dirclone)")
	fmt.Println()
	fmt.Println("用法: dirclone [选项]")
	fmt.Println()
	fmt.Println("选项:")
	flag.PrintDefaults()
	fmt.Println()
	fmt.Println("示例:")
	fmt.Println("  dirclone -src /home/project -dst /backup/project")
	fmt.Println("  dirclone -src ./app -dst ./empty-app -depth 3")
	fmt.Println("  dirclone -src . -dst ../clone -ignore \".git,node_modules,tmp\" -simulate")
}
