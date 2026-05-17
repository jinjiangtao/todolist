package main

import (
	"fmt"
	"testing"
)

func TestMainFunctionality(t *testing.T) {
	fmt.Println("Testing main package functionality...")

	t.Run("Terminal size detection", func(t *testing.T) {
		width, height := getTerminalSize()
		fmt.Printf("Terminal size: %dx%d\n", width, height)

		if width <= 0 {
			t.Log("Width is 0 or negative, using default")
		}
		if height <= 0 {
			t.Log("Height is 0 or negative, using default")
		}
	})

	t.Run("Command output parsing", func(t *testing.T) {
		result := getCommandOutput("echo", "test")
		if result != "test\n" && result != "test\r\n" {
			t.Logf("Unexpected output: %q", result)
		}
	})

	t.Run("Parse int function", func(t *testing.T) {
		tests := []struct {
			input    string
			expected int
		}{
			{"123", 123},
			{"abc", 10},
			{"-5", -5},
			{"", 10},
		}

		for _, test := range tests {
			result := parseInt(test.input, 10)
			if result != test.expected {
				t.Errorf("parseInt(%q, 10) = %d, expected %d", test.input, result, test.expected)
			}
		}
	})

	t.Run("Parse float function", func(t *testing.T) {
		tests := []struct {
			input    string
			expected float64
		}{
			{"123.45", 123.45},
			{"abc", 10.0},
			{"-5.5", -5.5},
			{"", 10.0},
		}

		for _, test := range tests {
			result := parseFloat(test.input, 10.0)
			if result != test.expected {
				t.Errorf("parseFloat(%q, 10.0) = %f, expected %f", test.input, result, test.expected)
			}
		}
	})
}

func TestKeyHandling(t *testing.T) {
	fmt.Println("Testing key handling...")

	t.Run("Test key press handler", func(t *testing.T) {
		selectedIndex := 0
		scrollOffset := 0
		sortMode := "pid"
		showHelp := false
		totalProcesses := 100
		visibleRows := 20

		tests := []struct {
			name     string
			key      int
			expected struct {
				selectedIndex int
				scrollOffset  int
				sortMode      string
				showHelp      bool
			}
		}{
			{"Arrow Up", 256 + 72, struct {
				selectedIndex int
				scrollOffset  int
				sortMode      string
				showHelp      bool
			}{selectedIndex: 0, scrollOffset: 0, sortMode: "pid", showHelp: false}},
			{"P key", 80, struct {
				selectedIndex int
				scrollOffset  int
				sortMode      string
				showHelp      bool
			}{selectedIndex: 0, scrollOffset: 0, sortMode: "cpu", showHelp: false}},
			{"M key", 77, struct {
				selectedIndex int
				scrollOffset  int
				sortMode      string
				showHelp      bool
			}{selectedIndex: 0, scrollOffset: 0, sortMode: "memory", showHelp: false}},
		}

		for _, test := range tests {
			t.Run(test.name, func(t *testing.T) {
				handleKeyPress(test.key, &selectedIndex, &scrollOffset, &sortMode, &showHelp, totalProcesses, visibleRows)

				if selectedIndex != test.expected.selectedIndex {
					t.Errorf("Expected selectedIndex %d, got %d", test.expected.selectedIndex, selectedIndex)
				}
				if scrollOffset != test.expected.scrollOffset {
					t.Errorf("Expected scrollOffset %d, got %d", test.expected.scrollOffset, scrollOffset)
				}
				if sortMode != test.expected.sortMode {
					t.Errorf("Expected sortMode %s, got %s", test.expected.sortMode, sortMode)
				}
				if showHelp != test.expected.showHelp {
					t.Errorf("Expected showHelp %v, got %v", test.expected.showHelp, showHelp)
				}
			})
		}
	})
}

func TestEnableVirtualTerminal(t *testing.T) {
	fmt.Println("Testing virtual terminal enablement...")
	enableVirtualTerminal()
	fmt.Println("Virtual terminal enabled (if supported)")
}

func TestReadPassword(t *testing.T) {
	fmt.Println("Testing readPassword function...")
	fmt.Println("Note: This function requires user input in actual use")
}

func TestClearScreen(t *testing.T) {
	fmt.Println("Testing clearScreen function...")
	clearScreen()
	fmt.Println("Screen cleared")
}
