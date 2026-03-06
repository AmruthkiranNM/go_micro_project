package ui

import (
	"inventory-system/internal/services"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func makeDashboard() fyne.CanvasObject {
	stats, _ := services.GetDashboardStats()

	// ---- Stat cards ----
	productCard := makeLabelCard("📦 Total Products", IntToStr(stats.ProductCount))
	salesCard := makeLabelCard("🛒 Sales Transactions", IntToStr(stats.SalesCount))
	revenueCard := makeLabelCard("💰 Total Revenue", FormatINR(stats.TotalRevenue))
	lowStockCard := makeLabelCard("⚠️ Low Stock Items", IntToStr(stats.LowStockCount))

	cardsRow := container.NewGridWithColumns(4,
		productCard,
		salesCard,
		revenueCard,
		lowStockCard,
	)

	// ---- Recent Sales Table ----
	headers := []string{"#", "Product", "Qty", "Total", "Date"}
	var rows [][]string

	// Create local slice first (even though this runs synchronously, keeping consistency)
	var newRows [][]string
	for _, s := range stats.RecentSales {
		newRows = append(newRows, []string{
			IntToStr(s.ID),
			s.ProductName,
			IntToStr(s.Quantity),
			s.FormattedTotal,
			s.Date,
		})
	}
	if len(newRows) == 0 {
		newRows = [][]string{{"—", "No recent sales", "—", "—", "—"}}
	}
	rows = newRows

	tbl := MakeTable(headers, &rows)
	tbl.SetColumnWidth(0, 50)
	tbl.SetColumnWidth(1, 200)
	tbl.SetColumnWidth(2, 60)
	tbl.SetColumnWidth(3, 120)
	tbl.SetColumnWidth(4, 180)

	scrollTbl := container.NewScroll(tbl)
	scrollTbl.SetMinSize(fyne.NewSize(700, 200))

	refreshBtn := widget.NewButton("🔄 Refresh", func() {})

	title := MakeSectionTitle("📊 Dashboard")
	recentTitle := MakeSectionTitle("🕒 Recent Sales (Last 5)")

	return container.NewBorder(
		container.NewVBox(title, widget.NewSeparator(), cardsRow, widget.NewSeparator()),
		nil, nil, nil,
		container.NewBorder(
			recentTitle,
			refreshBtn,
			nil, nil,
			scrollTbl,
		),
	)
}

// makeLabelCard produces a simple card-like widget with a title and value.
func makeLabelCard(title, value string) fyne.CanvasObject {
	titleLbl := widget.NewLabelWithStyle(title, fyne.TextAlignCenter, fyne.TextStyle{Bold: false})
	valueLbl := widget.NewLabelWithStyle(value, fyne.TextAlignCenter, fyne.TextStyle{Bold: true})
	return widget.NewCard("", "", container.NewVBox(titleLbl, valueLbl))
}
