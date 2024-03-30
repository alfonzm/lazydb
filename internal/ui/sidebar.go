package ui

import (
	"fmt"

	"github.com/alfonzm/lazydb/internal/db"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type Sidebar struct {
	app     *tview.Application
	view    *tview.Flex
	list    *tview.List
	db      *db.DBClient
	results *Results
}

func NewSidebar(app *tview.Application, db *db.DBClient, results *Results) (*Sidebar, error) {
	list := tview.NewList()

	// Define container for the sidebar
	view := tview.NewFlex()
	view.SetTitle("Tables")
	view.SetBorder(true)
	view.SetDirection(tview.FlexRow).
		AddItem(list, 0, 1, false)

	view.SetBorder(true).SetTitle("Sidebar")

	sidebar := &Sidebar{
		view:    view,
		list:    list,
		db:      db,
		results: results,
		app:     app,
	}

	if err := sidebar.renderTableList(); err != nil {
		return nil, fmt.Errorf("Failed to render table list: %w", err)
	}

	sidebar.setKeyBindings()

	return sidebar, nil
}

func (sidebar *Sidebar) setKeyBindings() {
	sidebar.view.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyRune {
			switch event.Rune() {
			case 'j':
				// pressing j at the end of the list goes to the top
				if sidebar.list.GetItemCount()-1 == sidebar.list.GetCurrentItem() {
					sidebar.list.SetCurrentItem(0)
				} else {
					sidebar.list.SetCurrentItem(sidebar.list.GetCurrentItem() + 1)
				}
			case 'k':
				sidebar.list.SetCurrentItem(sidebar.list.GetCurrentItem() - 1)
			case '/':
				// focus filter UI
			}
		}
		return event
	})
}

func (s *Sidebar) renderTableList() error {
	tableNames, err := s.db.GetTables()
	if err != nil {
		return fmt.Errorf("Failed to get tables: %w", err)
	}

	s.list.ShowSecondaryText(false).SetHighlightFullLine(true).
		SetTitle("Tables")

	for _, table := range tableNames {
		s.list.AddItem(table, "", 0, func() {
			s.results.RenderTable(table)
			s.app.SetFocus(s.results.table)
		})
	}

	return nil
}
