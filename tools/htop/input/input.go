package input

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"unsafe"
)

type Key int

const (
	KeyUp         Key = 1
	KeyDown       Key = 2
	KeyEnter      Key = 13
	KeyEscape     Key = 27
	KeySpace      Key = 32
	KeyTab        Key = 9
	KeyBackspace  Key = 8
	KeyDelete     Key = 127
	KeyHome       Key = 71
	KeyEnd        Key = 79
	KeyPageUp     Key = 73
	KeyPageDown   Key = 81
	KeyArrowUp    Key = 256 + 72
	KeyArrowDown  Key = 256 + 80
	KeyArrowLeft  Key = 256 + 75
	KeyArrowRight Key = 256 + 77
	KeyInsert     Key = 82
	KeyF1         Key = 256 + 59
	KeyF2         Key = 256 + 60
	KeyF3         Key = 256 + 61
	KeyF4         Key = 256 + 62
	KeyF5         Key = 256 + 63
	KeyF6         Key = 256 + 64
	KeyF7         Key = 256 + 65
	KeyF8         Key = 256 + 66
	KeyF9         Key = 256 + 67
	KeyF10        Key = 256 + 68
	KeyCtrlC      Key = 3
	KeyCtrlQ      Key = 17
	KeyCtrlZ      Key = 26
	KeyNull       Key = 0
)

const (
	STD_INPUT_HANDLE       = -10
	ENABLE_ECHO_INPUT      = 0x0004
	ENABLE_LINE_INPUT      = 0x0002
	ENABLE_MOUSE_INPUT     = 0x0010
	ENABLE_WINDOW_INPUT    = 0x0008
	ENABLE_PROCESSED_INPUT = 0x0001
)

type InputHandler struct {
	oldState *termState
}

type termState struct {
	mode uint32
}

type KEY_EVENT_RECORD struct {
	bKeyDown          int32
	wRepeatCount      uint16
	wVirtualKeyCode   uint16
	wVirtualScanCode  uint16
	UnicodeChar       uint16
	dwControlKeyState uint32
}

type MOUSE_EVENT_RECORD struct {
	dwMousePosX       int32
	dwMousePosY       int32
	dwButtonState     uint32
	dwControlKeyState uint32
	dwEventFlags      uint32
}

type WINDOW_BUFFER_SIZE_RECORD struct {
	dwSizeX int32
	dwSizeY int32
}

type MENU_EVENT_RECORD struct {
	dwCommandId uint32
}

type FOCUS_EVENT_RECORD struct {
	bSetFocus bool
}

type INPUT_RECORD struct {
	EventType uint16
	Event     [16]byte
}

var (
	modkernel32           = syscall.NewLazyDLL("kernel32.dll")
	procGetStdHandle      = modkernel32.NewProc("GetStdHandle")
	procGetConsoleMode    = modkernel32.NewProc("GetConsoleMode")
	procSetConsoleMode    = modkernel32.NewProc("SetConsoleMode")
	procReadConsoleW      = modkernel32.NewProc("ReadConsoleW")
	procWriteConsoleW     = modkernel32.NewProc("WriteConsoleW")
	procGetAsyncKeyState  = modkernel32.NewProc("GetAsyncKeyState")
	procPeekConsoleInputW = modkernel32.NewProc("PeekConsoleInputW")
	procReadConsoleInputW = modkernel32.NewProc("ReadConsoleInputW")
)

func GetStdHandle(nStdHandle int) uintptr {
	ret, _, _ := procGetStdHandle.Call(uintptr(nStdHandle))
	return ret
}

func GetConsoleMode(hConsoleHandle uintptr, lpMode *uint32) bool {
	ret, _, _ := procGetConsoleMode.Call(hConsoleHandle, uintptr(unsafe.Pointer(lpMode)))
	return ret != 0
}

func SetConsoleMode(hConsoleHandle uintptr, dwMode uint32) bool {
	ret, _, _ := procSetConsoleMode.Call(hConsoleHandle, uintptr(dwMode))
	return ret != 0
}

