package process

import (
	"fmt"
	"syscall"
	"unsafe"
)

type ProcessEntry32 struct {
	Size              uint32
	Usage             uint32
	ProcessID         uint32
	DefaultHeapID     uintptr
	ModuleID          uint32
	Threads           uint32
	ParentProcessID   uint32
	PriorityClassBase int32
	Flags             uint32
	ExeFile           [260]uint16
}

type ProcessInfo struct {
	PID      uint32
	Name     string
	CPU      float64
	Memory   uint64
	MemUsage float64
	Threads  uint32
}

type ProcessList struct {
	Processes  []ProcessInfo
	TotalMem   uint64
	TotalCount int
}

var (
	modkernel32                  = syscall.NewLazyDLL("kernel32.dll")
	modpsapi                     = syscall.NewLazyDLL("psapi.dll")
	procCreateToolhelp32Snapshot = modkernel32.NewProc("CreateToolhelp32Snapshot")
	procProcess32First           = modkernel32.NewProc("Process32FirstW")
	procProcess32Next            = modkernel32.NewProc("Process32NextW")
	procOpenProcess              = modkernel32.NewProc("OpenProcess")
	procGetProcessMemoryInfo     = modpsapi.NewProc("GetProcessMemoryInfo")
	procGetProcessTimes          = modkernel32.NewProc("GetProcessTimes")
	procCloseHandle              = modkernel32.NewProc("CloseHandle")
)

const (
	TH32CS_SNAPPROCESS        = 0x00000002
	PROCESS_QUERY_INFORMATION = 0x0400
	INVALID_HANDLE_VALUE      = ^uintptr(0)
)

func CreateToolhelp32Snapshot(flags uint32, processID uint32) uintptr {
	ret, _, _ := procCreateToolhelp32Snapshot.Call(
		uintptr(flags),
		uintptr(processID),
	)
	return ret
}

func Process32First(snapshot uintptr, processEntry *ProcessEntry32) bool {
	processEntry.Size = uint32(unsafe.Sizeof(*processEntry))
	ret, _, _ := procProcess32First.Call(
		snapshot,
		uintptr(unsafe.Pointer(processEntry)),
	)
	return ret != 0
}

func Process32Next(snapshot uintptr, processEntry *ProcessEntry32) bool {
	processEntry.Size = uint32(unsafe.Sizeof(*processEntry))
	ret, _, _ := procProcess32Next.Call(
		snapshot,
		uintptr(unsafe.Pointer(processEntry)),
	)
	return ret != 0
}

func OpenProcess(desiredAccess uint32, inheritHandle bool, processID uint32) uintptr {
	ret, _, _ := procOpenProcess.Call(
		uintptr(desiredAccess),
		uintptr(boolToUintptr(inheritHandle)),
		uintptr(processID),
	)
	return ret
}

func CloseHandle(handle uintptr) bool {
	ret, _, _ := procCloseHandle.Call(handle)
	return ret != 0
}

func boolToUintptr(b bool) uintptr {
	if b {
		return 1
	}
	return 0
}

func GetProcessList() (*ProcessList, error) {
	snapshot := CreateToolhelp32Snapshot(TH32CS_SNAPPROCESS, 0)
	if snapshot == INVALID_HANDLE_VALUE {
		return nil, fmt.Errorf("failed to create process snapshot")
	}
	defer CloseHandle(snapshot)

	var processEntry ProcessEntry32
	processes := make([]ProcessInfo, 0)

	if !Process32First(snapshot, &processEntry) {
		return nil, fmt.Errorf("failed to get first process")
	}

	for {
		name := syscall.UTF16ToString(processEntry.ExeFile[:])

		cpuUsage := GetProcessCPUUsage(processEntry.ProcessID)
		memUsage := GetProcessMemoryUsage(processEntry.ProcessID)

		if name != "" {
			processes = append(processes, ProcessInfo{
				PID:     processEntry.ProcessID,
				Name:    name,
				CPU:     cpuUsage,
				Memory:  memUsage,
				Threads: processEntry.Threads,
			})
		}

		if !Process32Next(snapshot, &processEntry) {
			break
		}
	}

	totalMem := calculateTotalMemory(processes)

	return &ProcessList{
		Processes:  processes,
		TotalMem:   totalMem,
		TotalCount: len(processes),
	}, nil
}

func calculateTotalMemory(processes []ProcessInfo) uint64 {
	var total uint64
	for _, p := range processes {
		total += p.Memory
	}
	return total
}

