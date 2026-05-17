package process

import (
	"fmt"
	"testing"
	"time"
)

func TestGetProcessList(t *testing.T) {
	processList, err := GetProcessList()
	if err != nil {
		t.Fatalf("Failed to get process list: %v", err)
	}

	fmt.Printf("Total processes: %d\n", processList.TotalCount)
	fmt.Printf("Total memory usage: %d bytes\n", processList.TotalMem)

	if len(processList.Processes) == 0 {
		t.Error("Process list should not be empty")
	}

	for i, p := range processList.Processes[:5] {
		fmt.Printf("Process %d: PID=%d, Name=%s, Memory=%d bytes\n",
			i, p.PID, p.Name, p.Memory)
	}
}

func TestGetProcessMemoryUsage(t *testing.T) {
	processList, err := GetProcessList()
	if err != nil {
		t.Fatalf("Failed to get process list: %v", err)
	}

	for _, p := range processList.Processes {
		mem := GetProcessMemoryUsage(p.PID)
		if mem == 0 && p.PID != 0 {
			t.Logf("Process %d (%s) memory usage is 0", p.PID, p.Name)
		} else {
			fmt.Printf("Process %d (%s): Memory = %d bytes\n", p.PID, p.Name, mem)
		}
	}
}

func TestGetProcessCPUUsage(t *testing.T) {
	processList, err := GetProcessList()
	if err != nil {
		t.Fatalf("Failed to get process list: %v", err)
	}

	fmt.Println("First CPU reading:")
	for i, p := range processList.Processes[:5] {
		cpu := GetProcessCPUUsage(p.PID)
		fmt.Printf("Process %d (%s): CPU = %.2f%%\n", p.PID, p.Name, cpu)
		if i == 0 {
			break
		}
	}

	time.Sleep(1100 * time.Millisecond)

	processList2, err := GetProcessList()
	if err != nil {
		t.Fatalf("Failed to get process list: %v", err)
	}

	fmt.Println("\nSecond CPU reading (after 1.1s):")
	for i, p := range processList2.Processes[:5] {
		cpu := GetProcessCPUUsage(p.PID)
		fmt.Printf("Process %d (%s): CPU = %.2f%%\n", p.PID, p.Name, cpu)
		if i == 0 {
			break
		}
	}
}

func TestSortByCPU(t *testing.T) {
	processList, err := GetProcessList()
	if err != nil {
		t.Fatalf("Failed to get process list: %v", err)
	}

	sorted := SortByCPU(processList.Processes)

	fmt.Println("Processes sorted by CPU:")
	for i, p := range sorted[:10] {
		fmt.Printf("%d. PID=%d, Name=%s, CPU=%.2f%%\n",
			i+1, p.PID, p.Name, p.CPU)
	}

	for i := 0; i < len(sorted)-1; i++ {
		if sorted[i].CPU < sorted[i+1].CPU {
			t.Error("Processes are not properly sorted by CPU")
			break
		}
	}
}

func TestSortByMemory(t *testing.T) {
	processList, err := GetProcessList()
	if err != nil {
		t.Fatalf("Failed to get process list: %v", err)
	}

	sorted := SortByMemory(processList.Processes)

	fmt.Println("Processes sorted by Memory:")
	for i, p := range sorted[:10] {
		fmt.Printf("%d. PID=%d, Name=%s, Memory=%d bytes\n",
			i+1, p.PID, p.Name, p.Memory)
	}

	for i := 0; i < len(sorted)-1; i++ {
		if sorted[i].Memory < sorted[i+1].Memory {
			t.Error("Processes are not properly sorted by Memory")
			break
		}
	}
}

func TestSortByPID(t *testing.T) {
	processList, err := GetProcessList()
	if err != nil {
		t.Fatalf("Failed to get process list: %v", err)
	}

	sorted := SortByPID(processList.Processes)

	fmt.Println("Processes sorted by PID:")
	for i, p := range sorted[:10] {
		fmt.Printf("%d. PID=%d, Name=%s\n", i+1, p.PID, p.Name)
	}

	for i := 0; i < len(sorted)-1; i++ {
		if sorted[i].PID < sorted[i+1].PID {
			t.Error("Processes are not properly sorted by PID")
			break
		}
	}
}

func TestCalculateMemoryPercentage(t *testing.T) {
	processList, err := GetProcessList()
	if err != nil {
		t.Fatalf("Failed to get process list: %v", err)
	}

	totalMem := processList.TotalMem
	if totalMem == 0 {
		totalMem = 8 * 1024 * 1024 * 1024
	}

	CalculateMemoryPercentage(processList.Processes, totalMem)

	fmt.Println("Memory percentages:")
	var totalPercentage float64
	for i, p := range processList.Processes[:10] {
		fmt.Printf("%d. PID=%d, Name=%s, Memory=%.2f%%\n",
			i+1, p.PID, p.Name, p.MemUsage)
		totalPercentage += p.MemUsage
	}

	fmt.Printf("Total percentage of first 10: %.2f%%\n", totalPercentage)
}

func TestProcessInfoStructure(t *testing.T) {
	info := ProcessInfo{
		PID:      1234,
		Name:     "test.exe",
		CPU:      10.5,
		Memory:   1024000,
		MemUsage: 0.01,
		Threads:  5,
	}

	fmt.Printf("Process Info: %+v\n", info)

	if info.PID != 1234 {
		t.Error("PID should be 1234")
	}
	if info.Name != "test.exe" {
		t.Error("Name should be test.exe")
	}
}
