package main

import (
	"fmt"

	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/widget"
	"github.com/SinfulDragon/dice-calculator/internal/core/parser"
)

func main() {
	Init()
	a := app.New()
	w := a.NewWindow("Dice Calculator")
	label := widget.NewLabel("Hello world!")

	w.SetContent(label)

	formula := "2d12 + 1d6.reroll(rerollexact, values:[1,2])"

	tree, err := parser.ParseFormula(formula)

	if err != nil {
		fmt.Println(err)
		return
	}

	tree.Roll()
	fmt.Println(tree.Evaluate())
	w.ShowAndRun()
}
