package services

import (
	"database/sql"
	"errors"
	"inventory-system/database"
	"inventory-system/models"
	"time"
)

// ListSales returns all sales joined with product names, newest first.
func ListSales() ([]models.Sale, error) {
	rows, err := database.DB.Query(`
		SELECT s.id, p.name, s.quantity, s.total, s.date
		FROM sales s
		JOIN products p ON s.product_id = p.id
		ORDER BY s.date DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var sales []models.Sale
	for rows.Next() {
		var s models.Sale
		if err := rows.Scan(&s.ID, &s.ProductName, &s.Quantity, &s.Total, &s.Date); err != nil {
			continue
		}
		s.FormattedTotal = models.FormatIndianRupees(s.Total)
		sales = append(sales, s)
	}
	return sales, rows.Err()
}

// RecordSale records a sale transaction and deducts stock atomically.
func RecordSale(productID, quantity int) error {
	if quantity <= 0 {
		return errors.New("quantity must be greater than 0")
	}

	tx, err := database.DB.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	var price float64
	var currentQty int
	err = tx.QueryRow("SELECT price, quantity FROM products WHERE id = ?", productID).
		Scan(&price, &currentQty)
	if err != nil {
		if err == sql.ErrNoRows {
			return errors.New("product not found")
		}
		return err
	}

	if currentQty < quantity {
		return errors.New("insufficient stock")
	}

	total := price * float64(quantity)
	date := time.Now().Format("2006-01-02 15:04:05")

	if _, err = tx.Exec(
		"INSERT INTO sales (product_id, quantity, total, date) VALUES (?, ?, ?, ?)",
		productID, quantity, total, date,
	); err != nil {
		return err
	}

	if _, err = tx.Exec(
		"UPDATE products SET quantity = quantity - ? WHERE id = ?", quantity, productID,
	); err != nil {
		return err
	}

	return tx.Commit()
}

// DashboardStats holds the stats shown on the dashboard.
type DashboardStats struct {
	ProductCount  int
	SalesCount    int
	LowStockCount int
	TotalRevenue  float64
	RecentSales   []models.Sale
}

// GetDashboardStats queries the DB for all dashboard statistics.
func GetDashboardStats() (DashboardStats, error) {
	var s DashboardStats

	database.DB.QueryRow("SELECT COUNT(*) FROM products").Scan(&s.ProductCount)
	database.DB.QueryRow("SELECT COUNT(*) FROM sales").Scan(&s.SalesCount)
	database.DB.QueryRow("SELECT COALESCE(SUM(total), 0) FROM sales").Scan(&s.TotalRevenue)
	database.DB.QueryRow("SELECT COUNT(*) FROM products WHERE quantity <= 10").Scan(&s.LowStockCount)

	rows, err := database.DB.Query(`
		SELECT s.id, p.name, s.quantity, s.total, s.date
		FROM sales s
		JOIN products p ON s.product_id = p.id
		ORDER BY s.date DESC LIMIT 5
	`)
	if err != nil {
		return s, nil // return partial stats, not fatal
	}
	defer rows.Close()

	for rows.Next() {
		var sale models.Sale
		if err := rows.Scan(&sale.ID, &sale.ProductName, &sale.Quantity, &sale.Total, &sale.Date); err == nil {
			sale.FormattedTotal = models.FormatIndianRupees(sale.Total)
			s.RecentSales = append(s.RecentSales, sale)
		}
	}
	return s, nil
}
