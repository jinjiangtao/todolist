package terminal

import (
	"fmt"
	"htop/process"
	"htop/system"
	"os"
)

const (
	ESC            = "\033"
	CSI            = ESC + "["
	BOLD           = CSI + "1m"
	RESET          = CSI + "0m"
	RED            = CSI + "31m"
	GREEN          = CSI + "32m"
	YELLOW         = CSI + "33m"
	BLUE           = CSI + "34m"
	MAGENTA        = CSI + "35m"
	CYAN           = CSI + "36m"
	WHITE          = CSI + "37m"
	BRIGHT_RED     = CSI + "91m"
	BRIGHT_GREEN   = CSI + "92m"
	BRIGHT_YELLOW  = CSI + "93m"
	BRIGHT_BLUE    = CSI + "94m"
	BRIGHT_MAGENTA = CSI + "95m"
	BRIGHT_CYAN    = CSI + "96m"
	BRIGHT_WHITE   = CSI + "97m"
	BG_RED         = CSI + "41m"
	BG_GREEN       = CSI + "42m"
	BG_YELLOW      = CSI + "43m"
	BG_BLUE        = CSI + "44m"
	BG_MAGENTA     = CSI + "45m"
	BG_CYAN        = CSI + "46m"
	BG_WHITE       = CSI + "47m"
	HIDE_CURSOR    = CSI + "?25l"
	SHOW_CURSOR    = CSI + "?25h"
	CLEAR_SCREEN   = CSI + "2J"
	CLEAR_LINE     = CSI + "2K"
	MOVE_HOME      = CSI + "H"
	MOVE_UP        = CSI + "A"
	MOVE_DOWN      = CSI + "B"
	MOVE_FORWARD   = CSI + "C"
	MOVE_BACKWARD  = CSI + "D"
)

type Colors struct {
	Header    string
	Body      string
	Highlight string
	CPU       string
	Memory    string
	Swap      string
	Process   string
	Warning   string
	Error     string
	Success   string
}

var DefaultColors = Colors{
	Header:    BOLD + CYAN,
	Body:      WHITE,
	Highlight: BRIGHT_GREEN,
	CPU:       BRIGHT_GREEN,
	Memory:    BRIGHT_BLUE,
	Swap:      BRIGHT_YELLOW,
	Process:   BRIGHT_WHITE,
	Warning:   BRIGHT_YELLOW,
	Error:     BRIGHT_RED,
	Success:   BRIGHT_GREEN,
}

type Renderer struct {
	width  int
	height int
	colors Colors
	offset int
}

func NewRenderer(width, height int) *Renderer {
	return &Renderer{
		width:  width,
		height: height,
		colors: DefaultColors,
		offset: 0,
	}
}

func (r *Renderer) SetOffset(offset int) {
	r.offset = offset
}

func (r *Renderer) GetOffset() int {
	return r.offset
}

func (r *Renderer) ClearScreen() {
	fmt.Print(CLEAR_SCREEN + MOVE_HOME)
}

func (r *Renderer) HideCursor() {
	fmt.Print(HIDE_CURSOR)
}

func (r *Renderer) ShowCursor() {
	fmt.Print(SHOW_CURSOR)
}

func (r *Renderer) MoveTo(row, col int) {
	fmt.Printf("%s%d;%dH", CSI, row, col)
}

func (r *Renderer) ClearLine() {
	fmt.Print(CLEAR_LINE)
}

func (r *Renderer) DrawHeader() {
	r.MoveTo(1, 1)
	r.ClearLine()
	fmt.Printf("%sHTop - Windows Task Manager%s", r.colors.Header+BOLD, RESET)

	r.MoveTo(1, r.width-20)
	r.ClearLine()
	fmt.Printf("%s[L] 帮助  [Q] 退出%s", r.colors.Body, RESET)
}

