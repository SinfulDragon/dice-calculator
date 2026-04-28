package main

import (
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/widget"
)

func main() {
	a := app.New()
	w := a.NewWindow("Dice Calculator")
	label := widget.NewLabel("Hello world!")

	w.SetContent(label)
	w.ShowAndRun()
}
