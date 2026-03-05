package services

import (
	"inventory-system/database"
	"inventory-system/models"
)

// ProductSummary holds per-product sales aggregates for the reports screen.
type ProductSummary struct {
	Name          string
	TotalQuantity int
	TotalRevenue  float64
	FormattedRev  string
}

// ReportData holds all data needed for the reports screen.
type ReportData struct {
	TotalRevenue       float64
	FormattedRevenue   string
	ProductSummaries   []ProductSummary
	RecentSalesHistory []models.Sale
}

// GetReports queries all report data from the database.
func GetReports() (ReportData, error) {
	var rd ReportData

	database.DB.QueryRow("SELECT COALESCE(SUM(total), 0) FROM sales").Scan(&rd.TotalRevenue)
	rd.FormattedRevenue = models.FormatIndianRupees(rd.TotalRevenue)

	rows, err := database.DB.Query(`
		SELECT p.name, SUM(s.quantity), SUM(s.total)
		FROM sales s
		JOIN products p ON s.product_id = p.id
		GROUP BY p.id
		ORDER BY SUM(s.total) DESC
	`)
	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var ps ProductSummary
			if err := rows.Scan(&ps.Name, &ps.TotalQuantity, &ps.TotalRevenue); err == nil {
				ps.FormattedRev = models.FormatIndianRupees(ps.TotalRevenue)
				rd.ProductSummaries = append(rd.ProductSummaries, ps)
			}
		}
	}

	historyRows, err := database.DB.Query(`
		SELECT s.id, p.name, s.quantity, s.total, s.date
		FROM sales s
		JOIN products p ON s.product_id = p.id
		ORDER BY s.date DESC LIMIT 10
	`)
	if err == nil {
		defer historyRows.Close()
		for historyRows.Next() {
			var s models.Sale
			if err := historyRows.Scan(&s.ID, &s.ProductName, &s.Quantity, &s.Total, &s.Date); err == nil {
				s.FormattedTotal = models.FormatIndianRupees(s.Total)
				rd.RecentSalesHistory = append(rd.RecentSalesHistory, s)
			}
		}
	}

	return rd, nil
}
