package ui

import (
	"github.com/alfonzm/lazydb/internal/db"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func Start(db *db.DBClient) error {
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

	// Setup results component
  results, err := NewResults(db)

	// Setup sidebar components
  sidebar, err := NewSidebar(app, db, results)
	if err != nil {
		return err
	}

	flex := tview.NewFlex().
		AddItem(sidebar.view, 0, 1, false).
		AddItem(results.view, 0, 6, false)

  // Setup Pages
	pages := tview.NewPages()
	pages.AddPage("main", flex, true, true)

	// Run the app
	if err := app.SetRoot(pages, true).SetFocus(sidebar.list).Run(); err != nil {
		return err
	}

	return nil
}
