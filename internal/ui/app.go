package ui

import (
	"strconv"

	"github.com/alfonzm/lazydb/internal/db"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type App struct {
	*tview.Application
	container       *tview.Flex
	tabHeaders      *tview.Table
	tabPages        *tview.Pages
	tabs            []*Tab
	currentTabIndex int
	dbClient        *db.DBClient
}

func Start() error {
	container := tview.NewFlex().SetDirection(tview.FlexRow)
	tabHeaders := tview.NewTable().SetSelectable(false, true)
	tabPages := tview.NewPages()

	app := &App{
		Application: tview.NewApplication(),
		container:   container,
		tabHeaders:  tabHeaders,
		tabPages:    tabPages,
	}

	app.addNewTab()

	container.AddItem(tabHeaders, 1, 0, false)
	container.AddItem(tabPages, 0, 1, true)

	app.setKeyBindings()

	if err := app.SetRoot(container, true).Run(); err != nil {
		return err
	}
	return nil
}

func (app *App) addNewTab() {
	currentTab := app.currentTab()
	if currentTab != nil {
		currentTab.OnDeactivate()
	}

	tab, err := NewTab(app, app.dbClient)
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
		}

		return event
	})
}
