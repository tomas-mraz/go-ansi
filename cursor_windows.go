//go:build windows

package ansi

import (
	"os"
	"syscall"
	"unsafe"
)

func CursorUp(n int) {
	cursorMove(0, -n)
}

func CursorDown(n int) {
	cursorMove(0, n)
}

func CursorForward(n int) {
	cursorMove(n, 0)
}

func CursorBack(n int) {
	cursorMove(-n, 0)
}

func cursorMove(x int, y int) {
	handle := syscall.Handle(os.Stdout.Fd())

	var csbi consoleScreenBufferInfo
	procGetConsoleScreenBufferInfo.Call(uintptr(handle), uintptr(unsafe.Pointer(&csbi)))

	var cursor coord
	cursor.x = csbi.cursorPosition.x + Short(x)
	cursor.y = csbi.cursorPosition.y + Short(y)

	procSetConsoleCursorPosition.Call(uintptr(handle), uintptr(*(*int32)(unsafe.Pointer(&cursor))))
}

func CursorNextLine(n int) {
	CursorUp(n)
	CursorHorizontalAbsolute(0)
}

func CursorPreviousLine(n int) {
	CursorDown(n)
	CursorHorizontalAbsolute(0)
}

func LastLine() {
	screen := getScreen()
	CursorAbsolute(0, int(screen.window.bottom))
}

func BottomOffset() int {
	screen := getScreen()
	return int(screen.cursorPosition.y - screen.window.bottom)
}

func IsLastLine() bool {
	screen := getScreen()
	if screen.window.bottom == screen.cursorPosition.y {
		return true
	}
	return false
}

func CursorAbsolute(x int, y int) {
	handle := syscall.Handle(os.Stdout.Fd())

	var csbi consoleScreenBufferInfo
	procGetConsoleScreenBufferInfo.Call(uintptr(handle), uintptr(unsafe.Pointer(&csbi)))

	var cursor coord
	cursor.x = Short(x)
	cursor.y = Short(y)

	if csbi.size.x < cursor.x {
		cursor.x = csbi.size.x
	}
	if csbi.size.y < cursor.y {
		cursor.y = csbi.size.y
	}

	procSetConsoleCursorPosition.Call(uintptr(handle), uintptr(*(*int32)(unsafe.Pointer(&cursor))))
}

func CursorHorizontalAbsolute(x int) {
	handle := syscall.Handle(os.Stdout.Fd())

	var csbi consoleScreenBufferInfo
	procGetConsoleScreenBufferInfo.Call(uintptr(handle), uintptr(unsafe.Pointer(&csbi)))

	var cursor coord
	cursor.x = Short(x)
	cursor.y = csbi.cursorPosition.y

	if csbi.size.x < cursor.x {
		cursor.x = csbi.size.x
	}

	procSetConsoleCursorPosition.Call(uintptr(handle), uintptr(*(*int32)(unsafe.Pointer(&cursor))))
}

func GetCursorPosition() (int, int) {
	screen := getScreen()
	return int(screen.cursorPosition.x), int(screen.cursorPosition.y)
}

func CursorShow() {
	handle := syscall.Handle(os.Stdout.Fd())

	var cci consoleCursorInfo
	procGetConsoleCursorInfo.Call(uintptr(handle), uintptr(unsafe.Pointer(&cci)))
	cci.visible = 1

	procSetConsoleCursorInfo.Call(uintptr(handle), uintptr(unsafe.Pointer(&cci)))
}

func CursorHide() {
	handle := syscall.Handle(os.Stdout.Fd())

	var cci consoleCursorInfo
	procGetConsoleCursorInfo.Call(uintptr(handle), uintptr(unsafe.Pointer(&cci)))
	cci.visible = 0

	procSetConsoleCursorInfo.Call(uintptr(handle), uintptr(unsafe.Pointer(&cci)))
}
