package ui

import (
	"fmt"
	"strings"

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
	filter  *tview.InputField
}

func NewSidebar(app *tview.Application, db *db.DBClient, results *Results) (*Sidebar, error) {
	list := tview.NewList()

	// Filter
	filter := tview.NewInputField()

	// Define container for the sidebar
	view := tview.NewFlex()
	view.SetTitle("Tables")
	view.SetBorder(true)
	view.SetDirection(tview.FlexRow).
		AddItem(filter, 1, 1, false).
		AddItem(list, 0, 1, false)

	sidebar := &Sidebar{
		view:    view,
		list:    list,
		db:      db,
		results: results,
		app:     app,
		filter:  filter,
	}

	if err := sidebar.renderTableList(""); err != nil {
		return nil, fmt.Errorf("Failed to render table list: %w", err)
	}

	filter.SetLabel("Filter")
	filter.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEscape {
			filter.SetText("")
			sidebar.renderTableList("")
			app.SetFocus(list)
		}
		return event
	})
	filter.SetDoneFunc(func(key tcell.Key) {
		if key == tcell.KeyEnter {
			filterText := filter.GetText()
			sidebar.renderTableList(filterText)
			app.SetFocus(list)
		}
	})

	sidebar.setKeyBindings()

	return sidebar, nil
}

func (sidebar *Sidebar) setKeyBindings() {
	sidebar.view.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if sidebar.app.GetFocus() != sidebar.filter {
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
					sidebar.app.SetFocus(sidebar.filter)
					return nil // prevents adding '/' char to the input field
				}
			}

      // Clear filter when pressing escape
			if event.Key() == tcell.KeyEscape {
				sidebar.filter.SetText("")
				sidebar.renderTableList("")
			}
		}
		return event
	})

	sidebar.list.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyTab {
			// Hack: For some reason, tab key moves the focus to the next item
			// so we need to move it back to the previous item
			sidebar.list.SetCurrentItem(sidebar.list.GetCurrentItem() - 1)
		}
		return event
	})
}

func (s *Sidebar) renderTableList(filter string) error {
	s.list.Clear()

	tableNames, err := s.db.GetTables()
	if err != nil {
		return fmt.Errorf("Failed to get tables: %w", err)
	}

	s.list.ShowSecondaryText(false).SetHighlightFullLine(true).
		SetTitle("Tables")

		// filter by if table name contains filter
	for _, table := range tableNames {
		if filter != "" && !strings.Contains(strings.ToLower(table), strings.ToLower(filter)) {
			continue
		}

		s.list.AddItem(table, "", 0, func() {
			s.results.RenderTable(table)
			s.app.SetFocus(s.results.table)
		})
	}

	return nil
}
