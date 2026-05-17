package system

import (
	"fmt"
	"testing"
	"time"
	"unsafe"
)

func TestGetMemoryInfo(t *testing.T) {
	total, used, free, usage := GetMemoryInfo()

	fmt.Printf("Total: %.0f, Used: %.0f, Free: %.0f, Usage: %.2f%%\n", total, used, free, usage)

	if total == 0 {
		t.Error("Total memory should not be 0")
	}

	if used == 0 {
		t.Error("Used memory should not be 0")
	}

	if free == 0 {
		t.Error("Free memory should not be 0")
	}

	if usage < 0 || usage > 100 {
		t.Errorf("Memory usage should be between 0 and 100, got %.2f", usage)
	}

	expectedUsage := float64(used) / float64(total) * 100
	if usage != expectedUsage {
		t.Errorf("Memory usage calculation incorrect: expected %.2f, got %.2f", expectedUsage, usage)
	}
}

func TestGetSwapInfo(t *testing.T) {
	total, used, free, usage := GetSwapInfo()

	fmt.Printf("Swap Total: %.0f, Used: %.0f, Free: %.0f, Usage: %.2f%%\n", total, used, free, usage)

	if total == 0 {
		t.Log("Swap total is 0, which might be normal for some systems")
	}

	if used > total {
		t.Error("Swap used should not be greater than total")
	}

	if usage < 0 || usage > 100 {
		t.Errorf("Swap usage should be between 0 and 100, got %.2f", usage)
	}
}

func TestGetCPUUsage(t *testing.T) {
	usage := GetCPUUsage()

	fmt.Printf("CPU Usage: %.2f%%\n", usage)

	if usage < 0 || usage > 100 {
		t.Errorf("CPU usage should be between 0 and 100, got %.2f", usage)
	}
}

func TestGetCPUUsageSmooth(t *testing.T) {
	firstUsage := GetCPUUsageSmooth()
	fmt.Printf("First CPU Usage Smooth: %.2f%%\n", firstUsage)

	time.Sleep(1100 * time.Millisecond)

	secondUsage := GetCPUUsageSmooth()
	fmt.Printf("Second CPU Usage Smooth: %.2f%%\n", secondUsage)

	if firstUsage < 0 || firstUsage > 100 {
		t.Errorf("CPU usage should be between 0 and 100, got %.2f", firstUsage)
	}

	if secondUsage < 0 || secondUsage > 100 {
		t.Errorf("CPU usage should be between 0 and 100, got %.2f", secondUsage)
	}
}

func TestGetSystemInfo(t *testing.T) {
	info := GetSystemInfo()

	fmt.Printf("CPU Count: %d\n", info.CPUCount)
	fmt.Printf("CPU Usage: %.2f%%\n", info.CPUUsage)
	fmt.Printf("Memory Total: %d\n", info.MemTotal)
	fmt.Printf("Memory Used: %d\n", info.MemUsed)
	fmt.Printf("Memory Free: %d\n", info.MemFree)
	fmt.Printf("Memory Usage: %.2f%%\n", info.MemUsage)
	fmt.Printf("Swap Total: %d\n", info.SwapTotal)
	fmt.Printf("Swap Used: %d\n", info.SwapUsed)
	fmt.Printf("Swap Free: %d\n", info.SwapFree)
	fmt.Printf("Swap Usage: %.2f%%\n", info.SwapUsage)

	if info.CPUCount <= 0 {
		t.Error("CPU count should be positive")
	}

	if info.MemTotal == 0 {
		t.Error("Memory total should not be 0")
	}

	if info.CPUUsage < 0 || info.CPUUsage > 100 {
		t.Errorf("CPU usage should be between 0 and 100, got %.2f", info.CPUUsage)
	}

	if info.MemUsage < 0 || info.MemUsage > 100 {
		t.Errorf("Memory usage should be between 0 and 100, got %.2f", info.MemUsage)
	}

	if info.SwapUsage < 0 || info.SwapUsage > 100 {
		t.Errorf("Swap usage should be between 0 and 100, got %.2f", info.SwapUsage)
	}
}

func TestFormatBytes(t *testing.T) {
	tests := []struct {
		input    uint64
		expected string
	}{
		{0, "0 B"},
		{512, "512 B"},
		{1023, "1023 B"},
		{1024, "1.0 KB"},
		{1536, "1.5 KB"},
		{1048576, "1.0 MB"},
		{1572864, "1.5 MB"},
		{1073741824, "1.0 GB"},
		{1610612736, "1.5 GB"},
	}

	for _, test := range tests {
		result := FormatBytes(test.input)
		if result != test.expected {
			t.Errorf("FormatBytes(%d) = %s, expected %s", test.input, result, test.expected)
		}
	}
}

func TestMemoryStatusExSize(t *testing.T) {
	var memStatus MemoryStatusEx
	size := unsafe.Sizeof(memStatus)

	expectedSize := uint32(64)
	if size != uintptr(expectedSize) {
		t.Logf("MemoryStatusEx size is %d bytes (expected %d on 64-bit Windows)", size, expectedSize)
	}
}
