package main

import (
	"encoding/binary"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"sync/atomic"
	"unsafe"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"golang.org/x/sys/windows"

	"github.com/For-ACGN/Terraria-Editor"
)

var (
	hProcess atomic.Value // windows.Handle
)

func main() {
	App := app.New()

	window := App.NewWindow("Terraria Editor")
	window.Resize(fyne.Size{
		Width:  600,
		Height: 400,
	})
	window.CenterOnScreen()

	debugText := widget.NewMultiLineEntry()
	debugText.Move(fyne.NewPos(10, 10))
	debugText.Resize(fyne.Size{
		Width:  300,
		Height: 300,
	})

	open := widget.NewButton("Open", func() {
		err := openProcess()
		if err != nil {
			dialog.NewError(err, window).Show()
			return
		}
		_, err = scanRoleList(debugText)
		if err != nil {
			dialog.NewError(err, window).Show()
			return
		}
	})
	open.Move(fyne.NewPos(500, 10))
	open.Resize(fyne.Size{
		Width:  80,
		Height: 40,
	})

	cont := container.NewWithoutLayout()
	cont.Add(debugText)
	cont.Add(open)
	window.SetContent(cont)

	window.ShowAndRun()
}

func openProcess() error {
	if hProcess.Load() != nil {
		return nil
	}
	id, err := editor.GetProcessIDByName("Terraria.exe")
	if err != nil {
		return err
	}
	da := windows.PROCESS_VM_READ |
		windows.PROCESS_VM_WRITE |
		windows.PROCESS_QUERY_INFORMATION
	h, err := editor.OpenProcess(uint32(da), false, id[0])
	if err != nil {
		return err
	}
	hProcess.Store(h)
	return nil
}

func scanRoleList(dt *widget.Entry) (interface{}, error) {
	h := hProcess.Load().(windows.Handle)
	// scan runtime data pointer about role structure
	roleRT := []byte{
		0x00, 0x02, 0x00, 0x01, 0x58, 0x0A, 0x00, 0x00,
		0x88, 0x25, 0x40, 0x00, 0x05, 0x00, 0x00, 0x00,
		0xF0, 0x65,
	}
	result, err := editor.ScanMemory(h, 0x00000000, 0x7FFFFFFF, roleRT, nil)
	if err != nil {
		return nil, err
	}
	if len(result) != 1 {
		dt.SetText(strconv.Itoa(len(result)))
		return nil, errors.New("failed to scan runtime data pointer about role structure")
	}
	roleStructPtr := make([]byte, 4)
	binary.LittleEndian.PutUint32(roleStructPtr, uint32(result[0].Address))
	// scan role list
	result, err = editor.ScanMemory(h, 0x00000000, 0x7FFFFFFF, roleStructPtr, nil)
	if err != nil {
		return nil, err
	}
	text := strings.Builder{}
	for i := 0; i < len(result); i++ {
		val := fmt.Sprintf("0x%X  %v\n", result[i].Address, result[i].Value)
		text.WriteString(val)
	}
	dt.SetText(text.String())
	// TODO find self role structure pointer
	selfRoleAddr := uintptr(0x3933AFF4)

	err = readRoleInfo(dt, selfRoleAddr)
	if err != nil {
		return nil, err
	}
	return nil, nil
}

type role struct {
	_ [0xD0]byte

	Item uint32

	_ [0x0310]byte

	HP uint32
	MP uint32
}

type item struct {
}

func readRoleInfo(dt *widget.Entry, addr uintptr) error {
	h := hProcess.Load().(windows.Handle)
	selfRole := role{}

	_, err := editor.ReadProcessMemory(h, addr,
		(*byte)(unsafe.Pointer(&selfRole)), unsafe.Sizeof(selfRole),
	)
	if err != nil {
		return err
	}

	basic := fmt.Sprintf("HP: %d, MP: %d\n", selfRole.HP, selfRole.MP)
	basic += fmt.Sprintf("0x%X\n", selfRole.Item)

	// read role item
	var items [5*10 + 8 + 1]uint32
	_, err = editor.ReadProcessMemory(h, uintptr(selfRole.Item+8),
		(*byte)(unsafe.Pointer(&items[0])), uintptr(len(items)*4),
	)
	if err != nil {
		return err
	}
	var itemAddrs string
	for i := 0; i < len(items); i++ {
		itemAddrs += fmt.Sprintf("0x%X\n", items[i])
	}

	text := strings.Builder{}
	text.WriteString(basic)
	text.WriteString(itemAddrs)
	dt.SetText(text.String())
	return nil
}

// func readRoleItem(dt *widget.Entry, addr uintptr) error {
//
// }
