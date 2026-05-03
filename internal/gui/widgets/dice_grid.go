package widgets

import (
	"fmt"
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/SinfulDragon/dice-calculator/internal/core/common"
)

func dieColor(sides int) color.Color {
	switch sides {
	case 4:
		return color.NRGBA{R: 255, G: 0, B: 0, A: 255}
	case 6:
		return color.White
	case 8:
		return color.NRGBA{R: 0, G: 255, B: 0, A: 255}
	case 10:
		return color.NRGBA{R: 0, G: 0, B: 255, A: 255}
	case 12:
		return color.NRGBA{R: 128, G: 0, B: 128, A: 255}
	case 20:
		return color.NRGBA{R: 255, G: 165, B: 0, A: 255}
	default:
		return color.Gray{Y: 128}
	}
}

func NewDiceGrid(dice []*common.Die) fyne.CanvasObject {
	if len(dice) == 0 {
		return widget.NewLabel("No dice")
	}
	items := make([]fyne.CanvasObject, len(dice))
	for i, d := range dice {
		bg := canvas.NewRectangle(dieColor(d.Sides))
		bg.StrokeWidth = 2
		bg.StrokeColor = color.Black

		txt := canvas.NewText(fmt.Sprintf("%d", d.Value), color.Black)
		txt.TextSize = 18
		txt.TextStyle = fyne.TextStyle{Bold: true}

		items[i] = container.NewStack(bg, container.NewCenter(txt))
	}
	grid := container.NewGridWrap(fyne.NewSize(54, 54), items...)
	scroll := container.NewScroll(grid)
	scroll.SetMinSize(fyne.NewSize(540, 120))
	return scroll
}
