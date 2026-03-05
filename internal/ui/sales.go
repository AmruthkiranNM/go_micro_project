package ui

import (
	"fmt"
	"inventory-system/internal/services"
	"inventory-system/models"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

func makeSales() fyne.CanvasObject {
	var sales []models.Sale
	var products []models.Product

	headers := []string{"#", "Product", "Qty", "Total", "Date"}

	var rows [][]string
	tbl := MakeTable(headers, &rows)
	tbl.SetColumnWidth(0, 50)
	tbl.SetColumnWidth(1, 200)
	tbl.SetColumnWidth(2, 70)
	tbl.SetColumnWidth(3, 130)
	tbl.SetColumnWidth(4, 200)

	scrollTbl := container.NewScroll(tbl)
	scrollTbl.SetMinSize(fyne.NewSize(700, 400))

	statusLabel := widget.NewLabel("Loading…")

	refreshTable := func() {
		var err error
		sales, err = services.ListSales()
		if err != nil {
			statusLabel.SetText("❌ " + err.Error())
			return
		}
		products, _ = services.ListProductsInStock()

		var newRows [][]string
		for _, s := range sales {
			newRows = append(newRows, []string{
				IntToStr(s.ID),
				s.ProductName,
				IntToStr(s.Quantity),
				s.FormattedTotal,
				s.Date,
			})
		}
		if len(newRows) == 0 {
			newRows = [][]string{{"—", "No sales recorded yet", "—", "—", "—"}}
		}

		rows = newRows
		tbl.Refresh()
		statusLabel.SetText(fmt.Sprintf("%d sales", len(sales)))
	}

	// ---- Record Sale Dialog ----
	showRecordSaleDialog := func(win fyne.Window) {
		products, _ = services.ListProductsInStock()
		if len(products) == 0 {
			dialog.ShowInformation("No Stock", "All products are out of stock.", win)
			return
		}

		productNames := make([]string, len(products))
		for i, p := range products {
			productNames[i] = fmt.Sprintf("%s (stock: %d)", p.Name, p.Quantity)
		}

		productSelect := widget.NewSelect(productNames, nil)
		productSelect.PlaceHolder = "Select Product"
		qtyEntry := widget.NewEntry()
		qtyEntry.SetPlaceHolder("Quantity to sell")
		errLbl := widget.NewLabel("")

		formItems := []*widget.FormItem{
			{Text: "Product", Widget: productSelect},
			{Text: "Quantity", Widget: qtyEntry},
			{Text: "", Widget: errLbl},
		}

		dlg := dialog.NewForm("Record Sale", "Record", "Cancel", formItems, func(ok bool) {
			if !ok {
				return
			}
			idx := productSelect.SelectedIndex()
			if idx < 0 {
				errLbl.SetText("❌ Please select a product.")
				return
			}
			qty := ParseInt(qtyEntry.Text)
			if qty <= 0 {
				errLbl.SetText("❌ Quantity must be > 0.")
				return
			}
			productID := products[idx].ID
			if err := services.RecordSale(productID, qty); err != nil {
				errLbl.SetText("❌ " + err.Error())
				return
			}
			refreshTable()
		}, win)
		dlg.Resize(fyne.NewSize(450, 280))
		dlg.Show()
	}

	getWin := func() fyne.Window {
		if a := fyne.CurrentApp(); a != nil {
			wins := a.Driver().AllWindows()
			if len(wins) > 0 {
				return wins[0]
			}
		}
		return nil
	}

	addBtn := widget.NewButton("🛒 Record Sale", func() {
		if w := getWin(); w != nil {
			showRecordSaleDialog(w)
		}
	})
	addBtn.Importance = widget.HighImportance

	refreshBtn := widget.NewButton("🔄 Refresh", func() { refreshTable() })

	toolbar := container.NewHBox(addBtn, refreshBtn, statusLabel)

	// Refresh UI load synchronously
	refreshTable()

	return container.NewBorder(
		container.NewVBox(MakeSectionTitle("📈 Sales Tracking"), widget.NewSeparator(), toolbar),
		nil, nil, nil,
		scrollTbl,
	)
}
