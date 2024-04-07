package ui

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/alfonzm/lazydb/internal/db"
	"github.com/atotto/clipboard"
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
	indexesTable         *tview.Table
	filter               *tview.InputField
	editor               *Editor
	selectedTable        string
	sortColumn           SortColumn
	dbColumns            []db.Column
	selectedRowForDelete int
}

func NewResults(app *tview.Application, pages *tview.Pages, db *db.DBClient) (*Results, error) {
	// Setup Results page
	resultsTable := tview.NewTable()
	filter := tview.NewInputField()
	filter.SetAutocompleteStyles(
		tcell.Color237,
		tcell.StyleDefault,
		tcell.StyleDefault.Foreground(tcell.Color237).Background(tcell.Color246),
	)

	resultsPage := tview.NewFlex()
	resultsPage.SetBorder(true).
		SetTitle("Results").
		SetBorder(true)
	resultsPage.SetDirection(tview.FlexRow).
		AddItem(filter, 1, 1, false).
		AddItem(resultsTable, 0, 1, false)

	// Setup Columns page
	columnsTable := tview.NewTable()
	indexesTable := tview.NewTable()

	columnsPage := tview.NewFlex()
	columnsPage.SetBorder(false).
		SetTitle("Columns")
	columnsPage.SetDirection(tview.FlexRow)
	columnsPage.AddItem(columnsTable, 0, 4, false)
	columnsPage.AddItem(indexesTable, 0, 1, false)

	view := tview.NewPages()
	view.AddPage("results", resultsPage, true, true)
	view.AddPage("columns", columnsPage, true, false)

	results := &Results{
		app:          app,
		resultsTable: resultsTable,
		columnsTable: columnsTable,
		indexesTable: indexesTable,
		view:         view,
		db:           db,
		filter:       filter,
		pages:        pages,
	}

	results.renderFilterField()
	results.setKeyBindings()

	return results, nil
}

// RenderTable renders the table with the given name and optional where clause
// It will also render the columns and indexes tables
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
		r.pages.SwitchToPage("editor")
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

	r.RenderColumnsTable(table, dbColumns)

	return nil
}

func (r *Results) RenderColumnsTable(table string, dbColumns []db.Column) error {
	r.columnsTable.Clear()

	columnMetaColumns := []string{"Field", "Type", "Null", "Key", "Default", "Extra"}

	// set headers from columns
	for i, column := range columnMetaColumns {
		r.columnsTable.SetCell(
			0,
			i,
			tview.NewTableCell(column).SetAlign(tview.AlignCenter).SetSelectable(true),
		)
	}

	// Iterate over records and fill table
	for i, col := range dbColumns {
		defaultValue := ""
		if col.Default.Valid {
			defaultValue = col.Default.String
		} else if col.Null {
			defaultValue = "NULL"
		}

		r.columnsTable.SetCell(
			i+1,
			0,
			tview.NewTableCell(col.Name).SetAlign(tview.AlignLeft).SetSelectable(true),
		)
		r.columnsTable.SetCell(
			i+1,
			1,
			tview.NewTableCell(col.DataType).SetAlign(tview.AlignLeft).SetSelectable(true),
		)
		r.columnsTable.SetCell(
			i+1,
			2,
			tview.NewTableCell(fmt.Sprintf("%t", col.Null)).
				SetAlign(tview.AlignLeft).
				SetSelectable(true),
		)
		r.columnsTable.SetCell(
			i+1,
			3,
			tview.NewTableCell(col.Key).SetAlign(tview.AlignLeft).SetSelectable(true),
		)
		r.columnsTable.SetCell(
			i+1,
			4,
			tview.NewTableCell(defaultValue).SetAlign(tview.AlignLeft).SetSelectable(true),
		)
		r.columnsTable.SetCell(
			i+1,
			5,
			tview.NewTableCell(col.Extra).SetAlign(tview.AlignLeft).SetSelectable(true),
		)
	}

	r.columnsTable.SetBorder(true)
	r.columnsTable.SetSelectable(true, true)
	r.columnsTable.SetFixed(1, 0)
	r.columnsTable.ScrollToBeginning()
	r.columnsTable.Select(0, 0)

	r.RenderIndexesTable(table)

	return nil
}

