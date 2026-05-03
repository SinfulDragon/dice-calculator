package views

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/SinfulDragon/dice-calculator/internal/gui/models"
)

type HistoryView struct {
	container       *fyne.Container
	history         *[]models.RollResult
	onLoad          func(string)
	onStats         func(string)
	scrollContainer *fyne.Container
}

func NewHistoryView(history *[]models.RollResult, onLoad, onStats func(string), onClear func()) *HistoryView {
	hv := &HistoryView{
		history: history,
		onLoad:  onLoad,
		onStats: onStats,
	}

	hv.scrollContainer = container.NewVBox()
	scroll := container.NewScroll(hv.scrollContainer)
	scroll.SetMinSize(fyne.NewSize(400, 300))

	clearBtn := widget.NewButton("Clear History", func() {
		if onClear != nil {
			onClear()
		}
	})
	clearBtn.Importance = widget.DangerImportance

	hv.container = container.NewBorder(nil, clearBtn, nil, nil, scroll)
	hv.refreshList()
	return hv
}

func (hv *HistoryView) refreshList() {
	hv.scrollContainer.Objects = nil
	if len(*hv.history) == 0 {
		hv.scrollContainer.Add(widget.NewLabelWithStyle("No history yet", fyne.TextAlignCenter, fyne.TextStyle{Italic: true}))
	} else {
		for i := len(*hv.history) - 1; i >= 0; i-- {
			res := (*hv.history)[i]
			lbl := widget.NewLabel(res.FormulaStr)
			lbl.TextStyle = fyne.TextStyle{Bold: true}
			totalLbl := widget.NewLabel(fmt.Sprintf("= %d  (%s)", res.Total, res.Time.Format("15:04:05")))
			formula := res.FormulaStr // capture for closure
			loadBtn := widget.NewButton("Load", func() {
				if hv.onLoad != nil {
					hv.onLoad(formula)
				}
			})
			statsBtn := widget.NewButton("Stats", func() {
				if hv.onStats != nil {
					hv.onStats(formula)
				}
			})
			row := container.NewHBox(lbl, totalLbl, layout.NewSpacer(), loadBtn, statsBtn)
			hv.scrollContainer.Add(row)
			hv.scrollContainer.Add(widget.NewSeparator())
		}
	}
	hv.scrollContainer.Refresh()
}

func (hv *HistoryView) Refresh() {
	hv.refreshList()
}

func (hv *HistoryView) CanvasObject() fyne.CanvasObject {
	return hv.container
}
