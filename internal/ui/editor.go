package ui

import (
	"github.com/alfonzm/lazydb/internal/db"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type Editor struct {
	app     *tview.Application
	results *Results
	view    *tview.Flex
}

func NewEditor(
	app *tview.Application,
	pages *tview.Pages,
	results *Results,
	db *db.DBClient,
) (*Editor, error) {
	textArea := tview.NewTextArea()
	textArea.SetBorder(true).SetTitle("Edit field")

	textArea.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		// On Ctrl+Enter, run the query
		if event.Rune() == 13 {
			newText := textArea.GetText()

			// update the record in the DB
			selectedRow, selectedColumn := results.table.GetSelection()
			id := results.table.GetCell(selectedRow, 0).Text
			colName := results.table.GetCell(0, selectedColumn).Text

			record := make(map[string]interface{})
			record[colName] = newText

			if err := db.UpdateRecordById(results.selectedTable, id, record); err != nil {
				return event
			}

			// refresh the records table
			results.RenderTable(results.selectedTable, "")
			pages.HidePage("editor")
			app.SetFocus(results.table)

			// stay on the same cell
			results.table.Select(selectedRow, selectedColumn)
		}

		// on press escape, hide the record editor
		if event.Key() == tcell.KeyEscape {
			pages.HidePage("editor")
			app.SetFocus(results.table)
		}

		return event
	})

	recordEditor := tview.NewFlex().
		AddItem(nil, 0, 1, false).
		AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
			AddItem(nil, 0, 1, false).
			AddItem(textArea, 15, 1, true).
			AddItem(nil, 0, 1, false), 100, 1, true).
		AddItem(nil, 0, 1, false)

	return &Editor{
		app:     app,
		view:    recordEditor,
		results: results,
	}, nil
}
