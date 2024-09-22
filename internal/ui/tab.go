package ui

import (
	"github.com/alfonzm/lazydb/internal/db"
	"github.com/rivo/tview"
)

type Tab struct {
	name      string
	lastFocus tview.Primitive
	pages     *tview.Pages
	app       *App

	sidebar     *Sidebar
	results     *Results
	connections *Connections

	dbClient       *db.DBClient
	databaseTables map[string][]string
}

func NewTab(app *App) (*Tab, error) {
	tab := &Tab{
		pages: tview.NewPages(),
		name:  "New Tab",
		app:   app,
	}

	conns, err := NewConnections(tab)
	if err != nil {
		return nil, err
	}

	tab.connections = conns

	tab.pages.AddPage("connections", conns.view, true, true)

	return tab, nil
}

func (t *Tab) ConnectDatabase(url string, dbName string) error {
	db, err := db.NewDBClient(url)
	if err != nil {
		return err
	}

	t.dbClient = db

	pages := t.pages

	// Setup results component
	results, err := NewResults(t.app, pages, db)

	// Setup sidebar components
	sidebar, err := NewSidebar(t, t.app.Application, db, results)
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
		AddItem(sidebar.view, 0, 1, true).
		AddItem(results.view, 0, 6, false)

	pages.AddPage("main", main, true, true)
	pages.AddPage("editor", cellEditor.view, true, false)

	t.sidebar = sidebar
	t.results = results

	// Switch to main page
	pages.SwitchToPage("main")
	t.app.SetFocus(sidebar.list)

	t.UpdateTabName(dbName)

	t.CacheDatabaseTables()

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
	case t.results.structure.columnsTable:
		t.app.SetFocus(t.sidebar.results.structure.indexesTable)
	case t.results.structure.indexesTable:
		t.app.SetFocus(t.sidebar.list)
	}
}

func (t *Tab) FocusFindTable() {
	t.app.SetFocus(t.sidebar.list)
	t.app.SetFocus(t.sidebar.filter)
}

func (t *Tab) UpdateTabName(name string) {
	t.name = name
	t.app.RenderTabHeaders()
}

func (t *Tab) CacheDatabaseTables() error {
	t.databaseTables = make(map[string][]string)

	dbNames, err := t.dbClient.GetDatabases()
	if err != nil {
		return err
	}

	// for each database, get tables
	for _, dbName := range dbNames {
		tableNames, err := t.dbClient.GetTablesOfDatabase(dbName)
		if err != nil {
			return err
		}

		t.databaseTables[dbName] = tableNames
	}

	return nil
}
