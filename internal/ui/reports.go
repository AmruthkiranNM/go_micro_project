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

	revCard := makeLabelCard("💰 Total Revenue", "—")

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
		*revCard.(*widget.Card) = *newCard.(*widget.Card)
		revCard.Refresh()

		// Rebuild summary table
		sumRows = nil
		for _, ps := range rd.ProductSummaries {
			sumRows = append(sumRows, []string{ps.Name, IntToStr(ps.TotalQuantity), ps.FormattedRev})
		}
		if len(sumRows) == 0 {
			sumRows = [][]string{{"No data", "—", "—"}}
		}
		summaryTbl.Refresh()

		// Rebuild history table
		histRows = nil
		for _, s := range rd.RecentSalesHistory {
			histRows = append(histRows, []string{
				IntToStr(s.ID), s.ProductName, IntToStr(s.Quantity), s.FormattedTotal, s.Date,
			})
		}
		if len(histRows) == 0 {
			histRows = [][]string{{"—", "No sales history", "—", "—", "—"}}
		}
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
