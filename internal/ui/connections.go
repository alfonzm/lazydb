package ui

import (
	"sort"

	"github.com/alfonzm/lazydb/internal/config"
	"github.com/alfonzm/lazydb/internal/db"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type Connections struct {
	view *tview.Flex
	list *tview.List
	tab  *Tab
}

func NewConnections(
	tab *Tab,
	db *db.DBClient,
) (*Connections, error) {
	connConfigurations, err := config.GetConnections()
	if err != nil {
		return nil, err
	}

	list := tview.NewList()
	view := tview.NewFlex()

	connections := &Connections{
		tab:  tab,
		view: view,
		list: list,
	}

	list.SetBorder(true)
	list.SetTitle("Select a connection")
	list.ShowSecondaryText(false)

	var connNames []string
	for k := range connConfigurations {
		connNames = append(connNames, k)
	}
	sort.Strings(connNames)

	for _, dbName := range connNames {
		conn := connConfigurations[dbName]
		list.AddItem(conn.Database, "", 0, connections.selectConnection(conn.String(), dbName))
	}

	view.SetDirection(tview.FlexRow).
		AddItem(list, 0, 1, true)

	connections.setKeyBindings()

	return connections, nil
}

func (c *Connections) setKeyBindings() {
	c.view.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyRune {
			switch event.Rune() {
			case 'j':
				// pressing j at the end of the list goes to the top
				if c.list.GetItemCount()-1 == c.list.GetCurrentItem() {
					c.list.SetCurrentItem(0)
				} else {
					c.list.SetCurrentItem(c.list.GetCurrentItem() + 1)
				}
			case 'k':
				c.list.SetCurrentItem(c.list.GetCurrentItem() - 1)
			}
		}
		return event
	})
}

func (c *Connections) selectConnection(url string, dbName string) func() {
	return func() {
		c.tab.ConnectDatabase(url, dbName)
	}
}
