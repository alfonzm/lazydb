package ui

import (
	"github.com/alfonzm/lazydb/internal/db"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type App struct {
	*tview.Application
	pages       *tview.Pages
	sidebar     *Sidebar
	results     *Results
	connections *Connections
	dbClient    *db.DBClient
}

func Start() error {
	// Setup Pages
	pages := tview.NewPages()

	app := &App{
		Application: tview.NewApplication(),
		pages:       pages,
	}

	conns, err := NewConnections(app, app.dbClient)
	if err != nil {
		return err
	}

	app.connections = conns

	pages.AddPage("connections", conns.view, true, true)
	app.setKeyBindings()

	if err := app.SetRoot(pages, true).SetFocus(conns.list).Run(); err != nil {
		return err
	}

	return nil
}

func (app *App) Connect(url string) error {
	db, err := db.NewDBClient(url)
	if err != nil {
		panic(err)
		// return err
	}

	pages := app.pages

	// Setup results component
	results, err := NewResults(app.Application, pages, db)

	// Setup sidebar components
	sidebar, err := NewSidebar(app.Application, db, results)
	if err != nil {
		return err
	}

	// Setup record cellEditor component
	cellEditor, err := NewCellEditor(app.Application, pages, results, db)
	if err != nil {
		return err
	}

	results.cellEditor = cellEditor

	main := tview.NewFlex().
		AddItem(sidebar.view, 0, 1, false).
		AddItem(results.view, 0, 6, false)

	pages.AddPage("main", main, true, false)
	pages.AddPage("editor", cellEditor.view, true, false)

	app.sidebar = sidebar
	app.results = results

	// Switch to main page
	pages.SwitchToPage("main")
	app.SetFocus(sidebar.list)

	return nil
}

// set keybindings
func (app *App) setKeyBindings() {
	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		// If the focus is on an input field, early return
		if _, ok := app.GetFocus().(*tview.InputField); ok {
			return event
		}

		// If the focus is on Editor, early return
		if app.GetFocus() == app.results.cellEditor.textArea {
			return event
		}

		if event.Key() == tcell.KeyTab {
			switch app.GetFocus() {
			case app.sidebar.list:
				app.results.Focus()
			case app.results.resultsTable:
				app.SetFocus(app.sidebar.list)
			case app.results.columnsTable:
				app.SetFocus(app.results.indexesTable)
			case app.results.indexesTable:
				app.SetFocus(app.sidebar.list)
			}
			return event
		}

		if event.Key() == tcell.KeyRune {
			switch event.Rune() {
			case 'q':
				app.Stop()
			case '0':
				app.pages.SwitchToPage("connections")
			}

			return event
		}

		switch event.Key() {
		case tcell.KeyCtrlF:
			app.SetFocus(app.sidebar.list)
			app.SetFocus(app.sidebar.filter)
		case tcell.KeyCtrlR:
			app.results.Focus()
		case tcell.KeyCtrlT:
			app.SetFocus(app.sidebar.list)
		}

		return event
	})
}
