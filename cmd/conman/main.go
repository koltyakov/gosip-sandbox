package main

import (
	"os"

	"fyne.io/fyne"
	"fyne.io/fyne/app"
	"fyne.io/fyne/widget"
)

func main() {
	os.Setenv("FYNE_THEME", "light")
	a := app.New()

	w := a.NewWindow("SharePoint Connection Manager")
	w.CenterOnScreen()
	w.SetFixedSize(true)
	w.Resize(fyne.NewSize(400, 200))

	c := widget.NewGroup(
		"Connection details",
		NewAuthForm(),
	)

	w.SetContent(c)
	w.ShowAndRun()
}
