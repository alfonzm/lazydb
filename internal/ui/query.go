package ui

import (
	"github.com/alfonzm/lazydb/internal/db"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type Query struct {
	app      *App
	db       *db.DBClient
	textArea *tview.TextArea
	table    *tview.Table
	view     *tview.Flex
}

func NewQuery(
	app *App,
	db *db.DBClient,
) (*Query, error) {
	view := tview.NewFlex()
	view.SetDirection(tview.FlexRow)

	// editor text area at the top and table at the bottom
	textArea := tview.NewTextArea()
	textArea.SetTitle("Query")
	textArea.SetBorder(true)

	table := tview.NewTable()
	table.SetBorder(true)
	table.SetTitle("Results")

	view.AddItem(textArea, 0, 1, true)
	view.AddItem(table, 0, 1, false)

	query := &Query{
		app:      app,
		db:       db,
		table:    table,
		textArea: textArea,
		view:     view,
	}

	query.setKeyBindings()

	return query, nil
}

func (q *Query) setKeyBindings() {
	q.textArea.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		// ctrl+R to run the query
		if event.Key() == tcell.KeyCtrlR {
			// TODO: Run the query
		}

		// escape to clear the text area
		if event.Key() == tcell.KeyEscape {
			q.textArea.SetText("", false)
		}

		return event
	})
}
