package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
)

func main() {
	App := app.New()

	window := App.NewWindow("Terraria Editor")
	window.Resize(fyne.Size{
		Width:  600,
		Height: 300,
	})
	window.CenterOnScreen()

	// role info runtime data
	// 00 02 00 01 58 0A 00 00 88 25 40 00 05 00 00 00 F0 65

	window.ShowAndRun()
}
