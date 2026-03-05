package ui

import (
	"inventory-system/internal/services"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func makeReports() fyne.CanvasObject {
	statusLabel := widget.NewLabel("")

	var sumRows [][]string
	summaryHeaders := []string{"Product", "Units Sold", "Revenue"}
	summaryTbl := MakeTable(summaryHeaders, &sumRows)
	summaryTbl.SetColumnWidth(0, 250)
	summaryTbl.SetColumnWidth(1, 100)
	summaryTbl.SetColumnWidth(2, 150)

	// History table
	var histRows [][]string
	historyHeaders := []string{"#", "Product", "Qty", "Total", "Date"}
	historyTbl := MakeTable(historyHeaders, &histRows)
	historyTbl.SetColumnWidth(0, 50)
	historyTbl.SetColumnWidth(1, 200)
	historyTbl.SetColumnWidth(2, 70)
	historyTbl.SetColumnWidth(3, 130)
	historyTbl.SetColumnWidth(4, 200)

	revCard := container.NewMax(makeLabelCard("💰 Total Revenue", "—"))

	scrollSumTbl := container.NewScroll(summaryTbl)
	scrollHistTbl := container.NewScroll(historyTbl)

	refresh := func() {
		rd, err := services.GetReports()
		if err != nil {
			statusLabel.SetText("❌ " + err.Error())
			return
		}

		// Rebuild revenue card
		newCard := makeLabelCard("💰 Total Revenue", rd.FormattedRevenue)
		revCard.Objects = []fyne.CanvasObject{newCard}
		revCard.Refresh()

		// Rebuild summary table
		var newSumRows [][]string
		for _, ps := range rd.ProductSummaries {
			newSumRows = append(newSumRows, []string{ps.Name, IntToStr(ps.TotalQuantity), ps.FormattedRev})
		}
		if len(newSumRows) == 0 {
			newSumRows = [][]string{{"No data", "—", "—"}}
		}
		sumRows = newSumRows
		summaryTbl.Refresh()

		// Rebuild history table
		var newHistRows [][]string
		for _, s := range rd.RecentSalesHistory {
			newHistRows = append(newHistRows, []string{
				IntToStr(s.ID), s.ProductName, IntToStr(s.Quantity), s.FormattedTotal, s.Date,
			})
		}
		if len(newHistRows) == 0 {
			newHistRows = [][]string{{"—", "No sales history", "—", "—", "—"}}
		}
		histRows = newHistRows
		historyTbl.Refresh()

		statusLabel.SetText("")
	}

	refreshBtn := widget.NewButton("🔄 Refresh", refresh)

	go refresh()

	return container.NewBorder(
		container.NewVBox(
			MakeSectionTitle("📈 Sales Reports"),
			widget.NewSeparator(),
			container.NewHBox(revCard, refreshBtn, statusLabel),
			widget.NewSeparator(),
			MakeSectionTitle("Product-wise Summary"),
		),
		nil, nil, nil,
		container.NewVSplit(
			scrollSumTbl,
			container.NewVBox(MakeSectionTitle("Recent Sales History (Last 10)"), scrollHistTbl),
		),
	)
}
