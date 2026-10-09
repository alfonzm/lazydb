package ui

import (
	"strings"

	"github.com/alfonzm/lazydb/internal/config"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type CommandModal struct {
	commandModal     tview.Primitive
	commandContainer *tview.Flex
	app              *App
	commands         *tview.List
	textField        *tview.InputField
	lastFocus        tview.Primitive

	connections map[string]config.Connection
}

func NewCommandModal() (*CommandModal, error) {
	textField := tview.NewInputField()
	commands := tview.NewList().
		ShowSecondaryText(false)

	commandModal := tview.NewBox()
	commandContainer := tview.NewFlex()

	commandPalette := &CommandModal{
		commandModal:     commandModal,
		commandContainer: commandContainer,
		commands:         commands,
		textField:        textField,
	}

	commandPalette.setKeyBindings()

	return commandPalette, nil
}

func (cm *CommandModal) Render(currentTab *Tab) {
	cm.lastFocus = cm.app.GetFocus()

	cm.RenderCommands(currentTab)

	// Create UI
	cm.commandContainer = tview.NewFlex().SetDirection(tview.FlexRow)
	cm.commandContainer.SetTitle("Command Palette").SetBorder(true)

	cm.commandContainer.AddItem(cm.textField, 1, 1, true)
	cm.commandContainer.AddItem(cm.commands, 0, 1, false)

	cm.commandModal = CreateModal(cm.commandContainer, 100, 15)

	// Switch to the command page
	cm.app.appPages.RemovePage("command")
	cm.app.appPages.AddPage("command", cm.commandModal, true, false)
	cm.app.appPages.ShowPage("command")

	cm.app.SetFocus(cm.textField)
}

func (cm *CommandModal) RenderCommands(currentTab *Tab) {
	// Prepare connection URL strings
	cm.connections = make(map[string]config.Connection)
	connectionStrings := make([]string, 0)

	cm.commands.Clear()
	filterText := cm.textField.GetText()

	// Sort database names
	for dbName, tables := range currentTab.databaseTables {
		for _, table := range tables {
			if filterText != "" && !strings.Contains(table, filterText) &&
				!strings.Contains(dbName, filterText) {
				continue
			}

			conn := currentTab.connection
			conn.Database = dbName
			conn.Table = table
			connectionString := dbName + "." + table
			cm.connections[connectionString] = conn
			connectionStrings = append(connectionStrings, connectionString)

			cm.commands.AddItem(connectionString, "", 0, func() {
				cm.SelectConnection(cm.connections[connectionString], table)
			})
		}
	}
}

func (cm *CommandModal) setKeyBindings() {
	cm.textField.SetDoneFunc(func(key tcell.Key) {
		if key == tcell.KeyEsc {
			cm.Close()
		}

		if key == tcell.KeyEnter {
			item := cm.commands.GetCurrentItem()
			itemText, _ := cm.commands.GetItemText(item)

			connection, ok := cm.connections[itemText]
			if ok {
				cm.SelectConnection(connection, connection.Table)
			}
		}
	})

	// On each key press, filter the list of commands
	cm.textField.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		// ctrl+p and n move up down
		if event.Key() == tcell.KeyCtrlP {
			cm.commands.SetCurrentItem(cm.commands.GetCurrentItem() - 1)
			return event
		}

		if event.Key() == tcell.KeyCtrlN {
			cm.commands.SetCurrentItem(cm.commands.GetCurrentItem() + 1)
			return event
		}

		// if new text empty, skip
		text := cm.textField.GetText()
		if text == "" {
			return event
		}
		cm.RenderCommands(cm.app.currentTab())
		return event
	})
}

func (cm *CommandModal) Close() {
	cm.app.appPages.SwitchToPage("app")
	cm.app.SetFocus(cm.lastFocus)
}

func (cm *CommandModal) CloseWithoutFocus() {
	cm.app.appPages.SwitchToPage("app")
}

func (cm *CommandModal) SelectConnection(conn config.Connection, table string) {
	cm.app.addNewTab()
	if err := cm.app.currentTab().ConnectDatabase(conn, conn.Database); err != nil {
		return
	}
	cm.app.currentTab().sidebar.SelectTable(table, true)
	cm.CloseWithoutFocus()
}
