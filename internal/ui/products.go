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

func makeProducts() fyne.CanvasObject {
	var products []models.Product

	headers := []string{"ID", "Name", "Category", "Price", "Qty", "Status"}

	// Table data builder
	buildRows := func() [][]string {
		var rows [][]string
		for _, p := range products {
			rows = append(rows, []string{
				IntToStr(p.ID),
				p.Name,
				p.Category,
				p.FormattedPrice,
				IntToStr(p.Quantity),
				p.Status,
			})
		}
		return rows
	}

	var rows [][]string
	tbl := MakeTable(headers, &rows)
	tbl.SetColumnWidth(0, 50)
	tbl.SetColumnWidth(1, 200)
	tbl.SetColumnWidth(2, 150)
	tbl.SetColumnWidth(3, 120)
	tbl.SetColumnWidth(4, 70)
	tbl.SetColumnWidth(5, 120)

	scrollTbl := container.NewScroll(tbl)
	scrollTbl.SetMinSize(fyne.NewSize(700, 400))

	statusLabel := widget.NewLabel("Loading products…")

	// Rebuild table after data changes
	refreshTable := func() {
		var err error
		products, err = services.ListProducts()
		if err != nil {
			statusLabel.SetText("❌ " + err.Error())
			return
		}
		newRows := buildRows()
		if len(newRows) == 0 {
			newRows = [][]string{{"—", "No products found", "—", "—", "—", "—"}}
		}
		rows = newRows
		tbl.Refresh()
		statusLabel.SetText(fmt.Sprintf("%d products", len(products)))
	}

	// ---- Add Product Dialog ----
	showAddDialog := func(win fyne.Window) {
		nameEntry := widget.NewEntry()
		nameEntry.SetPlaceHolder("Product Name")
		catEntry := widget.NewEntry()
		catEntry.SetPlaceHolder("Category")
		priceEntry := widget.NewEntry()
		priceEntry.SetPlaceHolder("Price (e.g. 299.99)")
		qtyEntry := widget.NewEntry()
		qtyEntry.SetPlaceHolder("Quantity")
		errLbl := widget.NewLabel("")

		formItems := []*widget.FormItem{
			{Text: "Name", Widget: nameEntry},
			{Text: "Category", Widget: catEntry},
			{Text: "Price (₹)", Widget: priceEntry},
			{Text: "Quantity", Widget: qtyEntry},
			{Text: "", Widget: errLbl},
		}

		dlg := dialog.NewForm("Add Product", "Add", "Cancel", formItems, func(ok bool) {
			if !ok {
				return
			}
			p := models.Product{
				Name:     nameEntry.Text,
				Category: catEntry.Text,
				Price:    ParseFloat(priceEntry.Text),
				Quantity: ParseInt(qtyEntry.Text),
			}
			if err := services.AddProduct(p); err != nil {
				errLbl.SetText("❌ " + err.Error())
				return
			}
			refreshTable()
		}, win)
		dlg.Resize(fyne.NewSize(420, 300))
		dlg.Show()
	}

	// ---- Edit Product Dialog ----
	var selectedRow int = -1
	tbl.OnSelected = func(id widget.TableCellID) {
		selectedRow = id.Row - 1 // row 0 is header
	}

	showEditDialog := func(win fyne.Window) {
		if selectedRow < 0 || selectedRow >= len(products) {
			dialog.ShowInformation("Select Row", "Please click on a product row first.", win)
			return
		}
		p := products[selectedRow]
		nameEntry := widget.NewEntry()
		nameEntry.SetText(p.Name)
		catEntry := widget.NewEntry()
		catEntry.SetText(p.Category)
		priceEntry := widget.NewEntry()
		priceEntry.SetText(FloatToStr(p.Price))
		qtyEntry := widget.NewEntry()
		qtyEntry.SetText(IntToStr(p.Quantity))
		errLbl := widget.NewLabel("")

		formItems := []*widget.FormItem{
			{Text: "Name", Widget: nameEntry},
			{Text: "Category", Widget: catEntry},
			{Text: "Price (₹)", Widget: priceEntry},
			{Text: "Quantity", Widget: qtyEntry},
			{Text: "", Widget: errLbl},
		}

		dlg := dialog.NewForm("Edit Product", "Save", "Cancel", formItems, func(ok bool) {
			if !ok {
				return
			}
			updated := models.Product{
				Name:     nameEntry.Text,
				Category: catEntry.Text,
				Price:    ParseFloat(priceEntry.Text),
				Quantity: ParseInt(qtyEntry.Text),
			}
			if err := services.UpdateProduct(p.ID, updated); err != nil {
				errLbl.SetText("❌ " + err.Error())
				return
			}
			selectedRow = -1
			refreshTable()
		}, win)
		dlg.Resize(fyne.NewSize(420, 300))
		dlg.Show()
	}

	// ---- Delete Product ----
	showDeleteConfirm := func(win fyne.Window) {
		if selectedRow < 0 || selectedRow >= len(products) {
			dialog.ShowInformation("Select Row", "Please click on a product row first.", win)
			return
		}
		p := products[selectedRow]
		dialog.ShowConfirm(
			"Delete Product",
			fmt.Sprintf("Delete '%s'? This cannot be undone.", p.Name),
			func(ok bool) {
				if !ok {
					return
				}
				if err := services.DeleteProduct(p.ID); err != nil {
					dialog.ShowError(fmt.Errorf("cannot delete: %w", err), win)
					return
				}
				selectedRow = -1
				refreshTable()
			},
			win,
		)
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

	addBtn := widget.NewButton("➕ Add Product", func() {
		if w := getWin(); w != nil {
			showAddDialog(w)
		}
	})
	addBtn.Importance = widget.HighImportance

	editBtn := widget.NewButton("✏️ Edit Selected", func() {
		if w := getWin(); w != nil {
			showEditDialog(w)
		}
	})

	deleteBtn := widget.NewButton("🗑️ Delete Selected", func() {
		if w := getWin(); w != nil {
			showDeleteConfirm(w)
		}
	})
	deleteBtn.Importance = widget.DangerImportance

	refreshBtn := widget.NewButton("🔄 Refresh", func() { refreshTable() })

	toolbar := container.NewHBox(addBtn, editBtn, deleteBtn, refreshBtn, statusLabel)

	screenContent := container.NewBorder(
		container.NewVBox(MakeSectionTitle("📦 Product Management"), widget.NewSeparator(), toolbar),
		nil, nil, nil,
		scrollTbl,
	)

	// Sync refresh load
	refreshTable()

	return screenContent
}
