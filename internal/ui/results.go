package ui

import (
	"fmt"

	"github.com/alfonzm/lazydb/internal/db"
	"github.com/rivo/tview"
)

type Results struct {
	view  *tview.Flex
	table *tview.Table
	db    *db.DBClient
}

func NewResults(db *db.DBClient) (*Results, error) {
	table := tview.NewTable()

	view := tview.NewFlex()
	view.SetBorder(true).
		SetTitle("Results").
		SetBorder(true)
	view.SetDirection(tview.FlexRow).
		AddItem(nil, 1, 1, false).
		AddItem(table, 0, 1, false)

	return &Results{
		table: table,
		view:  view,
		db:    db,
	}, nil
}

func (r *Results) RenderTable(table string) error {
	dbColumns, err := r.db.GetColumns(table)
	if err != nil {
		return fmt.Errorf("Error getting columns")
	}

	dbRecords, err := r.db.GetRecords(table)
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
		// edit field
	})

	// Iterate over records and fill table
	for rowIndex, record := range dbRecords {
		for columnIndex, colName := range dbColumns {
			recordValue, ok := record[colName]

			cellString := ""

			// if DB value is null, set valStr to "NULL"
			if ok && recordValue == nil {
				// TODO: For some reason this is not working
				// Maybe there is a better way to check for NULL DB values
				cellString = "NULL"
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
