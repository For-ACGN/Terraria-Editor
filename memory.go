//go:build windows
// +build windows

package editor

import (
	"unsafe"

	"golang.org/x/sys/windows"
)

// reference:
// https://docs.microsoft.com/en-us/windows/win32/api/memoryapi/nf-memoryapi-readprocessmemory
// https://docs.microsoft.com/en-us/windows/win32/api/memoryapi/nf-memoryapi-writeprocessmemory

// ReadProcessMemory is used to read memory from process. // #nosec
func ReadProcessMemory(hProcess windows.Handle, addr uintptr, buf *byte, size uintptr) (int, error) {
	const name = "ReadProcessMemory"
	var n uint
	ret, _, err := procReadProcessMemory.Call(
		uintptr(hProcess), addr,
		uintptr(unsafe.Pointer(buf)), size,
		uintptr(unsafe.Pointer(&n)),
	)
	if ret == 0 {
		return 0, newErrorf(name, err, "failed to read process memory at 0x%X", addr)
	}
	return int(n), nil
}

// WriteProcessMemory is used to write data to memory in process. // #nosec
func WriteProcessMemory(hProcess windows.Handle, addr uintptr, data []byte) (int, error) {
	const name = "WriteProcessMemory"
	var n uint
	ret, _, err := procWriteProcessMemory.Call(
		uintptr(hProcess), addr,
		uintptr(unsafe.Pointer(&data[0])), uintptr(len(data)),
		uintptr(unsafe.Pointer(&n)),
	)
	if ret == 0 {
		return 0, newErrorf(name, err, "failed to write process memory at 0x%X", addr)
	}
	return int(n), nil
}

// MemoryBasicInformation contains a range of pages in the virtual address space of a process.
// The VirtualQuery and VirtualQueryEx functions use this structure.
type MemoryBasicInformation struct {
	BaseAddress       uintptr
	AllocationBase    uintptr
	AllocationProtect uint32
	PartitionID       uint16
	RegionSize        uintptr
	State             uint32
	Protect           uint32
	Type              uint32
}

// VirtualQueryEx is used to retrieve information about a range of pages in the virtual address
// space of the calling process. To retrieve information about a range of pages in the address
// space of another process, use the VirtualQueryEx function. // #nosec
func VirtualQueryEx(hProcess windows.Handle, addr uintptr, mbi *MemoryBasicInformation) error {
	const name = "VirtualQueryEx"
	ret, _, err := procVirtualQueryEx.Call(
		uintptr(hProcess), addr, uintptr(unsafe.Pointer(mbi)), unsafe.Sizeof(MemoryBasicInformation{}),
	)
	if ret == 0 {
		return newErrorf(name, err, "failed to query memory information at 0x%X", addr)
	}
	return nil
}
