package ui

import (
	"strconv"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type App struct {
	*tview.Application
	appPages        *tview.Pages
	appContainer    *tview.Flex
	tabHeaders      *tview.Table
	tabPages        *tview.Pages
	tabs            []*Tab
	currentTabIndex int
	errorModal      *ErrorModal
	commandModal    *CommandModal
}

func Start() error {
	appPages := tview.NewPages()
	container := tview.NewFlex().SetDirection(tview.FlexRow)
	tabHeaders := tview.NewTable().SetSelectable(false, true)
	tabPages := tview.NewPages()
	application := tview.NewApplication()
	errorModal, err := NewErrorModal()
	commandModal, err := NewCommandModal()
	if err != nil {
		return err
	}

	app := &App{
		Application:  application,
		appPages:     appPages,
		appContainer: container,
		tabHeaders:   tabHeaders,
		tabPages:     tabPages,
		errorModal:   errorModal,
		commandModal: commandModal,
	}

	// TODO: Not sure if this is the best way to pass the app to the modals
	errorModal.app = app
	commandModal.app = app

	app.addNewTab()

	container.AddItem(tabHeaders, 1, 0, false)
	container.AddItem(tabPages, 0, 1, true)

	appPages.AddPage("app", container, true, true)

	app.setKeyBindings()

	if err := app.SetRoot(appPages, true).Run(); err != nil {
		return err
	}
	return nil
}

func (app *App) addNewTab() {
	currentTab := app.currentTab()
	if currentTab != nil {
		currentTab.OnDeactivate()
	}

	tab, err := NewTab(app)
	if err != nil {
		return
	}

	newTabIndex := len(app.tabs)

	app.tabs = append(app.tabs, tab)

	app.tabPages.AddPage(strconv.Itoa(newTabIndex), tab.pages, true, true)
	app.selectTab(newTabIndex)
	app.RenderTabHeaders()
}

func (app *App) RenderTabHeaders() {
	for i, tab := range app.tabs {
		app.tabHeaders.SetCell(0, i, tview.NewTableCell(tab.name))
	}
}

func (app *App) onDeactivateCurrentTab() {
	currentTab := app.currentTab()
	if currentTab != nil {
		currentTab.OnDeactivate()
	}
}

func (app *App) prevTab() {
	targetTabIndex := app.currentTabIndex - 1

	if targetTabIndex < 0 {
		targetTabIndex = len(app.tabs) - 1
	}

	app.onDeactivateCurrentTab()
	app.selectTab(targetTabIndex)
}

func (app *App) nextTab() {
	targetTabIndex := app.currentTabIndex + 1

	if targetTabIndex >= len(app.tabs) {
		targetTabIndex = 0
	}

	app.onDeactivateCurrentTab()
	app.selectTab(targetTabIndex)
}

func (app *App) currentTab() *Tab {
	if len(app.tabs) == 0 {
		return nil
	}
	return app.tabs[app.currentTabIndex]
}

func (app *App) selectTab(newTabIndex int) {
	if newTabIndex < 0 || newTabIndex >= len(app.tabs) {
		return
	}

	app.currentTabIndex = newTabIndex

	app.tabHeaders.Select(0, app.currentTabIndex)
	app.tabPages.SwitchToPage(strconv.Itoa(app.currentTabIndex))

	newSelectedTab := app.currentTab()
	newSelectedTab.OnActivate()
}

func (app *App) setKeyBindings() {
	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		currentTab := app.currentTab()

		// If the focus is on an input/textarea field, early return
		if _, ok := app.GetFocus().(*tview.InputField); ok {
			return event
		}

		if _, ok := app.GetFocus().(*tview.TextArea); ok {
			return event
		}

		if event.Key() == tcell.KeyRune {
			switch event.Rune() {

			// Tab management
			case '[':
				app.prevTab()
			case ']':
				app.nextTab()
			case 't':
				app.addNewTab()

				// App management
			case 'q':
				app.Stop()

			// Current tab hotkeys
			case '0':
				currentTab.pages.SwitchToPage("connections")
			}
		}

		switch event.Key() {
		case tcell.KeyCtrlF:
			currentTab.FocusFindTable()
		case tcell.KeyTab:
			currentTab.OnPressTab()
		case tcell.KeyCtrlP:
			app.commandModal.Render()
		}

		return event
	})
}

func (app *App) ShowError(errorText string) {
	app.errorModal.RenderError(errorText)
}

/* https://github.com/rivo/tview/wiki/CreateModal */
// Create modal container centered on screen
func CreateModal(p tview.Primitive, width, height int) tview.Primitive {
	return tview.NewFlex().
		AddItem(nil, 0, 1, false).
		AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
			AddItem(nil, 0, 1, false).
			AddItem(p, height, 1, true).
			AddItem(nil, 0, 1, false), width, 1, true).
		AddItem(nil, 0, 1, false)
}
