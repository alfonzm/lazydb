package ui

import (
	"fmt"

	"github.com/alfonzm/lazydb/internal/db"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type Structure struct {
	app          *tview.Application
	db           *db.DBClient
	view         *tview.Flex
	columnsView  *tview.Flex
	columnFilter *tview.InputField
	columnsTable *tview.Table
	indexesTable *tview.Table
}

func NewStructure(
	app *tview.Application,
	db *db.DBClient,
) (*Structure, error) {
	// Setup Columns view
	columnsTable := tview.NewTable()

	columnFilter := tview.NewInputField()
	columnFilter.SetBorder(false)
	columnFilter.SetLabel("FILTER ").
		SetFieldBackgroundColor(tcell.ColorNone).
		SetDoneFunc(func(key tcell.Key) {
			if key == tcell.KeyEnter {
				columnFilterText := columnFilter.GetText()
				panic("implement me - columnFilterText: " + columnFilterText)
			}
		})

	columnsView := tview.NewFlex()
	columnsView.SetBorder(true).
		SetTitle("Structure")
	columnsView.SetDirection(tview.FlexRow).
		AddItem(columnFilter, 1, 1, false).
		AddItem(columnsTable, 0, 1, true)

	// Setup Indexes view
	indexesTable := tview.NewTable()
	indexesTable.SetBorder(false)

	view := tview.NewFlex()
	view.SetBorder(false)
	view.SetDirection(tview.FlexRow)
	view.AddItem(columnsView, 0, 4, true)
	view.AddItem(indexesTable, 0, 1, false)

	structure := &Structure{
		app:          app,
		view:         view,
		db:           db,
		columnsTable: columnsTable,
		indexesTable: indexesTable,
	}

	structure.setKeyBindings()

	return structure, nil
}

func (s *Structure) Render(table string, dbColumns []db.Column) error {
	s.columnsTable.Clear()

	columnMetaColumns := []string{"Field", "Type", "Null", "Key", "Default", "Extra"}

	// set headers from columns
	for i, column := range columnMetaColumns {
		s.columnsTable.SetCell(
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

		s.columnsTable.SetCell(
			i+1,
			0,
			tview.NewTableCell(col.Name).SetAlign(tview.AlignLeft).SetSelectable(true),
		)
		s.columnsTable.SetCell(
			i+1,
			1,
			tview.NewTableCell(col.DataType).SetAlign(tview.AlignLeft).SetSelectable(true),
		)
		s.columnsTable.SetCell(
			i+1,
			2,
			tview.NewTableCell(fmt.Sprintf("%t", col.Null)).
				SetAlign(tview.AlignLeft).
				SetSelectable(true),
		)
		s.columnsTable.SetCell(
			i+1,
			3,
			tview.NewTableCell(col.Key).SetAlign(tview.AlignLeft).SetSelectable(true),
		)
		s.columnsTable.SetCell(
			i+1,
			4,
			tview.NewTableCell(defaultValue).SetAlign(tview.AlignLeft).SetSelectable(true),
		)
		s.columnsTable.SetCell(
			i+1,
			5,
			tview.NewTableCell(col.Extra).SetAlign(tview.AlignLeft).SetSelectable(true),
		)
	}

	s.columnsTable.SetSelectable(true, true)
	s.columnsTable.SetFixed(1, 0)
	s.columnsTable.ScrollToBeginning()
	s.columnsTable.Select(0, 0)

	s.RenderIndexesTable(table)

	return nil
}

func (s *Structure) RenderIndexesTable(table string) error {
	indexes, err := s.db.GetIndexes(table)
	if err != nil {
		return fmt.Errorf("Error getting indexes")
	}

	s.indexesTable.Clear()

	// render indexes
	for i, index := range indexes {
		for j, indexColumn := range index {
			s.indexesTable.SetCell(
				i,
				j,
				tview.NewTableCell(indexColumn).SetAlign(tview.AlignLeft).SetSelectable(true),
			)
		}
	}

	s.indexesTable.SetBorder(true)
	s.indexesTable.SetSelectable(true, true)
	s.indexesTable.SetFixed(1, 0)
	s.indexesTable.ScrollToBeginning()

	return nil
}

func (s *Structure) setKeyBindings() {
	s.app.SetFocus(s.columnsTable)

	s.columnsTable.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyRune {
			switch event.Rune() {
			case '2':
				panic("TODO: go back to results table")
			case '/':
				panic("TODO: Implement filter for columns table")
			}
		}

		return event
	})
}
