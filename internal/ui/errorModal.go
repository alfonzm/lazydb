package ui

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type ErrorModal struct {
	alertModal     tview.Primitive
	alertContainer *tview.Flex
	app            *App
	lastFocus      tview.Primitive
}

func NewErrorModal() (*ErrorModal, error) {
	alertModal := tview.NewBox()
	alertFlex := tview.NewFlex()

	errorModal := &ErrorModal{
		alertModal:     alertModal,
		alertContainer: alertFlex,
	}

	errorModal.setKeyBindings()

	return errorModal, nil
}

func (e *ErrorModal) RenderError(errorText string) {
	/* https://github.com/rivo/tview/wiki/Modal */
	// Create modal container centered on screen
	modal := func(p tview.Primitive, width, height int) tview.Primitive {
		return tview.NewFlex().
			AddItem(nil, 0, 1, false).
			AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
				AddItem(nil, 0, 1, false).
				AddItem(p, height, 1, true).
				AddItem(nil, 0, 1, false), width, 1, true).
			AddItem(nil, 0, 1, false)
	}

	// Modal text
	text := tview.NewTextView().
		SetText(errorText).
		SetTextColor(tcell.ColorRed)

	// Instructions at the bottom most of the modal
	// saying press Y to copy the error, ESC or Q to close modal
	legend := tview.NewTextView().
		SetText("[Y] Copy error / [Esc] Close").
		SetTextColor(tcell.ColorYellow).
		SetTextAlign(tview.AlignCenter)

	// Modal body
	alertFlex := tview.NewFlex()
	alertFlex.SetBorder(true).
		SetTitle("ERROR")
		// SetBorderColor(tcell.ColorRed).
		// SetTitleColor(tcell.ColorRed)
	alertFlex.SetDirection(tview.FlexRow)

	alertFlex.AddItem(text, 0, 1, false)
	alertFlex.AddItem(legend, 1, 1, false)

	// Modal
	e.alertModal = modal(alertFlex, 100, 15)
	e.alertContainer = alertFlex

	e.setKeyBindings()

	e.lastFocus = e.app.GetFocus()

	e.app.appPages.RemovePage("modal")
	e.app.appPages.AddPage("modal", e.app.errorModal.alertModal, true, false)
	e.app.appPages.ShowPage("modal")
	e.app.SetFocus(e.alertContainer)
}

func (e *ErrorModal) setKeyBindings() {
	e.alertContainer.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyEscape:
			e.app.appPages.RemovePage("modal")
			e.app.appPages.SwitchToPage("app")
			e.app.SetFocus(e.lastFocus)
		}
		return event
	})
}
