package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/driver/desktop"
)

func main() {
	a := app.New()
	w := a.NewWindow("gchat")

	if desk, ok := a.(desktop.App); ok {
		menu := fyne.NewMenu("gchat", nil)
		desk.SetSystemTrayMenu(m)
	}

	w.Resize(fyne.NewSize(200, 100))
	w.ShowAndRun()
}
