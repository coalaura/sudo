package main

import (
	"fmt"
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows"
)

const (
	SEE_MASK_NOCLOSEPROCESS = 0x00000040
	SW_NORMAL               = 1
)

type SHELLEXECUTEINFO struct {
	CbSize       uint32
	FMask        uint32
	Hwnd         uintptr
	LpVerb       *uint16
	LpFile       *uint16
	LpParameters *uint16
	LpDirectory  *uint16
	NShow        int32
	HInstApp     uintptr
	LpIDList     uintptr
	LpClass      *uint16
	HkeyClass    uintptr
	DwHotKey     uint32
	HIcon        uintptr
	HProcess     uintptr
}

func ShellExecuteEx(info *SHELLEXECUTEINFO) error {
	r, _, err := syscall.SyscallN(
		windows.NewLazySystemDLL("shell32.dll").NewProc("ShellExecuteExW").Addr(),
		uintptr(unsafe.Pointer(info)),
	)
	if r == 0 {
		return err
	}
	return nil
}

func elevate(app, params, cwd string) error {
	verbPtr, err := windows.UTF16PtrFromString("runas")
	if err != nil {
		return err
	}

	appPtr, err := windows.UTF16PtrFromString(app)
	if err != nil {
		return err
	}

	paramsPtr, err := windows.UTF16PtrFromString(params)
	if err != nil {
		return err
	}

	dirPtr, err := windows.UTF16PtrFromString(cwd)
	if err != nil {
		return err
	}

	execInfo := &SHELLEXECUTEINFO{
		CbSize:       uint32(unsafe.Sizeof(SHELLEXECUTEINFO{})),
		FMask:        SEE_MASK_NOCLOSEPROCESS,
		Hwnd:         0,
		LpVerb:       verbPtr,
		LpFile:       appPtr,
		LpParameters: paramsPtr,
		LpDirectory:  dirPtr,
		NShow:        SW_NORMAL,
	}

	err = ShellExecuteEx(execInfo)
	if err != nil {
		return fmt.Errorf("ShellExecuteEx error: %v", err)
	}

	return nil
}
