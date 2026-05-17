package main

import (
	"bufio"
	"fmt"
	"htop/input"
	"htop/process"
	"htop/system"
	"htop/terminal"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"syscall"
	"time"
	"unsafe"
)

var (
	kernel32                        = syscall.NewLazyDLL("kernel32.dll")
	procGetConsoleScreenBufferInfo  = kernel32.NewProc("GetConsoleScreenBufferInfo")
	procSetConsoleCursorPosition    = kernel32.NewProc("SetConsoleCursorPosition")
	procFillConsoleOutputAttribute  = kernel32.NewProc("FillConsoleOutputAttribute")
	procFillConsoleOutputCharacterW = kernel32.NewProc("FillConsoleOutputCharacterW")
)

type COORD struct {
	X int16
	Y int16
}

type SMALL_RECT struct {
	Left   int16
	Top    int16
	Right  int16
	Bottom int16
}

type CONSOLE_SCREEN_BUFFER_INFO struct {
	dwSize              COORD
	dwCursorPosition    COORD
	wAttributes         uint16
	srWindow            SMALL_RECT
	dwMaximumWindowSize COORD
}

func GetConsoleScreenBufferInfo(handle uintptr) (CONSOLE_SCREEN_BUFFER_INFO, bool) {
	var info CONSOLE_SCREEN_BUFFER_INFO
	ret, _, _ := procGetConsoleScreenBufferInfo.Call(handle, uintptr(unsafe.Pointer(&info)))
	return info, ret != 0
}

func main() {
	fmt.Printf("%sWindows HTOP - 进程监控工具%s\n\n", terminal.BOLD+terminal.CYAN, terminal.RESET)
	fmt.Println("正在初始化...")

	inputHandler, err := input.NewInputHandler()
	if err != nil {
		fmt.Printf("警告: 无法初始化输入处理: %v\n", err)
	}

	width, height := getTerminalSize()
	if width == 0 || height == 0 {
		width = 120
		height = 30
	}

	fmt.Printf("终端大小: %dx%d\n", width, height)
	fmt.Println("正在获取系统信息...")

	sysInfo := system.GetSystemInfo()
	fmt.Printf("CPU 核心数: %d\n", sysInfo.CPUCount)
	fmt.Printf("总内存: %s\n", system.FormatBytes(sysInfo.MemTotal))

	time.Sleep(1 * time.Second)

	renderer := terminal.NewRenderer(width, height)
	renderer.HideCursor()

	fmt.Print(terminal.CLEAR_SCREEN + terminal.MOVE_HOME)

	sortMode := "pid"
	selectedIndex := 0
	scrollOffset := 0
	showHelp := false

	processList, err := process.GetProcessList()
	if err != nil {
		fmt.Printf("错误: 无法获取进程列表: %v\n", err)
		processList = &process.ProcessList{
			Processes:  []process.ProcessInfo{},
			TotalMem:   0,
			TotalCount: 0,
		}
	}

	visibleRows := renderer.GetVisibleRows()

	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	inputChan := make(chan input.Key, 10)

	go func() {
		for {
			key := input.GetKeyPressed()
			inputChan <- key
		}
	}()

	renderTicker := time.NewTicker(100 * time.Millisecond)
	defer renderTicker.Stop()

	refreshTicker := time.NewTicker(2 * time.Second)
	defer refreshTicker.Stop()

	for {
		select {
		case key := <-inputChan:
			if showHelp {
				showHelp = false
				renderer.ShowCursor()
				renderer.ClearScreen()
				renderer.HideCursor()
				continue
			}

			handleKeyPress(key, &selectedIndex, &scrollOffset, &sortMode, &showHelp, len(processList.Processes), visibleRows)

			if key == input.Key('Q') || key == input.KeyCtrlQ {
				renderer.ShowCursor()
				fmt.Print(terminal.CLEAR_SCREEN + terminal.MOVE_HOME)
				fmt.Println("\n感谢使用 Windows HTOP！")
				if inputHandler != nil {
					inputHandler.Restore()
				}
				return
			}

		case <-ticker.C:
			if !showHelp {
				sysInfo = system.GetSystemInfo()
			}

		case <-refreshTicker.C:
			if !showHelp {
				newProcessList, err := process.GetProcessList()
				if err == nil {
					processList = newProcessList
				}

				visibleRows = renderer.GetVisibleRows()
				if scrollOffset > 0 && scrollOffset > len(processList.Processes)-visibleRows {
					scrollOffset = len(processList.Processes) - visibleRows
				}
				if scrollOffset < 0 {
					scrollOffset = 0
				}
			}

		case <-renderTicker.C:
			if showHelp {
				renderer.DrawHelp()
			} else {
				width, height = getTerminalSize()
				if width > 0 && height > 0 {
					renderer.SetDimensions(width, height)
				}
				renderer.SetOffset(scrollOffset)

				renderer.DrawHeader()
				renderer.DrawSystemInfo(sysInfo)
				renderer.DrawProcesses(processList, sortMode, selectedIndex+scrollOffset)
			}
		}
	}
}

