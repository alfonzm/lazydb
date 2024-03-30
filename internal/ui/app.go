package ui

import (
	"github.com/alfonzm/lazydb/internal/db"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type App struct {
	*tview.Application
	sidebar *Sidebar
	results *Results
}

func Start(db *db.DBClient) error {
	app := &App{Application: tview.NewApplication()}

	// Setup results component
	results, err := NewResults(app.Application, db)

	// Setup sidebar components
	sidebar, err := NewSidebar(app.Application, db, results)
	if err != nil {
		return err
	}

	flex := tview.NewFlex().
		AddItem(sidebar.view, 0, 1, false).
		AddItem(results.view, 0, 6, false)

		// Setup Pages
	pages := tview.NewPages()
	pages.AddPage("main", flex, true, true)

	app.setKeyBindings()
	app.sidebar = sidebar
	app.results = results

	// Run the app
	if err := app.SetRoot(pages, true).SetFocus(sidebar.list).Run(); err != nil {
		return err
	}

	return nil
}

// set keybindings
func (app *App) setKeyBindings() {
	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		// If the focus is on an input field, early return
		if _, ok := app.GetFocus().(*tview.InputField); ok {
			return event
		}

		if event.Key() == tcell.KeyRune {
			switch event.Rune() {
			case 'q':
				app.Stop()
			}
		}

		if event.Key() == tcell.KeyTab {
			switch app.GetFocus() {
			case app.sidebar.list:
				app.SetFocus(app.results.table)
			case app.results.table:
				app.SetFocus(app.sidebar.list)
			}
		}

		return event
	})
}
