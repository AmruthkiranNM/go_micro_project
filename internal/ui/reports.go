package ui

import (
	"inventory-system/internal/services"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func makeReports() fyne.CanvasObject {
	statusLabel := widget.NewLabel("")

	// Summary table
	summaryHeaders := []string{"Product", "Units Sold", "Revenue"}
	summaryTbl := MakeTable(summaryHeaders, nil)
	summaryTbl.SetColumnWidth(0, 250)
	summaryTbl.SetColumnWidth(1, 100)
	summaryTbl.SetColumnWidth(2, 150)

	// History table
	historyHeaders := []string{"#", "Product", "Qty", "Total", "Date"}
	historyTbl := MakeTable(historyHeaders, nil)
	historyTbl.SetColumnWidth(0, 50)
	historyTbl.SetColumnWidth(1, 200)
	historyTbl.SetColumnWidth(2, 70)
	historyTbl.SetColumnWidth(3, 130)
	historyTbl.SetColumnWidth(4, 200)

	revCard := makeLabelCard("💰 Total Revenue", "—")

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
		var sumRows [][]string
		for _, ps := range rd.ProductSummaries {
			sumRows = append(sumRows, []string{ps.Name, IntToStr(ps.TotalQuantity), ps.FormattedRev})
		}
		if len(sumRows) == 0 {
			sumRows = [][]string{{"No data", "—", "—"}}
		}
		newSumTbl := MakeTable(summaryHeaders, sumRows)
		newSumTbl.SetColumnWidth(0, 250)
		newSumTbl.SetColumnWidth(1, 100)
		newSumTbl.SetColumnWidth(2, 150)
		*summaryTbl = *newSumTbl
		summaryTbl.Refresh()

		// Rebuild history table
		var histRows [][]string
		for _, s := range rd.RecentSalesHistory {
			histRows = append(histRows, []string{
				IntToStr(s.ID), s.ProductName, IntToStr(s.Quantity), s.FormattedTotal, s.Date,
			})
		}
		if len(histRows) == 0 {
			histRows = [][]string{{"—", "No sales history", "—", "—", "—"}}
		}
		newHistTbl := MakeTable(historyHeaders, histRows)
		newHistTbl.SetColumnWidth(0, 50)
		newHistTbl.SetColumnWidth(1, 200)
		newHistTbl.SetColumnWidth(2, 70)
		newHistTbl.SetColumnWidth(3, 130)
		newHistTbl.SetColumnWidth(4, 200)
		*historyTbl = *newHistTbl
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
			summaryTbl,
			container.NewVBox(MakeSectionTitle("Recent Sales History (Last 10)"), historyTbl),
		),
	)
}
