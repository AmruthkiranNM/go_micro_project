package handlers

import (
	"net/http"

	"inventory-system/database"
	"inventory-system/models"

	"github.com/gin-gonic/gin"
)

func ShowReports(c *gin.Context) {
	// Total Revenue
	var totalRevenue float64
	database.DB.QueryRow("SELECT COALESCE(SUM(total), 0) FROM sales").Scan(&totalRevenue)

	// Product-wise sales summary
	rows, err := database.DB.Query(`
		SELECT p.name, SUM(s.quantity), SUM(s.total)
		FROM sales s
		JOIN products p ON s.product_id = p.id
		GROUP BY p.id
	`)

	type ProductSummary struct {
		Name          string
		TotalQuantity int
		TotalRevenue  float64
	}

	var summaries []ProductSummary
	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var ps ProductSummary
			if err := rows.Scan(&ps.Name, &ps.TotalQuantity, &ps.TotalRevenue); err == nil {
				summaries = append(summaries, ps)
			}
		}
	}

	// Recent Sales History (already in sales page, but can be added here too)
	historyRows, _ := database.DB.Query(`
		SELECT s.id, p.name, s.quantity, s.total, s.date 
		FROM sales s 
		JOIN products p ON s.product_id = p.id
		ORDER BY s.date DESC LIMIT 10
	`)
	var history []models.Sale
	if historyRows != nil {
		defer historyRows.Close()
		for historyRows.Next() {
			var s models.Sale
			if err := historyRows.Scan(&s.ID, &s.ProductName, &s.Quantity, &s.Total, &s.Date); err == nil {
				history = append(history, s)
			}
		}
	}

	c.HTML(http.StatusOK, "reports.html", gin.H{
		"title":        "Reports",
		"totalRevenue": totalRevenue,
		"summaries":    summaries,
		"history":      history,
	})
}
