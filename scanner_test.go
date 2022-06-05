package editer

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	"golang.org/x/sys/windows"
)

func TestScanMemory(t *testing.T) {
	const da = windows.PROCESS_VM_READ | windows.PROCESS_VM_WRITE | windows.PROCESS_QUERY_INFORMATION

	var hProcess windows.Handle
	processes, err := GetProcessList()
	require.NoError(t, err)
	for i := 0; i < len(processes); i++ {
		if processes[i].Name == "notepad.exe" {
			hProcess, err = OpenProcess(da, false, processes[i].PID)
			require.NoError(t, err)
		}
	}

	value := []byte{0x31, 0x00, 0x32, 0x00, 0x33, 0x00, 0x34, 0x00, 0x35, 0x00, 0x36, 0x00, 0x37, 0x00, 0x38, 0x00, 0x39, 0x00, 0x30, 0x00, 0x31, 0x00, 0x32, 0x00}
	result, err := ScanMemory(hProcess, 0x00000000, 0x7FFFFFFE0000, value, nil)
	require.NoError(t, err)

	for i := 0; i < len(result); i++ {
		fmt.Printf("0x%X  %v\n", result[i].Address, result[i].Value)

		data := bytes.Repeat([]byte{'F', 0x00}, 12)

		_, err = WriteProcessMemory(hProcess, result[i].Address, data)
		require.NoError(t, err)
	}

}
