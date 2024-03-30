package ui

import (
	"fmt"

	"github.com/alfonzm/lazydb/internal/db"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type Sidebar struct {
	view *tview.Flex
	list *tview.List
	db   *db.DBClient
}

func NewSidebar(db *db.DBClient) (*Sidebar, error) {
	tableNames, err := db.GetTables()
	if err != nil {
		return nil, fmt.Errorf("Failed to get tables: %w", err)
	}

	// Display the tables in a list
	list := tview.NewList()
	list.ShowSecondaryText(false).SetHighlightFullLine(true).
		SetTitle("Tables")

	for _, table := range tableNames {
		list.AddItem(table, "", 0, nil)
	}

	// Define container for the sidebar
	view := tview.NewFlex()
	view.SetTitle("Tables")
	view.SetBorder(true)
	view.SetDirection(tview.FlexRow).
		AddItem(list, 0, 1, false)

	view.SetBorder(true).SetTitle("Sidebar")

	sidebar := &Sidebar{
		view: view,
		list: list,
		db:   db,
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
				// app.SetFocus(filterUI)
			}
		}
		return event
	})
}