func (r *Renderer) DrawSystemInfo(sysInfo *system.SysInfo) {
	r.MoveTo(2, 1)
	r.ClearLine()
	fmt.Printf("%s_CPU: %s%.1f%%%s  |  %s_内存: %s%s / %s (%.1f%%)%s  |  %sSwap: %s%s / %s (%.1f%%)%s",
		r.colors.Header,
		r.colors.CPU, sysInfo.CPUUsage,
		r.colors.Header,
		r.colors.Header,
		r.colors.Memory, system.FormatBytes(sysInfo.MemUsed),
		system.FormatBytes(sysInfo.MemTotal), sysInfo.MemUsage,
		r.colors.Header,
		r.colors.Header,
		r.colors.Swap, system.FormatBytes(sysInfo.SwapUsed),
		system.FormatBytes(sysInfo.SwapTotal), sysInfo.SwapUsage,
		r.colors.Header,
	)

	r.MoveTo(3, 1)
	r.ClearLine()
	fmt.Printf("%s_CPU 核心数: %d%s", r.colors.Body, sysInfo.CPUCount, RESET)

	r.DrawCPUBar(sysInfo.CPUUsage)
	r.DrawMemoryBar(sysInfo.MemUsage)
	r.DrawSwapBar(sysInfo.SwapUsage)

	r.DrawProcessHeader()
}

func (r *Renderer) DrawCPUBar(usage float64) {
	barWidth := 30
	filled := int(usage / 100.0 * float64(barWidth))

	bar := "["
	for i := 0; i < barWidth; i++ {
		if i < filled {
			if usage > 80 {
				bar += r.colors.Error + "█"
			} else if usage > 50 {
				bar += r.colors.Warning + "█"
			} else {
				bar += r.colors.Success + "█"
			}
		} else {
			bar += r.colors.Body + "░"
		}
	}
	bar += RESET + "]"

	r.MoveTo(4, 1)
	r.ClearLine()
	fmt.Printf("%sCPU:%s %s %.1f%%", r.colors.Header, RESET, bar, usage)
}

func (r *Renderer) DrawMemoryBar(usage float64) {
	barWidth := 30
	filled := int(usage / 100.0 * float64(barWidth))

	bar := "["
	for i := 0; i < barWidth; i++ {
		if i < filled {
			if usage > 80 {
				bar += r.colors.Error + "█"
			} else if usage > 50 {
				bar += r.colors.Warning + "█"
			} else {
				bar += r.colors.Success + "█"
			}
		} else {
			bar += r.colors.Body + "░"
		}
	}
	bar += RESET + "]"

	r.MoveTo(5, 1)
	r.ClearLine()
	fmt.Printf("%s内存:%s %s %.1f%%", r.colors.Header, RESET, bar, usage)
}

func (r *Renderer) DrawSwapBar(usage float64) {
	barWidth := 30
	filled := int(usage / 100.0 * float64(barWidth))

	bar := "["
	for i := 0; i < barWidth; i++ {
		if i < filled {
			if usage > 80 {
				bar += r.colors.Error + "█"
			} else if usage > 50 {
				bar += r.colors.Warning + "█"
			} else {
				bar += r.colors.Success + "█"
			}
		} else {
			bar += r.colors.Body + "░"
		}
	}
	bar += RESET + "]"

	r.MoveTo(6, 1)
	r.ClearLine()
	fmt.Printf("%sSwap:%s %s %.1f%%", r.colors.Header, RESET, bar, usage)
}

func (r *Renderer) DrawProcessHeader() {
	r.MoveTo(7, 1)
	r.ClearLine()
	fmt.Printf("%s  PID   | 进程名称                    |  CPU  | 内存        | 线程数%s",
		r.colors.Header+BOLD, RESET)

	for i := 8; i <= 10; i++ {
		r.MoveTo(i, 1)
		r.ClearLine()
	}
}

