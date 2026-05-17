package system

import (
	"fmt"
	"runtime"
	"syscall"
	"time"
	"unsafe"
)

type MemoryStatusEx struct {
	dwLength                uint32
	dwMemoryLoad            uint32
	ullTotalPhys            uint64
	ullAvailPhys            uint64
	ullTotalPageFile        uint64
	ullAvailPageFile        uint64
	ullTotalVirtual         uint64
	ullAvailVirtual         uint64
	ullAvailExtendedVirtual uint64
}

var (
	modkernel32              = syscall.NewLazyDLL("kernel32.dll")
	procGlobalMemoryStatusEx = modkernel32.NewProc("GlobalMemoryStatusEx")
	procGetSystemTimes       = modkernel32.NewProc("GetSystemTimes")
)

func GlobalMemoryStatusEx(m *MemoryStatusEx) bool {
	ret, _, _ := procGlobalMemoryStatusEx.Call(
		uintptr(unsafe.Pointer(m)),
	)
	return ret != 0
}

func GetSystemTimes(idleTime, kernelTime, userTime *uint64) bool {
	ret, _, _ := procGetSystemTimes.Call(
		uintptr(unsafe.Pointer(idleTime)),
		uintptr(unsafe.Pointer(kernelTime)),
		uintptr(unsafe.Pointer(userTime)),
	)
	return ret != 0
}

func GetCPUUsage() float64 {
	var idleTime, kernelTime, userTime uint64

	if !GetSystemTimes(&idleTime, &kernelTime, &userTime) {
		return 0
	}

	idleTime = idleTime / 10000000
	kernelTime = kernelTime / 10000000
	userTime = userTime / 10000000

	totalTime := kernelTime + userTime + idleTime

	if totalTime == 0 {
		return 0
	}

	return float64(totalTime-idleTime) / float64(totalTime) * 100.0
}

func GetMemoryInfo() (total, used, free, usage float64) {
	var memStatus MemoryStatusEx
	memStatus.dwLength = uint32(unsafe.Sizeof(memStatus))

	if !GlobalMemoryStatusEx(&memStatus) {
		return 0, 0, 0, 0
	}

	total = float64(memStatus.ullTotalPhys)
	used = float64(memStatus.ullTotalPhys - memStatus.ullAvailPhys)
	free = float64(memStatus.ullAvailPhys)

	if total > 0 {
		usage = (used / total) * 100.0
	}

	return total, used, free, usage
}

func GetSwapInfo() (total, used, free, usage float64) {
	var memStatus MemoryStatusEx
	memStatus.dwLength = uint32(unsafe.Sizeof(memStatus))

	if !GlobalMemoryStatusEx(&memStatus) {
		return 0, 0, 0, 0
	}

	total = float64(memStatus.ullTotalPageFile)
	used = float64(memStatus.ullTotalPageFile - memStatus.ullAvailPageFile)
	free = float64(memStatus.ullAvailPageFile)

	if total > 0 {
		usage = (used / total) * 100.0
	}

	return total, used, free, usage
}

type SysInfo struct {
	CPUUsage  float64
	MemTotal  uint64
	MemUsed   uint64
	MemFree   uint64
	MemUsage  float64
	SwapTotal uint64
	SwapUsed  uint64
	SwapFree  uint64
	SwapUsage float64
	CPUCount  int
}

func GetSystemInfo() *SysInfo {
	cpuUsage := GetCPUUsage()
	memTotal, memUsed, memFree, memUsage := GetMemoryInfo()
	swapTotal, swapUsed, swapFree, swapUsage := GetSwapInfo()

	cpuCount := runtime.NumCPU()

	return &SysInfo{
		CPUUsage:  cpuUsage,
		MemTotal:  uint64(memTotal),
		MemUsed:   uint64(memUsed),
		MemFree:   uint64(memFree),
		MemUsage:  memUsage,
		SwapTotal: uint64(swapTotal),
		SwapUsed:  uint64(swapUsed),
		SwapFree:  uint64(swapFree),
		SwapUsage: swapUsage,
		CPUCount:  cpuCount,
	}
}

var lastCPUCheck time.Time
var lastCPUUsage float64

func GetCPUUsageSmooth() float64 {
	now := time.Now()
	elapsed := now.Sub(lastCPUCheck)

	if elapsed < time.Second {
		return lastCPUUsage
	}

	currentUsage := GetCPUUsage()

	if lastCPUCheck.IsZero() {
		lastCPUUsage = currentUsage
		lastCPUCheck = now
		return currentUsage
	}

	smoothing := float64(elapsed) / float64(time.Second)
	lastCPUUsage = lastCPUUsage*(1-smoothing) + currentUsage*smoothing
	lastCPUCheck = now

	return lastCPUUsage
}

func FormatBytes(bytes uint64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}

	div, exp := uint64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}

	units := []string{"KB", "MB", "GB", "TB", "PB"}
	if exp >= len(units) {
		exp = len(units) - 1
	}

	value := float64(bytes) / float64(div)
	return fmt.Sprintf("%.1f %s", value, units[exp])
}
