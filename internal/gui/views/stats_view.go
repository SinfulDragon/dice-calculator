package views

import (
	"fmt"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"github.com/SinfulDragon/dice-calculator/internal/core/parser"
	"github.com/SinfulDragon/dice-calculator/internal/core/stats"
)

type StatsView struct {
	container       *fyne.Container
	formulaEntry    *widget.Entry
	iterationsEntry *widget.Entry
	resultArea      *fyne.Container
	window          fyne.Window
}

func NewStatsView(window fyne.Window) *StatsView {
	sv := &StatsView{
		window: window,
	}

	sv.formulaEntry = widget.NewEntry()
	sv.formulaEntry.SetPlaceHolder("Enter formula")

	sv.iterationsEntry = widget.NewEntry()
	sv.iterationsEntry.SetText("100000")
	sv.iterationsEntry.SetPlaceHolder("Iterations")

	analyzeBtn := widget.NewButton("Analyze", func() {
		sv.analyze()
	})
	analyzeBtn.Importance = widget.HighImportance

	formulaRow := container.NewBorder(nil, nil, nil, container.NewHBox(analyzeBtn), sv.formulaEntry)
	iterRow := container.NewHBox(widget.NewLabel("Iterations"), sv.iterationsEntry)

	top := container.NewVBox(formulaRow, iterRow)

	sv.resultArea = container.NewVBox()
	sv.resultArea.Hide()

	sv.container = container.NewVBox(top, sv.resultArea)
	return sv
}

func (sv *StatsView) analyze() {
	formula := sv.formulaEntry.Text
	if formula == "" {
		dialog.ShowError(fmt.Errorf("formula is empty"), sv.window)
		return
	}
	node, err := parser.ParseFormula(formula)
	if err != nil {
		dialog.ShowError(err, sv.window)
		return
	}

	iterStr := sv.iterationsEntry.Text
	if iterStr == "" {
		iterStr = "100000"
	}
	iterations, err := strconv.Atoi(iterStr)
	if err != nil || iterations <= 0 {
		dialog.ShowError(fmt.Errorf("invalid iterations: %s", iterStr), sv.window)
		return
	}

	analyzer := stats.NewAnalyzer(node)
	dist, err := analyzer.MonteCarlo(iterations)
	if err != nil {
		dialog.ShowError(err, sv.window)
		return
	}

	summary := dist.Summary()
	summaryBox := container.NewVBox(
		widget.NewLabel(fmt.Sprintf("Min: %.0f", summary.Min)),
		widget.NewLabel(fmt.Sprintf("Max: %.0f", summary.Max)),
		widget.NewLabel(fmt.Sprintf("Mean: %.2f", summary.Mean)),
		widget.NewLabel(fmt.Sprintf("StdDev: %.2f", summary.StdDev)),
	)

	outcomes := dist.Outcomes()
	maxProb := 0.0
	for _, o := range outcomes {
		if p := dist.Probability(o); p > maxProb {
			maxProb = p
		}
	}

	histBox := container.NewVBox()
	for _, o := range outcomes {
		prob := dist.Probability(o)
		bar := widget.NewProgressBar()
		if maxProb > 0 {
			bar.Max = maxProb
		}
		bar.Value = prob

		row := container.NewBorder(nil, nil,
			widget.NewLabel(fmt.Sprintf("%d", o)),
			widget.NewLabel(fmt.Sprintf("%.2f%%", prob*100)),
			bar,
		)
		histBox.Add(row)
	}

	scroll := container.NewScroll(histBox)
	scroll.SetMinSize(fyne.NewSize(400, 300))

	sv.resultArea.Objects = []fyne.CanvasObject{
		widget.NewLabelWithStyle("Summary", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		summaryBox,
		widget.NewLabelWithStyle("Distribution", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		scroll,
	}
	sv.resultArea.Show()
	sv.resultArea.Refresh()
}

func (sv *StatsView) SetFormula(formula string) {
	sv.formulaEntry.SetText(formula)
	sv.resultArea.Hide()
	sv.resultArea.Objects = nil
	sv.resultArea.Refresh()
}

func (sv *StatsView) CanvasObject() fyne.CanvasObject {
	return sv.container
}
