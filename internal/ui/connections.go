package ui

import (
	"github.com/alfonzm/lazydb/internal/db"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type Connections struct {
	app  *App
	view *tview.Flex
	list *tview.List
}

type Connection struct {
	name string
	url  string
}

func NewConnections(
	app *App,
	db *db.DBClient,
) (*Connections, error) {
	// TODO: Implement reading connections from a configuration file
	connConfigurations := []Connection{
		{
			name: "foo",
			url:  "root:root@/foo",
		},
	}

	list := tview.NewList()
	view := tview.NewFlex()

	connections := &Connections{
		app:  app,
		view: view,
		list: list,
	}

	list.SetBorder(true)
	list.SetTitle("Select a connection")
	list.ShowSecondaryText(false)

	for _, conn := range connConfigurations {
		list.AddItem(conn.name, "", 0, connections.selectConnection(conn))
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

func (c *Connections) selectConnection(conn Connection) func() {
	return func() {
		c.app.Connect(conn)
	}
}
