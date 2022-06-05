//go:build windows
// +build windows

package editer

import (
	"golang.org/x/sys/windows"
)

var (
	modKernel32 = windows.NewLazySystemDLL("kernel32.dll")

	procReadProcessMemory  = modKernel32.NewProc("ReadProcessMemory")
	procWriteProcessMemory = modKernel32.NewProc("WriteProcessMemory")
	procVirtualQueryEx     = modKernel32.NewProc("VirtualQueryEx")
)

// CloseHandle is used to close handle, it will not return error.
func CloseHandle(handle windows.Handle) {
	_ = windows.CloseHandle(handle)
}
