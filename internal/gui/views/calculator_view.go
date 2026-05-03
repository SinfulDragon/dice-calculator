package views

import (
	"fmt"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"github.com/SinfulDragon/dice-calculator/internal/core/common"
	"github.com/SinfulDragon/dice-calculator/internal/core/parser"
	"github.com/SinfulDragon/dice-calculator/internal/core/tree"
	"github.com/SinfulDragon/dice-calculator/internal/gui/models"
	"github.com/SinfulDragon/dice-calculator/internal/gui/widgets"
)

type CalculatorView struct {
	container     *fyne.Container
	formulaEntry  *widget.Entry
	resultArea    *fyne.Container
	builderScroll *container.Scroll
	modeRadio     *widget.RadioGroup

	currentNode  tree.FormulaNode
	currentDice  []*common.Die
	currentTotal int

	onAddHistory func(models.RollResult)
	window       fyne.Window
}

func NewCalculatorView(onAddHistory func(models.RollResult), window fyne.Window) *CalculatorView {
	cv := &CalculatorView{
		onAddHistory: onAddHistory,
		window:       window,
	}

	cv.formulaEntry = widget.NewEntry()
	cv.formulaEntry.SetPlaceHolder("Enter formula (e.g. 2d6 + 3)")

	rollBtn := widget.NewButton("Roll", func() {
		cv.performRoll()
	})
	rollBtn.Importance = widget.HighImportance
	clearBtn := widget.NewButton("Clear", func() {
		cv.formulaEntry.SetText("")
		cv.hideResult()
	})

	btnBox := container.NewHBox(rollBtn, clearBtn)
	top := container.NewBorder(nil, nil, nil, btnBox, cv.formulaEntry)

	cv.resultArea = container.NewVBox()
	cv.resultArea.Hide()

	cv.builderScroll = container.NewScroll(widget.NewLabel(""))
	cv.builderScroll.Hide()

	cv.modeRadio = widget.NewRadioGroup([]string{"Text", "Builder"}, func(selected string) {
		cv.switchMode(selected)
	})
	cv.modeRadio.Horizontal = true
	cv.modeRadio.Required = true
	cv.modeRadio.Selected = "Text"

	// Center area: modeRadio on top, builderScroll fills the rest
	center := container.NewBorder(cv.modeRadio, nil, nil, nil, cv.builderScroll)

	cv.container = container.NewBorder(top, cv.resultArea, nil, nil, center)
	return cv
}

func (cv *CalculatorView) performRoll() {
	formula := cv.formulaEntry.Text
	if formula == "" {
		return
	}
	node, err := parser.ParseFormula(formula)
	if err != nil {
		dialog.ShowError(err, cv.window)
		return
	}
	cv.currentNode = node
	cv.currentDice = node.Roll()
	cv.currentTotal = node.Evaluate()

	cv.showResult()
}

func (cv *CalculatorView) showResult() {
	cv.resultArea.Objects = nil

	formulaLabel := widget.NewLabelWithStyle(cv.formulaEntry.Text, fyne.TextAlignLeading, fyne.TextStyle{Bold: true})
	totalLabel := widget.NewLabelWithStyle(
		fmt.Sprintf("Total: %d", cv.currentTotal),
		fyne.TextAlignLeading,
		fyne.TextStyle{Bold: true},
	)

	diceGrid := widgets.NewDiceGrid(cv.currentDice)

	addHistoryBtn := widget.NewButton("Add to History", func() {
		if cv.onAddHistory != nil {
			cv.onAddHistory(models.RollResult{
				FormulaStr: cv.formulaEntry.Text,
				Node:       cv.currentNode,
				Dice:       cv.currentDice,
				Total:      cv.currentTotal,
				Time:       time.Now(),
			})
		}
	})

	cv.resultArea.Add(container.NewVBox(formulaLabel, diceGrid, totalLabel, addHistoryBtn))
	cv.resultArea.Show()
	cv.resultArea.Refresh()
}

func (cv *CalculatorView) hideResult() {
	cv.resultArea.Hide()
	cv.resultArea.Objects = nil
	cv.resultArea.Refresh()
}

func (cv *CalculatorView) switchMode(selected string) {
	if selected == "Builder" {
		formula := cv.formulaEntry.Text
		if formula == "" {
			formula = "0"
		}
		node, err := parser.ParseFormula(formula)
		if err != nil {
			dialog.ShowError(err, cv.window)
			cv.modeRadio.SetSelected("Text")
			return
		}
		cv.formulaEntry.Disable()
		bp := widgets.NewBuilderPanel(node, func(updated tree.FormulaNode) {
			cv.formulaEntry.SetText(widgets.NodeString(updated))
		})
		cv.builderScroll.Content = bp.CanvasObject()
		cv.builderScroll.Show()
		cv.builderScroll.Refresh()
	} else {
		cv.formulaEntry.Enable()
		cv.builderScroll.Hide()
	}
	cv.container.Refresh()
}

func (cv *CalculatorView) SetFormula(formula string) {
	cv.formulaEntry.SetText(formula)
	cv.modeRadio.SetSelected("Text")
	cv.formulaEntry.Enable()
	cv.builderScroll.Hide()
	cv.hideResult()
}

func (cv *CalculatorView) CanvasObject() fyne.CanvasObject {
	return cv.container
}
