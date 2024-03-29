package ui

import "github.com/rivo/tview"

func NewResults() tview.Primitive {
	return tview.NewBox().SetBorder(true).SetTitle("Results")
}
