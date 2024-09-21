package ui

import (
	"fmt"

	"github.com/alfonzm/lazydb/internal/db"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type CellEditor struct {
	app      *App
	results  *Results
	view     *tview.Flex
	textArea *tview.TextArea
}

func NewCellEditor(
	app *App,
	pages *tview.Pages,
	results *Results,
	db *db.DBClient,
) (*CellEditor, error) {
	textArea := tview.NewTextArea()
	textArea.SetBorder(true).SetTitle("Edit field")

	textArea.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		// On Enter, run the update query
		if event.Key() == tcell.KeyEnter {
			newText := textArea.GetText()

			// update the record in the DB
			selectedRow, selectedColumn := results.resultsTable.GetSelection()
			id := results.resultsTable.GetCell(selectedRow, 0).Text
			colName := results.resultsTable.GetCell(0, selectedColumn).Text

			record := make(map[string]interface{})
			record[colName] = newText

			if err := db.UpdateRecordById(results.selectedTable, id, record); err != nil {
				// TODO: Show error message in the UI
				app.ShowError(fmt.Sprintf("%v", err))
				return event
			}

			// refresh the records table
			results.RenderTable(results.selectedTable, results.filter.GetText())
			pages.HidePage("editor")
			app.SetFocus(results.resultsTable)

			// stay on the same cell
			results.resultsTable.Select(selectedRow, selectedColumn)
		}

		// on press escape, hide the record editor
		if event.Key() == tcell.KeyEscape {
			pages.HidePage("editor")
			app.SetFocus(results.resultsTable)
		}

		return event
	})

	recordEditor := tview.NewFlex().
		AddItem(nil, 0, 1, false).
		AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
			AddItem(nil, 0, 1, false).
			AddItem(textArea, 10, 1, true).
			AddItem(nil, 0, 1, false), 80, 1, true).
		AddItem(nil, 0, 1, false)

	return &CellEditor{
		app:      app,
		view:     recordEditor,
		results:  results,
		textArea: textArea,
	}, nil
}
