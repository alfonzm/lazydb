package ui

import (
	"fmt"

	"github.com/alfonzm/lazydb/internal/db"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type SortColumn struct {
	Name      string
	Ascending bool
}

type Results struct {
	app           *tview.Application
	view          *tview.Flex
	table         *tview.Table
	filter        *tview.InputField
	pages         *tview.Pages
	db            *db.DBClient
	editor        *Editor
	selectedTable string
	sortColumn    SortColumn
	dbColumns     []db.Column
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

	r.dbColumns = dbColumns

	orderBy := ""
	if r.sortColumn.Name != "" {
		orderBy = r.sortColumn.Name
		if !r.sortColumn.Ascending {
			orderBy = fmt.Sprintf("%s DESC", orderBy)
		}
	}

	dbRecords, err := r.db.GetRecords(table, where, orderBy)
	if err != nil {
		return fmt.Errorf("Error getting records")
	}

	r.table.Clear()

	// set headers from columns
	for i, column := range dbColumns {
		var columnName string = column.Name

		// append sort arrow to column name
		if r.sortColumn.Name == column.Name {
			if r.sortColumn.Ascending {
				columnName = fmt.Sprintf("%s ↑", column.Name)
			} else {
				columnName = fmt.Sprintf("%s ↓", column.Name)
			}
		}

		r.table.SetCell(
			0,
			i,
			tview.NewTableCell(columnName).SetAlign(tview.AlignCenter).SetSelectable(true),
		)
	}
	r.table.SetSelectable(true, true)

	// set cell selection function
	r.table.SetSelectedFunc(func(row, column int) {
		// handle sort if headers
		if row == 0 {
			r.toggleSort(dbColumns[column].Name)
			return
		}

		// else show cell editor
		r.pages.ShowPage("editor")
		r.editor.textArea.SetText(r.table.GetCell(row, column).Text, true)
	})

	// Iterate over records and fill table
	for rowIndex, record := range dbRecords {
		for columnIndex, column := range dbColumns {
			recordValue, ok := record[column.Name]

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
		SetFieldBackgroundColor(tcell.ColorNone).
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
		if event.Key() == tcell.KeyEscape {
			r.clearFilter()
		}

		if event.Key() == tcell.KeyRune {
			switch {
			case event.Rune() == '/':
				r.app.SetFocus(r.filter)
			case event.Rune() == 's':
				r.toggleSortForCell()
			// D will delete the selected row
			case event.Rune() == 'd':
				r.deleteRow()
			}
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

func (r *Results) toggleSort(columnName string) {
	row, col := r.table.GetSelection()

	if columnName == "" {
		columnName = r.dbColumns[col].Name
	}

	if r.sortColumn.Name == columnName {
		// toggle from ASC, DESC, and none
		switch r.sortColumn.Ascending {
		case true:
			r.sortColumn.Ascending = false
		case false:
			r.sortColumn.Name = ""
		}
	} else {
		r.sortColumn.Name = columnName
		r.sortColumn.Ascending = true
	}

	// re-render table
	r.RenderTable(r.selectedTable, r.filter.GetText())

	// reselect current cell
	r.table.Select(row, col)

	return
}

func (r *Results) toggleSortForCell() {
	r.toggleSort("")
}

func (r *Results) deleteRow() {
	row, col := r.table.GetSelection()
	if row == 0 {
		return
	}

	columns := r.dbColumns

	where := ""

	for i, col := range columns {
		if col.DataType == "longtext" || col.DataType == "text" || col.DataType == "blob" {
			continue
		}

		cell := r.table.GetCell(row, i)

		if cell.Text == "" {
			continue
		}

		whereClause := fmt.Sprintf("%s = '%s'", col.Name, cell.Text)

		if i == 0 {
			where = whereClause
		} else {
			where = fmt.Sprintf("%s AND %s", where, whereClause)
		}
	}

	if err := r.db.DeleteRecord(r.selectedTable, where); err != nil {
		fmt.Printf("Error deleting record: %v\n", err)
		return
	}

	r.RenderTable(r.selectedTable, r.filter.GetText())
	r.table.Select(row, col)
}
