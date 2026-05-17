package terminal

import (
	"fmt"
	"testing"
)

func TestNewRenderer(t *testing.T) {
	r := NewRenderer(80, 24)

	if r.width != 80 {
		t.Errorf("Expected width 80, got %d", r.width)
	}

	if r.height != 24 {
		t.Errorf("Expected height 24, got %d", r.height)
	}

	if r.offset != 0 {
		t.Errorf("Expected offset 0, got %d", r.offset)
	}
}

func TestRendererOffset(t *testing.T) {
	r := NewRenderer(80, 24)

	r.SetOffset(10)
	if r.GetOffset() != 10 {
		t.Errorf("Expected offset 10, got %d", r.GetOffset())
	}

	r.SetOffset(-5)
	if r.GetOffset() != 0 {
		t.Errorf("Expected offset 0 after negative set, got %d", r.GetOffset())
	}

	r.SetOffset(100)
	if r.GetOffset() != 100 {
		t.Errorf("Expected offset 100, got %d", r.GetOffset())
	}
}

func TestRendererDimensions(t *testing.T) {
	r := NewRenderer(100, 30)

	r.SetDimensions(120, 40)

	if r.GetWidth() != 120 {
		t.Errorf("Expected width 120, got %d", r.GetWidth())
	}

	if r.GetHeight() != 40 {
		t.Errorf("Expected height 40, got %d", r.GetHeight())
	}
}

func TestGetVisibleRows(t *testing.T) {
	r := NewRenderer(80, 24)

	visibleRows := r.GetVisibleRows()
	expected := 24 - 12

	if visibleRows != expected {
		t.Errorf("Expected %d visible rows, got %d", expected, visibleRows)
	}

	r.SetDimensions(80, 30)
	visibleRows = r.GetVisibleRows()
	expected = 30 - 12

	if visibleRows != expected {
		t.Errorf("Expected %d visible rows, got %d", expected, visibleRows)
	}
}

func TestColors(t *testing.T) {
	r := NewRenderer(80, 24)

	colors := r.colors

	fmt.Println("Testing colors:")
	fmt.Printf("Header: %sTest%s\n", colors.Header, RESET)
	fmt.Printf("Body: %sTest%s\n", colors.Body, RESET)
	fmt.Printf("Highlight: %sTest%s\n", colors.Highlight, RESET)
	fmt.Printf("CPU: %sTest%s\n", colors.CPU, RESET)
	fmt.Printf("Memory: %sTest%s\n", colors.Memory, RESET)
	fmt.Printf("Swap: %sTest%s\n", colors.Swap, RESET)
	fmt.Printf("Process: %sTest%s\n", colors.Process, RESET)
	fmt.Printf("Warning: %sTest%s\n", colors.Warning, RESET)
	fmt.Printf("Error: %sTest%s\n", colors.Error, RESET)
	fmt.Printf("Success: %sTest%s\n", colors.Success, RESET)

	if colors.Header == "" {
		t.Error("Header color should not be empty")
	}

	if colors.Body == "" {
		t.Error("Body color should not be empty")
	}
}

func TestANSIConstants(t *testing.T) {
	fmt.Println("Testing ANSI constants:")
	fmt.Printf("ESC: %q\n", ESC)
	fmt.Printf("CSI: %q\n", CSI)
	fmt.Printf("BOLD: %q\n", BOLD)
	fmt.Printf("RESET: %q\n", RESET)
	fmt.Printf("RED: %q\n", RED)
	fmt.Printf("GREEN: %q\n", GREEN)
	fmt.Printf("YELLOW: %q\n", YELLOW)
	fmt.Printf("BLUE: %q\n", BLUE)
	fmt.Printf("CLEAR_SCREEN: %q\n", CLEAR_SCREEN)
	fmt.Printf("HIDE_CURSOR: %q\n", HIDE_CURSOR)
	fmt.Printf("SHOW_CURSOR: %q\n", SHOW_CURSOR)

	if ESC != "\033" {
		t.Error("ESC should be \\033")
	}

	if CSI != "\033[" {
		t.Error("CSI should be \\033[")
	}

	if RESET != "\033[0m" {
		t.Error("RESET should be \\033[0m")
	}
}

func TestMoveTo(t *testing.T) {
	r := NewRenderer(80, 24)

	fmt.Println("Testing MoveTo function:")
	r.MoveTo(1, 1)
	fmt.Printf("Moved to row 1, col 1\n")

	r.MoveTo(10, 20)
	fmt.Printf("Moved to row 10, col 20\n")
}

func TestClearFunctions(t *testing.T) {
	fmt.Println("Testing clear functions:")
	fmt.Printf("CLEAR_SCREEN: %s\n", CLEAR_SCREEN)
	fmt.Printf("CLEAR_LINE: %s\n", CLEAR_LINE)
}

func TestMovementSequences(t *testing.T) {
	fmt.Println("Testing movement sequences:")
	fmt.Printf("MOVE_UP: %s\n", MOVE_UP)
	fmt.Printf("MOVE_DOWN: %s\n", MOVE_DOWN)
	fmt.Printf("MOVE_FORWARD: %s\n", MOVE_FORWARD)
	fmt.Printf("MOVE_BACKWARD: %s\n", MOVE_BACKWARD)
	fmt.Printf("MOVE_HOME: %s\n", MOVE_HOME)
}