func (r *Results) RenderIndexesTable(table string) error {
	indexes, err := r.db.GetIndexes(table)
	if err != nil {
		return fmt.Errorf("Error getting indexes")
	}

	r.indexesTable.Clear()

	// render indexes
	for i, index := range indexes {
		for j, indexColumn := range index {
			r.indexesTable.SetCell(
				i,
				j,
				tview.NewTableCell(indexColumn).SetAlign(tview.AlignLeft).SetSelectable(true),
			)
		}
	}

	r.indexesTable.SetBorder(true)
	r.indexesTable.SetSelectable(true, true)
	r.indexesTable.SetFixed(1, 0)
	r.indexesTable.ScrollToBeginning()

	return nil
}

func (r *Results) renderFilterField() {
	// Handle autocomplete
	r.filter.SetAutocompleteFunc(func(currentText string) (entries []string) {
		if len(currentText) == 0 {
			return
		}

		// Split the currentText into words and use the last word for suggestions
		words := strings.Fields(currentText)
		lastWord := words[len(words)-1]

		// Remove special characters from the last word using regex
		regex := regexp.MustCompile("[^a-zA-Z0-9]+")
		lastWord = regex.ReplaceAllString(lastWord, "")

		// Prepare suggestions
		suggestions := make([]string, len(r.dbColumns))
		for i, col := range r.dbColumns {
			// Add padding so it looks nice on the UI
			suggestions[i] = " " + col.Name + " "
		}

		for _, suggestion := range suggestions {
			if strings.Contains(strings.ToLower(suggestion), strings.ToLower(lastWord)) {
				entries = append(entries, suggestion)
			}
		}

		if len(entries) == 0 {
			entries = nil
		}

		return entries
	})

	r.filter.SetAutocompletedFunc(func(selectedSuggestion string, index int, source int) bool {
		if source != tview.AutocompletedEnter {
			return false
		}

		newText := replaceLastWordWithSuggestion(r.filter.GetText(), selectedSuggestion)
		r.filter.SetText(newText)

		// Return true to indicate the autocompleted text has been handled
		return true
	})

	r.filter.SetLabel("WHERE ").
		SetFieldBackgroundColor(tcell.ColorNone).
		SetDoneFunc(func(key tcell.Key) {
			if key == tcell.KeyEnter {
				where := r.filter.GetText()
				r.RenderTable(r.selectedTable, where)

				r.app.SetFocus(r.resultsTable)
			}
		})
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

func (r *Results) setKeyBindings() {
	// Resutls Table key bindings
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
				r.view.SwitchToPage("columns")
				r.app.SetFocus(r.columnsTable)
			case event.Rune() == 'y':
				// Yank the cell text to clipboard
				row, col := r.resultsTable.GetSelection()
				cell := r.resultsTable.GetCell(row, col)
				clipboard.WriteAll(cell.Text)

				// On yank, Highlight the cell for a short time
				oldBgColor := cell.BackgroundColor
				cell.SetBackgroundColor(tcell.ColorGreen)
				r.resultsTable.SetSelectedStyle(
					tcell.StyleDefault.Background(tcell.ColorYellow).Foreground(tcell.ColorBlack),
				)

				time.AfterFunc(75*time.Millisecond, func() {
					cell.SetBackgroundColor(oldBgColor)
					r.resultsTable.SetSelectedStyle(
						tcell.StyleDefault.Background(tcell.ColorWhite).
							Foreground(tcell.ColorBlack),
					)
					r.app.Draw()
				})
			}
		}

		return event
	})

	r.columnsTable.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyRune {
			switch event.Rune() {
			case '2':
				r.view.SwitchToPage("results")
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
	r.RenderTable(r.selectedTable, "")
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
	r.RenderTable(r.selectedTable, r.filter.GetText())

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

	r.RenderTable(r.selectedTable, r.filter.GetText())
	r.resultsTable.Select(rowToDelete, col)
}

func replaceLastWordWithSuggestion(originalText, suggestion string) string {
	// Handle select autocomplete suggestion:
	// If the user selects a suggestion, replace the last word
	// in the filter text with the selected suggestion.
	splitRegex := regexp.MustCompile(`(\W)|(\w+)`)
	words := splitRegex.FindAllString(originalText, -1)

	// filter out empty strings
	words = words[:0]
	for _, word := range splitRegex.FindAllString(originalText, -1) {
		word = strings.TrimSpace(word)
		if word != "" {
			words = append(words, word)
		}
	}

	// Trim whitespace from the selected suggestion
	suggestion = strings.TrimSpace(suggestion)

	if len(words) == 0 {
		return suggestion
	} else {
		// Replace the last word with the selected suggestion
		words[len(words)-1] = suggestion

		// Join the words back into a string and set it as the new text
		newText := strings.Join(words, " ")
		return newText
	}
}
