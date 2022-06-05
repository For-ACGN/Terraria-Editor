package editer

import (
	"bytes"
	"syscall"

	"golang.org/x/sys/windows"
)

var pageSize = uintptr(syscall.Getpagesize())

// ScanOptions contains scan options.
type ScanOptions struct {
	Writable     bool // TODO wait
	Executable   bool // TODO wait
	PauseProcess bool // TODO wait
}

// ScanResult contains scan result.
type ScanResult struct {
	Address uintptr
	Value   []byte
}

// ScanMemory is used to scan value in the target process memory.
func ScanMemory(hProcess windows.Handle, start, stop uintptr, value []byte, opts *ScanOptions) ([]*ScanResult, error) {
	if opts == nil {
		opts = new(ScanOptions) // not used
	}
	var (
		mbi MemoryBasicInformation
		err error
	)

	sr := make([]*ScanResult, 0, 1024)

	buf := make([]byte, pageSize)
	for addr := start; addr < stop; addr += mbi.RegionSize {
		err = VirtualQueryEx(hProcess, addr, &mbi)
		if err != nil {
			return nil, err
		}

		// fmt.Println(mbi.RegionSize)

		if mbi.RegionSize == 0 {
			break
		}
		// // MEM_COMMIT 0x1000
		if mbi.State != 0x1000 {
			continue
		}
		// PAGE_EXECUTE_READWRITE 0x40
		// PAGE_EXECUTE_WRITECOPY 0x80
		// PAGE_READWRITE         0x04
		// PAGE_WRITECOPY         0x08
		if mbi.Protect != 0x04 { // TODO update it
			continue
		}

		for block := uintptr(0); block < mbi.RegionSize/pageSize; block++ {
			startAddr := addr + block*pageSize
			n, err := ReadProcessMemory(hProcess, startAddr, &buf[0], pageSize)
			if err != nil {
				return nil, err
			}
			result := sundaySearch(startAddr, buf[:n], value)
			if len(result) > 0 {
				sr = append(sr, result...)
			}
		}
	}
	return sr, nil
}

// TODO use new algo, sundaySearch
func sundaySearch(addr uintptr, buf, value []byte) []*ScanResult {
	result := make([]*ScanResult, 0, 16)
	l := len(value)
	var (
		skipped int
		idx     int
	)

	for {
		idx = bytes.Index(buf, value)
		if idx == -1 {
			break
		}
		cp := make([]byte, l)
		copy(cp, buf[idx:idx+l])
		buf = buf[idx+l:]

		result = append(result, &ScanResult{
			Address: addr + uintptr(skipped+idx),
			Value:   cp,
		})

		skipped += idx + l
	}
	return result
}