func (r *Renderer) DrawProcesses(processList *process.ProcessList, sortMode string, selectedIndex int) {
	var processes []process.ProcessInfo
	switch sortMode {
	case "cpu":
		processes = process.SortByCPU(processList.Processes)
	case "memory":
		processes = process.SortByMemory(processList.Processes)
	default:
		processes = process.SortByPID(processList.Processes)
	}

	process.CalculateMemoryPercentage(processes, processList.TotalMem)

	startRow := 8
	maxRows := r.height - 12
	endRow := startRow + maxRows

	if r.offset > len(processes)-maxRows {
		r.offset = len(processes) - maxRows
	}
	if r.offset < 0 {
		r.offset = 0
	}

	end := r.offset + maxRows
	if end > len(processes) {
		end = len(processes)
	}

	for i := startRow; i < endRow; i++ {
		r.MoveTo(i, 1)
		r.ClearLine()

		processIndex := i - startRow + r.offset
		if processIndex >= len(processes) {
			break
		}

		p := processes[processIndex]

		if processIndex == selectedIndex {
			fmt.Printf("%s>%s", r.colors.Highlight+BOLD, RESET)
		} else {
			fmt.Print(" ")
		}

		name := p.Name
		if len(name) > 26 {
			name = name[:23] + "..."
		}

		cpuStr := fmt.Sprintf("%5.1f%%", p.CPU)
		if p.CPU < 0.1 {
			cpuStr = "  0.0%"
		}

		memStr := system.FormatBytes(p.Memory)
		if len(memStr) < 10 {
			memStr = fmt.Sprintf("%-10s", memStr)
		}

		fmt.Printf("  %5d  | %-26s | %s | %s | %5d",
			p.PID,
			name,
			cpuStr,
			memStr,
			p.Threads)

		if processIndex == selectedIndex {
			fmt.Print(" " + r.colors.Highlight + "◄" + RESET)
		}
	}

	for i := endRow; i < r.height-1; i++ {
		r.MoveTo(i, 1)
		r.ClearLine()
	}

	r.DrawFooter(len(processes), selectedIndex, sortMode)
}

func (r *Renderer) DrawFooter(totalProcesses, selectedIndex int, sortMode string) {
	r.MoveTo(r.height-1, 1)
	r.ClearLine()

	sortLabel := "PID"
	switch sortMode {
	case "cpu":
		sortLabel = "CPU"
	case "memory":
		sortLabel = "内存"
	}

	fmt.Printf("%s[%s]%s 排序  %s[P]%s CPU排序  %s[M]%s 内存排序  %s [↑↓]%s 移动  %s[Q]%s 退出",
		r.colors.Warning, sortLabel, RESET,
		r.colors.Body, RESET,
		r.colors.Body, RESET,
		r.colors.Body, RESET,
		r.colors.Error+BOLD, RESET)

	r.MoveTo(r.height, 1)
	r.ClearLine()
	fmt.Printf("进程: %d/%d  |  排序: %s", selectedIndex+1, totalProcesses, sortLabel)
}

func (r *Renderer) DrawHelp() {
	r.ClearScreen()
	r.HideCursor()

	r.MoveTo(1, 1)
	fmt.Printf("%s=== Windows HTOP 帮助 ===%s\n\n", r.colors.Header+BOLD, RESET)

	fmt.Printf("%s快捷键:%s\n", r.colors.Highlight, RESET)
	fmt.Printf("  ↑ / ↓     - 上下滚动进程列表\n")
	fmt.Printf("  P          - 按 CPU 使用率排序\n")
	fmt.Printf("  M          - 按内存使用排序\n")
	fmt.Printf("  L          - 显示帮助信息\n")
	fmt.Printf("  Q          - 退出程序\n\n")

	fmt.Printf("%s颜色说明:%s\n", r.colors.Highlight, RESET)
	fmt.Printf("  %s绿色%s  - CPU/内存使用率低 (0-50%%)\n", r.colors.Success, RESET)
	fmt.Printf("  %s黄色%s  - CPU/内存使用率中 (50-80%%)\n", r.colors.Warning, RESET)
	fmt.Printf("  %s红色%s  - CPU/内存使用率高 (80-100%%)\n\n", r.colors.Error, RESET)

	fmt.Printf("%s按任意键返回...%s", r.colors.Body, RESET)

	os.Stdout.Sync()
}

func (r *Renderer) GetVisibleRows() int {
	return r.height - 12
}

func (r *Renderer) SetDimensions(width, height int) {
	r.width = width
	r.height = height
}

func (r *Renderer) GetWidth() int {
	return r.width
}

func (r *Renderer) GetHeight() int {
	return r.height
}
