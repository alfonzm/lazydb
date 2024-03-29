package ui

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func Start() error {
	app := tview.NewApplication()

	// Press q to quit
	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyRune {
			if event.Rune() == 'q' {
				app.Stop()
			}
		}
		return event
	})

  // Setup main page
	sidebar := NewSidebar()
	results := NewResults()
	flex := tview.NewFlex().
		AddItem(sidebar, 0, 1, false).
		AddItem(results, 0, 5, false)

  // Setup Pages
	pages := tview.NewPages()
	pages.AddPage("main", flex, true, true)

  // Run the app
	if err := app.SetRoot(pages, true).Run(); err != nil {
		return err
	}

	return nil
}
