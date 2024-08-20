//go:build windows

package ansi

import (
	"golang.org/x/term"
	"os"
	"syscall"
	"unsafe"
)

func GetTerminalSize() (int, int) {
	width, height, _ := term.GetSize(int(os.Stdout.Fd()))
	return width, height
}

func GetScreenSize() (int, int) {
	screen := getScreen()
	return int(screen.size.x), int(screen.size.y)
}

func getScreen() consoleScreenBufferInfo {
	handle := syscall.Handle(os.Stdout.Fd())

	var csbi consoleScreenBufferInfo
	procGetConsoleScreenBufferInfo.Call(uintptr(handle), uintptr(unsafe.Pointer(&csbi)))
	return csbi
}

func EraseInLine(mode int) {
	handle := syscall.Handle(os.Stdout.Fd())

	var csbi consoleScreenBufferInfo
	procGetConsoleScreenBufferInfo.Call(uintptr(handle), uintptr(unsafe.Pointer(&csbi)))

	var w uint32
	var x Short
	cursor := csbi.cursorPosition
	switch mode {
	case 1:
		x = csbi.size.x
	case 2:
		x = 0
	case 3:
		cursor.x = 0
		x = csbi.size.x
	}
	procFillConsoleOutputCharacter.Call(uintptr(handle), uintptr(' '), uintptr(x), uintptr(*(*int32)(unsafe.Pointer(&cursor))), uintptr(unsafe.Pointer(&w)))
}

func GetBottomScrollIndex() int {
	screen := getScreen()
	return int(screen.window.bottom)
}
