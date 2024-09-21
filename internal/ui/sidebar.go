package ui

import (
	"fmt"
	"strings"

	"github.com/alfonzm/lazydb/internal/db"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type Sidebar struct {
	tab     *Tab
	app     *tview.Application
	view    *tview.Flex
	list    *tview.List
	db      *db.DBClient
	results *Results
	filter  *tview.InputField
}

func NewSidebar(
	tab *Tab,
	app *tview.Application,
	db *db.DBClient,
	results *Results,
) (*Sidebar, error) {
	list := tview.NewList()
	filter := tview.NewInputField()

	// Sidebar main container
	view := tview.NewFlex()
	view.SetTitle("Tables")
	view.SetBorder(true)
	view.SetDirection(tview.FlexRow).
		AddItem(filter, 1, 1, false).
		AddItem(list, 0, 1, true)

	sidebar := &Sidebar{
		view:    view,
		list:    list,
		db:      db,
		results: results,
		app:     app,
		filter:  filter,
		tab:     tab,
	}

	// Render all components
	if err := sidebar.renderTableList(""); err != nil {
		return nil, fmt.Errorf("Failed to render table list: %w", err)
	}
	sidebar.renderFilterField()
	sidebar.setKeyBindings()

	return sidebar, nil
}

func (sidebar *Sidebar) setKeyBindings() {
	sidebar.view.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if sidebar.app.GetFocus() != sidebar.filter {
			// Ctrl+n / Ctrl+p to navigate the list
			if event.Key() == tcell.KeyCtrlN {
				sidebar.list.SetCurrentItem(sidebar.list.GetCurrentItem() + 1)
				tableName, _ := sidebar.list.GetItemText(sidebar.list.GetCurrentItem())
				sidebar.selectTable(tableName, false)
				return event
			}
			if event.Key() == tcell.KeyCtrlP {
				sidebar.list.SetCurrentItem(sidebar.list.GetCurrentItem() - 1)
				tableName, _ := sidebar.list.GetItemText(sidebar.list.GetCurrentItem())
				sidebar.selectTable(tableName, false)
				return event
			}

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
			if event.Key() == tcell.KeyEscape && sidebar.filter.GetText() != "" {
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

	for _, table := range tableNames {
		if filter != "" && !strings.Contains(strings.ToLower(table), strings.ToLower(filter)) {
			continue
		}

		s.list.AddItem(table, "", 0, func() {
			s.selectTable(table, true)
		})
	}

	return nil
}

func (s *Sidebar) selectTable(table string, focus bool) {
	s.results.ClearSort()
	s.results.RenderTable(table, "")

	if focus {
		s.results.Focus()
	}

	s.tab.UpdateTabName(table)
}

func (s *Sidebar) renderFilterField() {
	s.filter.SetLabel("Filter ")
	s.filter.SetFieldBackgroundColor(tcell.ColorNone)

	s.filter.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		currentItem := s.list.GetCurrentItem()

		// Set the focus back to the list on tab key
		if event.Key() == tcell.KeyTab {
			s.app.SetFocus(s.list)
			return event
		}

		// Ctrl+n / Ctrl+p to navigate the list
		if event.Key() == tcell.KeyCtrlN {
			s.list.SetCurrentItem(currentItem + 1)
			tableName, _ := s.list.GetItemText(s.list.GetCurrentItem())
			s.selectTable(tableName, false)
			return event
		}
		if event.Key() == tcell.KeyCtrlP {
			s.list.SetCurrentItem(currentItem - 1)
			tableName, _ := s.list.GetItemText(s.list.GetCurrentItem())
			s.selectTable(tableName, false)
			return event
		}

		// Filter the table list in real time
		currentText := s.filter.GetText()
		if event.Key() == tcell.KeyEscape {
			if currentText != "" {
				s.renderTableList(currentText)
				s.app.SetFocus(s.list)
			}
			return event
		}

		// Remove the last character if present
		if event.Key() == tcell.KeyBackspace || event.Key() == tcell.KeyBackspace2 {
			if len(currentText) > 0 {
				currentText = currentText[:len(currentText)-1]
			}
		}

		// Append the current rune to the filter text for rendering
		if event.Key() == tcell.KeyRune {
			currentText += string(event.Rune())
		}

		// Render the table list and filter in real time
		s.renderTableList(currentText)

		s.list.SetCurrentItem(currentItem)

		return event
	})

	s.filter.SetDoneFunc(func(key tcell.Key) {
		if key == tcell.KeyEnter {
			currentItem := s.list.GetCurrentItem()

			// if current item is > 0, select the item
			if currentItem > 0 {
				tableName, _ := s.list.GetItemText(currentItem)
				s.selectTable(tableName, true)
				return
			}

			s.app.SetFocus(s.list)
		}
	})
}
