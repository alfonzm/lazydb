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
	app                  *tview.Application
	pages                *tview.Pages
	db                   *db.DBClient
	view                 *tview.Pages
	resultsTable         *tview.Table
	columnsTable         *tview.Table
	filter               *tview.InputField
	editor               *Editor
	selectedTable        string
	sortColumn           SortColumn
	dbColumns            []db.Column
	selectedRowForDelete int
}

func NewResults(app *tview.Application, pages *tview.Pages, db *db.DBClient) (*Results, error) {
	// Setup Results page
	table := tview.NewTable()
	filter := tview.NewInputField()

	resultsPage := tview.NewFlex()
	resultsPage.SetBorder(true).
		SetTitle("Results").
		SetBorder(true)
	resultsPage.SetDirection(tview.FlexRow).
		AddItem(filter, 1, 1, false).
		AddItem(table, 0, 1, false)

		// Setup Columns page
	columnsTable := tview.NewTable()

	columnsPage := tview.NewFlex()
	columnsPage.SetBorder(true).
		SetTitle("Columns").
		SetBorder(true)
	columnsPage.SetDirection(tview.FlexRow)
	columnsPage.AddItem(columnsTable, 0, 1, false)

	view := tview.NewPages()
	view.AddPage("results", resultsPage, true, true)
	view.AddPage("columns", columnsPage, true, false)
	// view.ShowPage("columns")

	results := &Results{
		app:          app,
		resultsTable: table,
		columnsTable: columnsTable,
		view:         view,
		db:           db,
		filter:       filter,
		pages:        pages,
	}

	results.renderFilterField()
	results.setKeyBindings()

	return results, nil
}

func (r *Results) RenderResultsTable(table string, where string) error {
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

	r.resultsTable.Clear()

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

		r.resultsTable.SetCell(
			0,
			i,
			tview.NewTableCell(columnName).SetAlign(tview.AlignCenter).SetSelectable(true),
		)
	}
	r.resultsTable.SetSelectable(true, true)

	// set cell selection function
	r.resultsTable.SetSelectedFunc(func(row, column int) {
		// handle sort if headers
		if row == 0 {
			r.toggleSort(dbColumns[column].Name)
			return
		}

		// else show cell editor
		r.pages.ShowPage("editor")
		r.editor.textArea.SetText(r.resultsTable.GetCell(row, column).Text, true)
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
			r.resultsTable.SetCell(rowIndex+1, columnIndex, cell)
		}
	}

	r.resultsTable.SetFixed(1, 0)
	r.resultsTable.ScrollToBeginning()
	r.resultsTable.Select(0, 0)

	return nil
}

func (r *Results) RenderColumnsTable(table string) error {
	dbColumns, err := r.db.GetColumns(table)
	if err != nil {
		return fmt.Errorf("Error getting columns")
	}

	r.columnsTable.Clear()

	// set headers from columns
	for i, column := range dbColumns {
		r.columnsTable.SetCell(
			0,
			i,
			tview.NewTableCell(column.Name).SetAlign(tview.AlignCenter).SetSelectable(true),
		)
	}

	return nil
}

func (r *Results) Focus() {
	// focus  the active page content (table or columns)
	frontPage, _ := r.view.GetFrontPage()

	switch frontPage {
	case "results":
		r.app.SetFocus(r.resultsTable)
	case "columns":
		r.app.SetFocus(r.columnsTable)
	}
}

func (r *Results) renderFilterField() {
	r.filter.SetLabel("WHERE ").
		SetFieldBackgroundColor(tcell.ColorNone).
		SetDoneFunc(func(key tcell.Key) {
			if key == tcell.KeyEnter {
				where := r.filter.GetText()
				r.RenderResultsTable(r.selectedTable, where)

				r.app.SetFocus(r.resultsTable)
			}
		})
}

func (r *Results) setKeyBindings() {
	// Table key bindings
	r.resultsTable.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEscape {
			r.clearFilter()
		}

		if event.Key() == tcell.KeyRune {
			switch {
			case event.Rune() == '/':
				r.app.SetFocus(r.filter)
			case event.Rune() == 's':
				r.toggleSortForCell()
			case event.Rune() == 'd':
				r.attemptDeleteRow()
			case event.Rune() == '1':
				r.view.ShowPage("columns")
				r.app.SetFocus(r.columnsTable)
			case event.Rune() == '2':
				r.view.ShowPage("results")
				r.app.SetFocus(r.resultsTable)
			}
		}

		return event
	})

	// Filter field key bindings
	r.filter.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEscape {
			r.app.SetFocus(r.resultsTable)
		}
		return event
	})
}

func (r *Results) clearFilter() {
	if r.selectedTable == "" {
		return
	}

	r.filter.SetText("")
	r.RenderResultsTable(r.selectedTable, "")
	r.app.SetFocus(r.resultsTable)
}

func (r *Results) toggleSort(columnName string) {
	row, col := r.resultsTable.GetSelection()

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
	r.RenderResultsTable(r.selectedTable, r.filter.GetText())

	// reselect current cell
	r.resultsTable.Select(row, col)

	return
}

func (r *Results) toggleSortForCell() {
	r.toggleSort("")
}

func (r *Results) attemptDeleteRow() {
	row, _ := r.resultsTable.GetSelection()

	// if the selected row is already selected for delete, confirm deletion
	if r.selectedRowForDelete == row {
		r.deleteRow(r.selectedRowForDelete)
		r.selectedRowForDelete = 0
		return
	}

	if r.selectedRowForDelete != 0 {
		// clear the previous selected row for delete
		for i := 0; i < len(r.dbColumns); i++ {
			cell := r.resultsTable.GetCell(r.selectedRowForDelete, i)
			cell.SetBackgroundColor(tcell.ColorDefault)
		}
	}

	// if the selected row is not already selected for delete, highlight it
	r.selectedRowForDelete = row

	// set the selected row to red background
	for i := 0; i < len(r.dbColumns); i++ {
		cell := r.resultsTable.GetCell(row, i)
		cell.SetBackgroundColor(tcell.ColorRed)
	}
}

func (r *Results) deleteRow(row int) {
	rowToDelete, col := r.resultsTable.GetSelection()

	if row > 0 {
		rowToDelete = row
	}

	if rowToDelete == 0 {
		return
	}

	columns := r.dbColumns

	where := ""

	for i, col := range columns {
		if col.DataType == "longtext" || col.DataType == "text" || col.DataType == "blob" {
			continue
		}

		cell := r.resultsTable.GetCell(rowToDelete, i)

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

	r.RenderResultsTable(r.selectedTable, r.filter.GetText())
	r.resultsTable.Select(rowToDelete, col)
}
