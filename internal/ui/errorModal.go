package ui

import (
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"golang.design/x/clipboard"
)

type ErrorModal struct {
	alertModal     tview.Primitive
	alertContainer *tview.Flex
	app            *App
	lastFocus      tview.Primitive
	errorString    string
	errorText      *tview.TextView
}

func NewErrorModal() (*ErrorModal, error) {
	alertModal := tview.NewBox()
	alertContainer := tview.NewFlex()

	errorModal := &ErrorModal{
		alertModal:     alertModal,
		alertContainer: alertContainer,
	}

	errorModal.setKeyBindings()

	return errorModal, nil
}

func (e *ErrorModal) RenderError(errorText string) {
	// Modal text
	e.errorText = tview.NewTextView().
		SetText(errorText).
		SetTextColor(tcell.ColorRed).
		SetDynamicColors(true)

	// Instructions at the bottom most of the modal
	// saying press Y to copy the error, ESC or Q to close modal
	legend := tview.NewTextView().
		SetText("[Y] Copy error / [Esc] Close").
		SetTextColor(tcell.ColorYellow).
		SetTextAlign(tview.AlignCenter)

	// Modal body
	alertContainer := tview.NewFlex()
	alertContainer.SetBorder(true).
		SetTitle("ERROR")
		// SetBorderColor(tcell.ColorRed).
		// SetTitleColor(tcell.ColorRed)
	alertContainer.SetDirection(tview.FlexRow)

	alertContainer.AddItem(e.errorText, 0, 1, false)
	alertContainer.AddItem(legend, 1, 1, false)

	// Modal
	e.alertContainer = alertContainer
	e.alertModal = CreateModal(alertContainer, 100, 15)

	e.setKeyBindings()
	e.lastFocus = e.app.GetFocus()
	e.errorString = errorText

	e.app.appPages.RemovePage("modal")
	e.app.appPages.AddPage("modal", e.alertModal, true, false)
	e.app.appPages.ShowPage("modal")
	e.app.SetFocus(e.alertContainer)
}

func (e *ErrorModal) setKeyBindings() {
	e.alertContainer.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEscape {
			e.Close()
		}

		if event.Key() == tcell.KeyRune {
			switch event.Rune() {
			case 'y':
				// copy to clipboard errorTExt
				clipboard.Write(clipboard.FmtText, []byte(e.errorString))

				// highlight the error text, use the same logic as the code below
				e.errorText.SetText("[black:yellow]" + e.errorString)

				time.AfterFunc(75*time.Millisecond, func() {
					// return to default
					e.errorText.SetText("[::]" + e.errorString)
					e.app.Draw()
				})
			}
		}

		return event
	})
}

func (e *ErrorModal) Close() {
	e.app.appPages.SwitchToPage("app")
	e.app.SetFocus(e.lastFocus)
}
