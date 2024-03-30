package ui

import (
	"fmt"

	"github.com/alfonzm/lazydb/internal/db"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type Results struct {
	app           *tview.Application
	view          *tview.Flex
	table         *tview.Table
	db            *db.DBClient
	filter        *tview.InputField
	pages         *tview.Pages
	selectedTable string
}

func NewResults(app *tview.Application, pages *tview.Pages, db *db.DBClient) (*Results, error) {
	table := tview.NewTable()
	filter := tview.NewInputField()

	view := tview.NewFlex()
	view.SetBorder(true).
		SetTitle("Results").
		SetBorder(true)
	view.SetDirection(tview.FlexRow).
		AddItem(filter, 1, 1, false).
		AddItem(table, 0, 1, false)

	results := &Results{
		app:    app,
		table:  table,
		view:   view,
		db:     db,
		filter: filter,
		pages:  pages,
	}

	results.renderFilterField()
	results.setKeyBindings()

	return results, nil
}

func (r *Results) RenderTable(table string, where string) error {
	r.selectedTable = table

	dbColumns, err := r.db.GetColumns(table)
	if err != nil {
		return fmt.Errorf("Error getting columns")
	}

	dbRecords, err := r.db.GetRecords(table, where)
	if err != nil {
		return fmt.Errorf("Error getting records")
	}

	r.table.Clear()

	// set headers from columns
	for i, column := range dbColumns {
		r.table.SetCell(
			0,
			i,
			tview.NewTableCell(column).SetAlign(tview.AlignCenter).SetSelectable(true),
		)
	}
	r.table.SetSelectable(true, true)
	r.table.SetSelectedFunc(func(row, column int) {
		// show editor page
		r.pages.ShowPage("editor")
	})

	// Iterate over records and fill table
	for rowIndex, record := range dbRecords {
		for columnIndex, colName := range dbColumns {
			recordValue, ok := record[colName]

			cellString := ""

			// if DB value is null, set valStr to "NULL"
			if ok && recordValue == nil {
				// TODO: For some reason this is not working
				// Maybe there is a better way to check for NULL DB values cellString = "NULL"
			} else if ok && recordValue != nil {
				cellString = fmt.Sprintf("%v", recordValue)
			}

			cell := tview.NewTableCell(cellString).SetAlign(tview.AlignLeft).SetSelectable(true)
			r.table.SetCell(rowIndex+1, columnIndex, cell)
		}
	}

	r.table.SetFixed(1, 0)
	r.table.ScrollToBeginning()
	r.table.Select(0, 0)

	return nil
}

func (r *Results) renderFilterField() {
	r.filter.SetLabel("WHERE ").
		SetFieldBackgroundColor(tcell.ColorBlack).
		SetDoneFunc(func(key tcell.Key) {
			if key == tcell.KeyEnter {
				where := r.filter.GetText()
				r.RenderTable(r.selectedTable, where)

				r.app.SetFocus(r.table)
			}
		})
}

func (r *Results) setKeyBindings() {
	// Table key bindings
	r.table.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		// press /
		if event.Key() == tcell.KeyRune && event.Rune() == '/' {
			r.app.SetFocus(r.filter)
		}

		// press escape clears filter
		if event.Key() == tcell.KeyEscape {
			r.clearFilter()
		}
		return event
	})

	// Filter field key bindings
	r.filter.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEscape {
			r.clearFilter()
		}
		return event
	})
}

func (r *Results) clearFilter() {
	r.filter.SetText("")
	r.RenderTable(r.selectedTable, "")
	r.app.SetFocus(r.table)
}
