package input

import (
	"fmt"
	"testing"
	"unsafe"
)

func TestKeyConstants(t *testing.T) {
	fmt.Println("Testing key constants:")

	fmt.Printf("KeyUp: %d\n", KeyUp)
	fmt.Printf("KeyDown: %d\n", KeyDown)
	fmt.Printf("KeyEnter: %d\n", KeyEnter)
	fmt.Printf("KeyEscape: %d\n", KeyEscape)
	fmt.Printf("KeyArrowUp: %d\n", KeyArrowUp)
	fmt.Printf("KeyArrowDown: %d\n", KeyArrowDown)
	fmt.Printf("KeyF1: %d\n", KeyF1)
	fmt.Printf("KeyF2: %d\n", KeyF2)

	if KeyUp <= 0 {
		t.Error("KeyUp should be positive")
	}

	if KeyDown <= 0 {
		t.Error("KeyDown should be positive")
	}

	if KeyEnter != 13 {
		t.Errorf("KeyEnter should be 13, got %d", KeyEnter)
	}

	if KeyEscape != 27 {
		t.Errorf("KeyEscape should be 27, got %d", KeyEscape)
	}
}

func TestNewInputHandler(t *testing.T) {
	fmt.Println("Testing NewInputHandler...")

	handler, err := NewInputHandler()
	if err != nil {
		t.Errorf("Failed to create input handler: %v", err)
		return
	}

	if handler == nil {
		t.Error("Input handler should not be nil")
	}

	if handler.oldState == nil {
		t.Error("Old state should not be nil")
	}

	fmt.Printf("Input handler created successfully\n")
	fmt.Printf("Old state mode: %d\n", handler.oldState.mode)

	handler.Restore()
	fmt.Println("Input handler restored")
}

func TestInputHandlerRestore(t *testing.T) {
	fmt.Println("Testing InputHandler.Restore...")

	handler, err := NewInputHandler()
	if err != nil {
		t.Errorf("Failed to create input handler: %v", err)
		return
	}

	handler.Restore()
	fmt.Println("Input handler restored successfully")
}

func TestGetStdHandle(t *testing.T) {
	fmt.Println("Testing GetStdHandle...")

	handle := GetStdHandle(STD_INPUT_HANDLE)

	if handle == 0 {
		t.Error("Std handle should not be 0")
	}

	fmt.Printf("Std input handle: %d\n", handle)
}

func TestGetConsoleMode(t *testing.T) {
	fmt.Println("Testing GetConsoleMode...")

	stdinHandle := GetStdHandle(STD_INPUT_HANDLE)
	var mode uint32

	result := GetConsoleMode(stdinHandle, &mode)

	if !result {
		t.Error("GetConsoleMode should return true")
	}

	fmt.Printf("Console mode: %d\n", mode)
}

func TestConstants(t *testing.T) {
	fmt.Println("Testing constants:")
	fmt.Printf("STD_INPUT_HANDLE: %d\n", STD_INPUT_HANDLE)
	fmt.Printf("ENABLE_ECHO_INPUT: %d\n", ENABLE_ECHO_INPUT)
	fmt.Printf("ENABLE_LINE_INPUT: %d\n", ENABLE_LINE_INPUT)
	fmt.Printf("ENABLE_MOUSE_INPUT: %d\n", ENABLE_MOUSE_INPUT)
	fmt.Printf("ENABLE_WINDOW_INPUT: %d\n", ENABLE_WINDOW_INPUT)
	fmt.Printf("ENABLE_PROCESSED_INPUT: %d\n", ENABLE_PROCESSED_INPUT)

	if STD_INPUT_HANDLE != -10 {
		t.Errorf("STD_INPUT_HANDLE should be -10, got %d", STD_INPUT_HANDLE)
	}
}

func TestKeyType(t *testing.T) {
	fmt.Println("Testing Key type...")

	var key Key = KeyUp
	fmt.Printf("Key value: %d\n", key)

	if key <= 0 {
		t.Error("Key should be positive")
	}

	key = Key('A')
	if key != 65 {
		t.Errorf("Key('A') should be 65, got %d", key)
	}

	key = Key('Q')
	if key != 81 {
		t.Errorf("Key('Q') should be 81, got %d", key)
	}
}

func TestTermState(t *testing.T) {
	fmt.Println("Testing termState...")

	state := &termState{
		mode: 0x0103,
	}

	if state.mode == 0 {
		t.Error("Mode should not be 0")
	}

	fmt.Printf("Term state mode: %d\n", state.mode)
}

func TestINPUT_RECORD(t *testing.T) {
	fmt.Println("Testing INPUT_RECORD...")

	var input INPUT_RECORD
	input.EventType = 1

	fmt.Printf("Input record event type: %d\n", input.EventType)

	if input.EventType != 1 {
		t.Error("Event type should be 1 for keyboard event")
	}
}

func TestKEY_EVENT_RECORD(t *testing.T) {
	fmt.Println("Testing KEY_EVENT_RECORD...")

	keyEvent := (*KEY_EVENT_RECORD)(unsafe.Pointer(&[16]byte{}))
	keyEvent.bKeyDown = 1
	keyEvent.wVirtualKeyCode = 81
	keyEvent.UnicodeChar = 'Q'

	fmt.Printf("Key down: %d\n", keyEvent.bKeyDown)
	fmt.Printf("Virtual key code: %d\n", keyEvent.wVirtualKeyCode)
	fmt.Printf("Unicode char: %c\n", keyEvent.UnicodeChar)
}
