package ui

import (
	"inventory-system/models"
	"strconv"

	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

// FormatINR wraps the models helper for use across UI files.
func FormatINR(amount float64) string {
	return models.FormatIndianRupees(amount)
}

// MakeStatCard creates a styled stat card with a label and a bold value.
func MakeStatCard(label, value string, iconName fyne.ThemeIconName) fyne.CanvasObject {
	icon := widget.NewIcon(theme.Icon(iconName))
	valueLabel := canvas.NewText(value, color.White)
	valueLabel.TextSize = 22
	valueLabel.TextStyle = fyne.TextStyle{Bold: true}

	titleLabel := canvas.NewText(label, color.NRGBA{R: 200, G: 200, B: 200, A: 255})
	titleLabel.TextSize = 12

	bg := canvas.NewRectangle(color.NRGBA{R: 30, G: 35, B: 50, A: 255})

	content := container.NewVBox(
		icon,
		valueLabel,
		titleLabel,
	)
	return container.NewStack(bg, container.NewPadded(content))
}

// MakeSectionTitle returns a bold section heading.
func MakeSectionTitle(title string) fyne.CanvasObject {
	t := canvas.NewText(title, color.White)
	t.TextSize = 16
	t.TextStyle = fyne.TextStyle{Bold: true}
	return t
}

// MakeErrorLabel creates a red error label (initially hidden).
func MakeErrorLabel() *widget.Label {
	lbl := widget.NewLabel("")
	lbl.Wrapping = fyne.TextWrapWord
	return lbl
}

// MakeTable creates a simple table from headers and a pointer to rows of strings.
// rows is a *[][]string so it can be updated dynamically and refreshed.
func MakeTable(headers []string, rows *[][]string) *widget.Table {
	t := widget.NewTable(
		func() (int, int) {
			if rows == nil || len(*rows) == 0 {
				return 1, len(headers)
			}
			return len(*rows) + 1, len(headers)
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("                              ")
		},
		func(id widget.TableCellID, obj fyne.CanvasObject) {
			lbl := obj.(*widget.Label)
			if id.Row == 0 {
				lbl.TextStyle = fyne.TextStyle{Bold: true}
				lbl.SetText(headers[id.Col])
			} else {
				lbl.TextStyle = fyne.TextStyle{}
				if rows != nil && id.Row-1 < len(*rows) && id.Col < len((*rows)[id.Row-1]) {
					lbl.SetText((*rows)[id.Row-1][id.Col])
				} else {
					lbl.SetText("")
				}
			}
		},
	)
	return t
}

// IntToStr converts int to string.
func IntToStr(i int) string {
	return strconv.Itoa(i)
}

// FloatToStr converts float64 to string (2 decimals).
func FloatToStr(f float64) string {
	return strconv.FormatFloat(f, 'f', 2, 64)
}

// ParseFloat parses a string to float64, returns 0 on error.
func ParseFloat(s string) float64 {
	v, _ := strconv.ParseFloat(s, 64)
	return v
}

// ParseInt parses a string to int, returns 0 on error.
func ParseInt(s string) int {
	v, _ := strconv.Atoi(s)
	return v
}
