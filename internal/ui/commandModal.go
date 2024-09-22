package ui

import (
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
}

func NewCommandModal() (*CommandModal, error) {
	textField := tview.NewInputField()
	commands := tview.NewList()

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

func (cm *CommandModal) Render() {
	cm.lastFocus = cm.app.GetFocus()

	cm.commandContainer = tview.NewFlex().SetDirection(tview.FlexRow)
	cm.commandContainer.SetTitle("Command Palette").SetBorder(true)

	cm.commandContainer.AddItem(cm.textField, 1, 1, true)
	cm.commandContainer.AddItem(cm.commands, 0, 1, false)

	cm.commandModal = CreateModal(cm.commandContainer, 100, 15)

	cm.app.appPages.RemovePage("command")
	cm.app.appPages.AddPage("command", cm.commandModal, true, false)
	cm.app.appPages.ShowPage("command")

	cm.app.SetFocus(cm.textField)
}

func (cm *CommandModal) setKeyBindings() {
	cm.textField.SetDoneFunc(func(key tcell.Key) {
		if key == tcell.KeyEsc {
			cm.Close()
		}
	})
}

func (cm *CommandModal) Close() {
	cm.app.appPages.SwitchToPage("app")
	cm.app.SetFocus(cm.lastFocus)
}
