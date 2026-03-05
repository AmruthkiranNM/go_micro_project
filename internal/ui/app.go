package ui

import (
	"inventory-system/models"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

// AppState holds shared state across all UI screens.
type AppState struct {
	App         fyne.App
	Window      fyne.Window
	CurrentUser *models.User
}

// Run initialises the Fyne app and shows the login screen.
func Run() {
	a := app.New()
	a.Settings().SetTheme(theme.DarkTheme())

	w := a.NewWindow("Inventory & Sales Management System")
	w.Resize(fyne.NewSize(1100, 700))
	w.CenterOnScreen()

	state := &AppState{
		App:    a,
		Window: w,
	}

	showLogin(state)
	w.ShowAndRun()
}

// showMainApp builds the main navigation shell after successful login.
func showMainApp(state *AppState) {
	// Navigation sidebar items
	navItems := []string{
		"📊  Dashboard",
		"📦  Products",
		"🛒  Sales",
		"📈  Reports",
		"👤  Profile",
	}

	// Content area placeholder
	contentArea := container.NewStack()

	// Screen builders
	screens := map[string]func() fyne.CanvasObject{
		"📊  Dashboard": func() fyne.CanvasObject { return makeDashboard() },
		"📦  Products":  func() fyne.CanvasObject { return makeProducts() },
		"🛒  Sales":     func() fyne.CanvasObject { return makeSales() },
		"📈  Reports":   func() fyne.CanvasObject { return makeReports() },
		"👤  Profile":   func() fyne.CanvasObject { return makeProfile(state) },
	}

	var navList *widget.List
	navList = widget.NewList(
		func() int { return len(navItems) },
		func() fyne.CanvasObject {
			return widget.NewLabel("                    ")
		},
		func(id widget.ListItemID, obj fyne.CanvasObject) {
			obj.(*widget.Label).SetText(navItems[id])
		},
	)

	// Select Dashboard by default
	navList.OnSelected = func(id widget.ListItemID) {
		key := navItems[id]
		if builder, ok := screens[key]; ok {
			contentArea.Objects = []fyne.CanvasObject{builder()}
			contentArea.Refresh()
		}
	}

	// Logout button
	logoutBtn := widget.NewButton("🚪 Logout", func() {
		state.CurrentUser = nil
		showLogin(state)
	})

	userLabel := widget.NewLabel("👤 " + state.CurrentUser.Username)
	userLabel.TextStyle = fyne.TextStyle{Bold: true}

	sidebar := container.NewBorder(
		container.NewVBox(
			widget.NewLabelWithStyle("MENU", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
			widget.NewSeparator(),
		),
		container.NewVBox(widget.NewSeparator(), userLabel, logoutBtn),
		nil, nil,
		navList,
	)

	split := container.NewHSplit(sidebar, contentArea)
	split.SetOffset(0.22)

	state.Window.SetContent(split)

	// Load dashboard by default
	navList.Select(0)
}