func handleKeyPress(key input.Key, selectedIndex, scrollOffset *int, sortMode *string, showHelp *bool, totalProcesses, visibleRows int) {
	switch key {
	case input.KeyArrowUp, input.KeyUp:
		if *selectedIndex > 0 {
			*selectedIndex--
			if *scrollOffset > 0 && *selectedIndex < *scrollOffset {
				*scrollOffset--
			}
		}

	case input.KeyArrowDown, input.KeyDown:
		if *selectedIndex < totalProcesses-1 && *selectedIndex < visibleRows-1 {
			*selectedIndex++
			if *selectedIndex >= visibleRows && *scrollOffset+*selectedIndex < totalProcesses {
				*scrollOffset++
				*selectedIndex = visibleRows - 1
			}
		} else if *scrollOffset+*selectedIndex < totalProcesses-1 {
			*scrollOffset++
		}

	case input.Key('P'):
		*sortMode = "cpu"

	case input.Key('M'):
		*sortMode = "memory"

	case input.Key('L'), input.KeyF1:
		*showHelp = true

	case input.KeyHome:
		*selectedIndex = 0
		*scrollOffset = 0

	case input.KeyEnd:
		*scrollOffset = totalProcesses - visibleRows
		if *scrollOffset < 0 {
			*scrollOffset = 0
		}
		*selectedIndex = visibleRows - 1
		if *selectedIndex >= totalProcesses {
			*selectedIndex = totalProcesses - 1
		}

	case input.KeyPageUp:
		*scrollOffset -= visibleRows
		if *scrollOffset < 0 {
			*scrollOffset = 0
		}

	case input.KeyPageDown:
		maxScroll := totalProcesses - visibleRows
		if maxScroll < 0 {
			maxScroll = 0
		}
		*scrollOffset += visibleRows
		if *scrollOffset > maxScroll {
			*scrollOffset = maxScroll
		}
	}
}

func getTerminalSize() (int, int) {
	kernel32 := syscall.NewLazyDLL("kernel32.dll")
	procGetConsoleScreenBufferInfo := kernel32.NewProc("GetConsoleScreenBufferInfo")
	procGetStdHandle := kernel32.NewProc("GetStdHandle")

	stdoutHandle, _, _ := procGetStdHandle.Call(uintptr(11))

	var info CONSOLE_SCREEN_BUFFER_INFO
	ret, _, _ := procGetConsoleScreenBufferInfo.Call(stdoutHandle, uintptr(unsafe.Pointer(&info)))

	if ret == 0 {
		return 0, 0
	}

	width := int(info.srWindow.Right - info.srWindow.Left + 1)
	height := int(info.srWindow.Bottom - info.srWindow.Top + 1)

	return width, height
}

func clearScreen() {
	cmd := exec.Command("cmd", "/c", "cls")
	cmd.Stdout = os.Stdout
	cmd.Run()
}

func enableVirtualTerminal() {
	kernel32 := syscall.NewLazyDLL("kernel32.dll")
	procGetStdHandle := kernel32.NewProc("GetStdHandle")
	procGetConsoleMode := kernel32.NewProc("GetConsoleMode")
	procSetConsoleMode := kernel32.NewProc("SetConsoleMode")

	stdoutHandle, _, _ := procGetStdHandle.Call(uintptr(uint32(0xfffffff5)))
	stdinHandle, _, _ := procGetStdHandle.Call(uintptr(uint32(0xfffffff6)))

	var mode uint32
	procGetConsoleMode.Call(stdinHandle, uintptr(unsafe.Pointer(&mode)))

	VT_ENABLED := uint32(0x0004)
	procSetConsoleMode.Call(stdoutHandle, uintptr(mode|VT_ENABLED))
	procSetConsoleMode.Call(stdinHandle, uintptr(mode|VT_ENABLED))
}

func readPassword(prompt string) string {
	fmt.Print(prompt)
	reader := bufio.NewReader(os.Stdin)
	text, _ := reader.ReadString('\n')
	return strings.TrimSpace(text)
}

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
}

func getCommandOutput(command string, args ...string) string {
	cmd := exec.Command(command, args...)
	output, err := cmd.Output()
	if err != nil {
		return ""
	}
	return string(output)
}

func parseInt(s string, defaultVal int) int {
	val, err := strconv.Atoi(s)
	if err != nil {
		return defaultVal
	}
	return val
}

func parseFloat(s string, defaultVal float64) float64 {
	val, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return defaultVal
	}
	return val
}
