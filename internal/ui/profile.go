package ui

import (
	"inventory-system/internal/services"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

func makeProfile(state *AppState) fyne.CanvasObject {
	if state.CurrentUser == nil {
		return widget.NewLabel("Not logged in.")
	}

	user, err := services.GetUserByID(state.CurrentUser.ID)
	if err != nil {
		return widget.NewLabel("Failed to load profile: " + err.Error())
	}

	usernameEntry := widget.NewEntry()
	usernameEntry.SetText(user.Username)

	emailEntry := widget.NewEntry()
	emailEntry.SetText(user.Email)

	passwordEntry := widget.NewPasswordEntry()
	passwordEntry.SetPlaceHolder("Leave blank to keep current password")

	statusLabel := widget.NewLabel("")
	statusLabel.Wrapping = fyne.TextWrapWord

	saveBtn := widget.NewButton("💾 Save Changes", nil)
	saveBtn.Importance = widget.HighImportance

	var winRef fyne.Window
	saveBtn.OnTapped = func() {
		err := services.UpdateProfile(
			state.CurrentUser.ID,
			usernameEntry.Text,
			emailEntry.Text,
			passwordEntry.Text,
		)
		if err != nil {
			statusLabel.SetText("❌ " + err.Error())
			return
		}
		// Update cached username in state
		state.CurrentUser.Username = usernameEntry.Text
		state.CurrentUser.Email = emailEntry.Text
		passwordEntry.SetText("")
		statusLabel.SetText("✅ Profile updated successfully!")
		if winRef != nil {
			dialog.ShowInformation("Success", "Profile updated successfully!", winRef)
		}
	}

	go func() {
		if a := fyne.CurrentApp(); a != nil {
			windows := a.Driver().AllWindows()
			if len(windows) > 0 {
				winRef = windows[0]
			}
		}
	}()

	form := widget.NewForm(
		widget.NewFormItem("Username", usernameEntry),
		widget.NewFormItem("Email", emailEntry),
		widget.NewFormItem("New Password", passwordEntry),
	)

	scrollContent := container.NewScroll(
		container.NewPadded(
			container.NewVBox(
				form,
				widget.NewLabel(""), // Spacer
				saveBtn,
				statusLabel,
			),
		),
	)

	return container.NewBorder(
		container.NewVBox(
			MakeSectionTitle("👤 Admin Profile Settings"),
			widget.NewSeparator(),
		),
		nil, nil, nil,
		scrollContent,
	)
}
