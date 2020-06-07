package main

import (
	"fyne.io/fyne"
	"fyne.io/fyne/widget"
)

// FormField ...
type FormField struct {
	*widget.Box

	Entry fyne.CanvasObject
}

// NewFormField ...
func NewFormField(label string, entry fyne.CanvasObject) *FormField {
	b := widget.NewHBox(
		widget.NewVBox(widget.NewLabel(label)),
		widget.NewVBox(entry),
	)

	return &FormField{
		Box:   b,
		Entry: entry,
	}
}