func GetProcessMemoryUsage(processID uint32) uint64 {
	handle := OpenProcess(PROCESS_QUERY_INFORMATION, false, processID)
	if handle == 0 {
		return 0
	}
	defer CloseHandle(handle)

	type PROCESS_MEMORY_COUNTERS struct {
		CB                         uint32
		PageFaultCount             uint32
		PeakWorkingSetSize         uintptr
		WorkingSetSize             uintptr
		QuotaPeakPagedPoolUsage    uintptr
		PagedPoolUsage             uintptr
		PeakPagedPoolUsage         uintptr
		QuotaPeakNonPagedPoolUsage uintptr
		NonPagedPoolUsage          uintptr
		PageFileUsage              uintptr
		PeakPageFileUsage          uintptr
	}

	var memCounters PROCESS_MEMORY_COUNTERS
	memCounters.CB = uint32(unsafe.Sizeof(memCounters))

	ret, _, _ := procGetProcessMemoryInfo.Call(
		uintptr(handle),
		uintptr(unsafe.Pointer(&memCounters)),
		uintptr(memCounters.CB),
	)

	if ret == 0 {
		return 0
	}

	return uint64(memCounters.WorkingSetSize)
}

var lastProcessTimes = make(map[uint32]struct {
	lastTime uint64
	lastCPU  float64
})

func GetProcessCPUUsage(processID uint32) float64 {
	handle := OpenProcess(PROCESS_QUERY_INFORMATION, false, processID)
	if handle == 0 {
		return 0
	}
	defer CloseHandle(handle)

	type FILETIME struct {
		DwLowDateTime  uint32
		DwHighDateTime uint32
	}

	var creationTime, exitTime, kernelTime, userTime FILETIME

	ret, _, _ := procGetProcessTimes.Call(
		uintptr(handle),
		uintptr(unsafe.Pointer(&creationTime)),
		uintptr(unsafe.Pointer(&exitTime)),
		uintptr(unsafe.Pointer(&kernelTime)),
		uintptr(unsafe.Pointer(&userTime)),
	)

	if ret == 0 {
		return 0
	}

	currentTime := (uint64(kernelTime.DwHighDateTime)<<32 | uint64(kernelTime.DwLowDateTime)) +
		(uint64(userTime.DwHighDateTime)<<32 | uint64(userTime.DwLowDateTime))

	if last, exists := lastProcessTimes[processID]; exists {
		timeDiff := currentTime - last.lastTime
		if timeDiff > 0 {
			cpuUsage := float64(timeDiff) / 10000000.0 / 1.0 * 100.0
			if cpuUsage > 100 {
				cpuUsage = 100
			}
			lastProcessTimes[processID] = struct {
				lastTime uint64
				lastCPU  float64
			}{currentTime, cpuUsage}
			return cpuUsage
		}
	}

	lastProcessTimes[processID] = struct {
		lastTime uint64
		lastCPU  float64
	}{currentTime, 0}

	return 0
}

func SortByCPU(processes []ProcessInfo) []ProcessInfo {
	sorted := make([]ProcessInfo, len(processes))
	copy(sorted, processes)

	for i := 0; i < len(sorted)-1; i++ {
		for j := 0; j < len(sorted)-i-1; j++ {
			if sorted[j].CPU < sorted[j+1].CPU {
				sorted[j], sorted[j+1] = sorted[j+1], sorted[j]
			}
		}
	}

	return sorted
}

func SortByMemory(processes []ProcessInfo) []ProcessInfo {
	sorted := make([]ProcessInfo, len(processes))
	copy(sorted, processes)

	for i := 0; i < len(sorted)-1; i++ {
		for j := 0; j < len(sorted)-i-1; j++ {
			if sorted[j].Memory < sorted[j+1].Memory {
				sorted[j], sorted[j+1] = sorted[j+1], sorted[j]
			}
		}
	}

	return sorted
}

func SortByPID(processes []ProcessInfo) []ProcessInfo {
	sorted := make([]ProcessInfo, len(processes))
	copy(sorted, processes)

	for i := 0; i < len(sorted)-1; i++ {
		for j := 0; j < len(sorted)-i-1; j++ {
			if sorted[j].PID < sorted[j+1].PID {
				sorted[j], sorted[j+1] = sorted[j+1], sorted[j]
			}
		}
	}

	return sorted
}

func CalculateMemoryPercentage(processes []ProcessInfo, totalMem uint64) {
	if totalMem == 0 {
		return
	}

	for i := range processes {
		processes[i].MemUsage = float64(processes[i].Memory) / float64(totalMem) * 100.0
	}
}
