package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"syscall"

	"golang.org/x/sys/windows"
)

const PROCESS_QUERY_LIMITED_INFORMATION = 0x1000

func getParentProcessName() (string, error) {
	ppid := uint32(os.Getppid())

	hProcess, err := windows.OpenProcess(PROCESS_QUERY_LIMITED_INFORMATION, false, ppid)
	if err != nil {
		return "", err
	}

	defer windows.CloseHandle(hProcess)

	buffer := make([]uint16, windows.MAX_PATH)
	size := uint32(len(buffer))

	err = windows.QueryFullProcessImageName(hProcess, 0, &buffer[0], &size)
	if err != nil {
		return "", err
	}

	processName := syscall.UTF16ToString(buffer)

	return filepath.Base(processName), nil
}

func shell() string {
	parent, err := getParentProcessName()
	if err != nil {
		return "cmd.exe"
	}

	path, err := exec.LookPath(parent)
	if err != nil {
		return "cmd.exe"
	}

	return path
}

func shellWithWT() string {
	sh := shell()

	if wt, err := exec.LookPath("wt"); err == nil {
		sh = wt + " " + sh
	}

	return sh
}
