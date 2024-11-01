package main

import (
	"fmt"
	"unsafe"

	"golang.org/x/sys/windows"
)

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
