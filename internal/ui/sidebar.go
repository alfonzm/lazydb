package ui

import "github.com/rivo/tview"

func NewSidebar() tview.Primitive {
	return tview.NewBox().SetBorder(true).SetTitle("Sidebar")
}
