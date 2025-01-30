package main

import (
	"fmt"
	"os"
	"strings"
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows"
	"golang.org/x/term"
)

var (
	user32                       = syscall.NewLazyDLL("user32.dll")
	procGetForegroundWindow      = user32.NewProc("GetForegroundWindow")
	procGetWindowThreadProcessId = user32.NewProc("GetWindowThreadProcessId")
	procGetWindowRect            = user32.NewProc("GetWindowRect")
)

type RECT struct {
	Left, Top, Right, Bottom int32
}

type tokenElevation struct {
	TokenIsElevated uint32
}

func isElevated() (bool, error) {
	var token windows.Token

	err := windows.OpenProcessToken(windows.CurrentProcess(), windows.TOKEN_QUERY, &token)
	if err != nil {
		return false, fmt.Errorf("OpenProcessToken error: %v", err)
	}

	defer token.Close()

	var (
		elevation     tokenElevation
		elevationSize uint32
	)

	err = windows.GetTokenInformation(
		token,
		windows.TokenElevation,
		(*byte)(unsafe.Pointer(&elevation)),
		uint32(unsafe.Sizeof(elevation)),
		&elevationSize,
	)

	if err != nil {
		return false, fmt.Errorf("GetTokenInformation error: %v", err)
	}

	return elevation.TokenIsElevated != 0, nil
}

func terminalSize() (int, int) {
	width, height, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		return 0, 0
	}

	return width, height
}

func processExists(name string) (bool, *RECT) {
	handle, err := windows.CreateToolhelp32Snapshot(windows.TH32CS_SNAPPROCESS, 0)
	if err != nil {
		return false, nil
	}

	defer windows.CloseHandle(handle)

	var entry windows.ProcessEntry32

	entry.Size = uint32(unsafe.Sizeof(entry))

	if err := windows.Process32First(handle, &entry); err != nil {
		return false, nil
	}

	for {
		if strings.EqualFold(windows.UTF16ToString(entry.ExeFile[:]), name) {
			hwnd, _, _ := procGetForegroundWindow.Call()

			var pid uint32

			procGetWindowThreadProcessId.Call(hwnd, uintptr(unsafe.Pointer(&pid)))

			if pid != entry.ProcessID {
				return true, nil
			}

			rect := &RECT{}

			procGetWindowRect.Call(hwnd, uintptr(unsafe.Pointer(rect)))

			return true, rect
		}

		if err := windows.Process32Next(handle, &entry); err != nil {
			break
		}
	}

	return false, nil
}

func terminateParent() {
	pid := os.Getppid()

	hProcess, err := windows.OpenProcess(windows.PROCESS_TERMINATE, false, uint32(pid))
	if err != nil {
		return
	}

	defer windows.CloseHandle(hProcess)

	windows.TerminateProcess(hProcess, uint32(1))
}
