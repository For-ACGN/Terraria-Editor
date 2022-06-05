//go:build windows
// +build windows

package editer

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	"golang.org/x/sys/windows"
)

func TestGetProcessList(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		processes, err := GetProcessList()
		require.NoError(t, err)

		fmt.Println("Name    PID    PPID")
		for _, process := range processes {
			fmt.Printf("%s %d %d\n", process.Name, process.PID, process.PPID)
		}
	})
}

func TestGetProcessIDByName(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		pid, err := GetProcessIDByName("svchost.exe")
		require.NoError(t, err)

		require.NotEmpty(t, pid)
		for _, pid := range pid {
			t.Log("pid:", pid)
		}
	})
}

func TestOpenProcess(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		hProcess, err := OpenProcess(windows.PROCESS_QUERY_INFORMATION, false, uint32(os.Getpid()))
		require.NoError(t, err)

		CloseHandle(hProcess)
	})
}
