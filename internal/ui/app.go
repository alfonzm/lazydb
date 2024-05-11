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
	tab, err := NewTab(app.Application, app.dbClient)
	if err != nil {
		return
	}

	newTabIndex := len(app.tabs)

	app.tabs = append(app.tabs, tab)

	app.tabPages.AddPage(strconv.Itoa(newTabIndex), tab.pages, true, true)
	app.currentTabIndex = newTabIndex
	app.selectTab(app.currentTabIndex)
	app.renderTabHeaders()
}

func (app *App) renderTabHeaders() {
	for i, tab := range app.tabs {
		app.tabHeaders.SetCell(0, i, tview.NewTableCell(tab.name))
	}
}

func (app *App) prevTab() {
	if app.currentTabIndex > 0 {
		app.selectTab(app.currentTabIndex - 1)
	}
}

func (app *App) nextTab() {
	if app.currentTabIndex < len(app.tabs)-1 {
		app.selectTab(app.currentTabIndex + 1)
	}
}

func (app *App) currentTab() *Tab {
	return app.tabs[app.currentTabIndex]
}

func (app *App) selectTab(newTabIndex int) {
	currentTab := app.currentTab()
	if currentTab != nil {
		// currentTab.OnDeactivate()
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

		if event.Key() == tcell.KeyTab {
			switch app.GetFocus() {
			case currentTab.sidebar.list:
				currentTab.Focus("results")
				// currentTab.results.Focus()
			case currentTab.results.resultsTable:
				currentTab.Focus("sidebar")
			case currentTab.results.columnsTable:
				// app.SetFocus(currentTab.results.indexesTable)
				currentTab.Focus("indexes")
			case currentTab.results.indexesTable:
				currentTab.Focus("sidebar")
			}
			return event
		}

		switch event.Key() {
		case tcell.KeyCtrlF:
			app.SetFocus(currentTab.sidebar.list)
			app.SetFocus(currentTab.sidebar.filter)
		case tcell.KeyCtrlR:
			currentTab.results.Focus()
		case tcell.KeyCtrlT:
			app.SetFocus(currentTab.sidebar.list)
		}

		return event
	})
}
