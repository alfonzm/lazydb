package ui

import (
	"github.com/alfonzm/lazydb/internal/db"
	"github.com/rivo/tview"
)

type Tab struct {
	app       *tview.Application
	dbClient  *db.DBClient
	name      string
	lastFocus tview.Primitive
	pages     *tview.Pages

	sidebar     *Sidebar
	results     *Results
	connections *Connections
}

func NewTab(app *tview.Application, dbClient *db.DBClient) (*Tab, error) {
	tab := &Tab{
		app:   app,
		pages: tview.NewPages(),
		name:  "New Tab",
	}

	conns, err := NewConnections(tab, dbClient)
	if err != nil {
		return nil, err
	}

	tab.connections = conns

	tab.pages.AddPage("connections", conns.view, true, true)

	return tab, nil
}

func (t *Tab) Connect(url string) error {
	db, err := db.NewDBClient(url)
	if err != nil {
		return err
	}

	pages := t.pages

	// Setup results component
	results, err := NewResults(t.app, pages, db)

	// Setup sidebar components
	sidebar, err := NewSidebar(t.app, db, results)
	if err != nil {
		return err
	}

	// Setup record cellEditor component
	cellEditor, err := NewCellEditor(t.app, pages, results, db)
	if err != nil {
		return err
	}

	results.cellEditor = cellEditor

	main := tview.NewFlex().
		AddItem(sidebar.view, 0, 1, false).
		AddItem(results.view, 0, 6, false)

	pages.AddPage("main", main, true, false)
	pages.AddPage("editor", cellEditor.view, true, false)

	t.sidebar = sidebar
	t.results = results

	// Switch to main page
	pages.SwitchToPage("main")
	t.app.SetFocus(sidebar.list)

	return nil
}

func (t *Tab) OnActivate() {
	if t.lastFocus != nil {
		t.app.SetFocus(t.lastFocus)
	}
}

func (t *Tab) OnDeactivate() {
	t.lastFocus = t.app.GetFocus()
}

func (t *Tab) OnPressTab() {
	if t.sidebar == nil || t.results == nil {
		return
	}

	switch t.app.GetFocus() {
	case t.sidebar.list:
		t.results.Focus()
	case t.results.resultsTable:
		t.app.SetFocus(t.sidebar.list)
	case t.results.columnsTable:
		t.app.SetFocus(t.sidebar.results.indexesTable)
	case t.results.indexesTable:
		t.app.SetFocus(t.sidebar.list)
	}
}

func (t *Tab) FocusFindTable() {
	t.app.SetFocus(t.sidebar.list)
	t.app.SetFocus(t.sidebar.filter)
}