func NewInputHandler() (*InputHandler, error) {
	handler := &InputHandler{}

	stdinHandle := GetStdHandle(STD_INPUT_HANDLE)

	var mode uint32
	if !GetConsoleMode(stdinHandle, &mode) {
		return nil, fmt.Errorf("failed to get console mode")
	}

	handler.oldState = &termState{mode: mode}

	newMode := mode &^ (ENABLE_ECHO_INPUT | ENABLE_LINE_INPUT | ENABLE_MOUSE_INPUT)
	newMode |= (ENABLE_WINDOW_INPUT | ENABLE_PROCESSED_INPUT)

	if !SetConsoleMode(stdinHandle, newMode) {
		return nil, fmt.Errorf("failed to set console mode")
	}

	return handler, nil
}

func (h *InputHandler) Restore() {
	if h.oldState != nil {
		stdinHandle := GetStdHandle(STD_INPUT_HANDLE)
		SetConsoleMode(stdinHandle, h.oldState.mode)
	}
}

func (h *InputHandler) ReadKey() Key {
	var input INPUT_RECORD
	var eventsRead uint32
	stdinHandle := GetStdHandle(STD_INPUT_HANDLE)

	for {
		ret, _, _ := procReadConsoleInputW.Call(
			stdinHandle,
			uintptr(unsafe.Pointer(&input)),
			1,
			uintptr(unsafe.Pointer(&eventsRead)),
		)

		if ret == 0 {
			return KeyNull
		}

		if eventsRead == 0 {
			continue
		}

		if input.EventType == 1 {
			keyEvent := (*KEY_EVENT_RECORD)(unsafe.Pointer(&input.Event[0]))

			if keyEvent.bKeyDown == 1 {
				vkCode := keyEvent.wVirtualKeyCode

				if vkCode == 0x51 {
					return KeyCtrlQ
				}

				if vkCode == 0x50 {
					return Key('P')
				}

				if vkCode == 0x4D {
					return Key('M')
				}

				if vkCode == 0x4C {
					return Key('L')
				}

				if vkCode == 0x51 {
					return Key('Q')
				}

				if vkCode == 0x26 {
					return KeyArrowUp
				}

				if vkCode == 0x28 {
					return KeyArrowDown
				}

				if vkCode == 0x25 {
					return KeyArrowLeft
				}

				if vkCode == 0x27 {
					return KeyArrowRight
				}

				if vkCode == 0x70 {
					return KeyF1
				}

				if vkCode == 0x71 {
					return KeyF2
				}

				return Key(vkCode)
			}
		}
	}
}

func (h *InputHandler) NonBlockingRead() (Key, bool) {
	var input INPUT_RECORD
	var eventsRead uint32
	stdinHandle := GetStdHandle(STD_INPUT_HANDLE)

	ret, _, _ := procPeekConsoleInputW.Call(
		stdinHandle,
		uintptr(unsafe.Pointer(&input)),
		1,
		uintptr(unsafe.Pointer(&eventsRead)),
	)

	if ret == 0 || eventsRead == 0 {
		return KeyNull, false
	}

	key := h.ReadKey()
	return key, true
}

func GetKeyPressed() Key {
	for {
		var input INPUT_RECORD
		var eventsRead uint32
		stdinHandle := GetStdHandle(STD_INPUT_HANDLE)

		ret, _, _ := procReadConsoleInputW.Call(
			stdinHandle,
			uintptr(unsafe.Pointer(&input)),
			1,
			uintptr(unsafe.Pointer(&eventsRead)),
		)

		if ret == 0 || eventsRead == 0 {
			continue
		}

		if input.EventType == 1 {
			keyEvent := (*KEY_EVENT_RECORD)(unsafe.Pointer(&input.Event[0]))

			if keyEvent.bKeyDown == 1 {
				vkCode := keyEvent.wVirtualKeyCode

				if vkCode == 0x51 {
					return KeyCtrlQ
				}

				if vkCode == 0x50 {
					return Key('P')
				}

				if vkCode == 0x4D {
					return Key('M')
				}

				if vkCode == 0x4C {
					return Key('L')
				}

				if vkCode == 0x26 {
					return KeyArrowUp
				}

				if vkCode == 0x28 {
					return KeyArrowDown
				}

				if vkCode == 0x70 {
					return KeyF1
				}

				if vkCode == 0x71 {
					return KeyF2
				}

				return Key(vkCode)
			}
		}
	}
}

func SetupSignalHandler() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-c
		fmt.Println("\n收到退出信号，正在退出...")
		os.Exit(0)
	}()
}

func IsKeyPressed(vk int) bool {
	ret, _, _ := procGetAsyncKeyState.Call(uintptr(vk))
	return (ret & 0x8000) != 0
}
