package gui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"github.com/SinfulDragon/dice-calculator/internal/gui/models"
	"github.com/SinfulDragon/dice-calculator/internal/gui/views"
)

type App struct {
	history []models.RollResult
	tabs    *container.AppTabs
	calc    *views.CalculatorView
	hist    *views.HistoryView
	stats   *views.StatsView
	window  fyne.Window
}

func NewApp(a fyne.App) fyne.Window {
	app := &App{}
	app.window = a.NewWindow("Dice Calculator")

	app.calc = views.NewCalculatorView(func(r models.RollResult) {
		app.history = append(app.history, r)
		app.hist.Refresh()
	}, app.window)

	app.hist = views.NewHistoryView(
		&app.history,
		func(formula string) {
			app.calc.SetFormula(formula)
			app.tabs.SelectIndex(0)
		},
		func(formula string) {
			app.stats.SetFormula(formula)
			app.tabs.SelectIndex(2)
		},
		func() {
			app.history = nil
			app.hist.Refresh()
		},
	)

	app.stats = views.NewStatsView(app.window)

	app.tabs = container.NewAppTabs(
		container.NewTabItem("Calculator", app.calc.CanvasObject()),
		container.NewTabItem("History", app.hist.CanvasObject()),
		container.NewTabItem("Stats", app.stats.CanvasObject()),
	)

	app.window.SetContent(app.tabs)
	app.window.Resize(fyne.NewSize(1024, 768))
	app.window.CenterOnScreen()
	return app.window
}

func Show() {
	a := app.New()
	w := NewApp(a)
	w.ShowAndRun()
}
