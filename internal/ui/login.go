package ui

import (
	"inventory-system/internal/services"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

// showLogin renders the login screen and sets it as the window content.
func showLogin(state *AppState) {
	title := widget.NewLabelWithStyle(
		"🏪 Inventory & Sales Management",
		fyne.TextAlignCenter,
		fyne.TextStyle{Bold: true},
	)

	subtitle := widget.NewLabelWithStyle(
		"Please sign in to continue",
		fyne.TextAlignCenter,
		fyne.TextStyle{Italic: true},
	)

	usernameEntry := widget.NewEntry()
	usernameEntry.SetPlaceHolder("Username")

	passwordEntry := widget.NewPasswordEntry()
	passwordEntry.SetPlaceHolder("Password")

	errorLabel := widget.NewLabel("")
	errorLabel.Wrapping = fyne.TextWrapWord

	loginBtn := widget.NewButton("Login", nil)
	loginBtn.Importance = widget.HighImportance

	doLogin := func() {
		username := usernameEntry.Text
		password := passwordEntry.Text
		if username == "" || password == "" {
			errorLabel.SetText("⚠️ Please enter username and password.")
			return
		}
		user, err := services.AuthenticateUser(username, password)
		if err != nil {
			errorLabel.SetText("❌ " + err.Error())
			passwordEntry.SetText("")
			return
		}
		state.CurrentUser = user
		errorLabel.SetText("")
		showMainApp(state)
	}

	loginBtn.OnTapped = func() { doLogin() }

	// Allow Enter key in password field to submit
	passwordEntry.OnSubmitted = func(_ string) { doLogin() }
	usernameEntry.OnSubmitted = func(_ string) {
		state.Window.Canvas().Focus(passwordEntry)
	}

	form := container.NewVBox(
		widget.NewSeparator(),
		usernameEntry,
		passwordEntry,
		errorLabel,
		loginBtn,
	)

	centered := container.NewCenter(
		container.NewVBox(
			title,
			subtitle,
			widget.NewSeparator(),
			form,
		),
	)

	state.Window.SetContent(centered)
	state.Window.Canvas().Focus(usernameEntry)
}
