package handlers

import (
	"database/sql"
	"net/http"
	"strconv"
	"time"

	"inventory-system/database"
	"inventory-system/models"

	"github.com/gin-gonic/gin"
)

func ListSales(c *gin.Context) {
	rows, err := database.DB.Query(`
		SELECT s.id, p.name, s.quantity, s.total, s.date 
		FROM sales s 
		JOIN products p ON s.product_id = p.id
		ORDER BY s.date DESC
	`)
	if err != nil {
		c.String(http.StatusInternalServerError, "Database error: "+err.Error())
		return
	}
	defer rows.Close()

	var sales []models.Sale
	for rows.Next() {
		var s models.Sale
		if err := rows.Scan(&s.ID, &s.ProductName, &s.Quantity, &s.Total, &s.Date); err != nil {
			continue
		}
		sales = append(sales, s)
	}

	// Fetch products for the dropdown
	prodRows, _ := database.DB.Query("SELECT id, name FROM products WHERE quantity > 0")
	defer prodRows.Close()
	var products []models.Product
	for prodRows.Next() {
		var p models.Product
		prodRows.Scan(&p.ID, &p.Name)
		products = append(products, p)
	}

	c.HTML(http.StatusOK, "sales.html", gin.H{
		"title":    "Sales",
		"sales":    sales,
		"products": products,
	})
}

func RecordSale(c *gin.Context) {
	productID, _ := strconv.Atoi(c.PostForm("product_id"))
	quantity, _ := strconv.Atoi(c.PostForm("quantity"))

	if quantity <= 0 {
		c.String(http.StatusBadRequest, "Quantity must be greater than 0")
		return
	}

	// Start Transaction
	tx, err := database.DB.Begin()
	if err != nil {
		c.String(http.StatusInternalServerError, "Failed to start transaction")
		return
	}
	defer tx.Rollback()

	// 1. Get product price and current quantity (with lock)
	var price float64
	var currentQty int
	err = tx.QueryRow("SELECT price, quantity FROM products WHERE id = ?", productID).Scan(&price, &currentQty)
	if err != nil {
		if err == sql.ErrNoRows {
			c.String(http.StatusNotFound, "Product not found")
		} else {
			c.String(http.StatusInternalServerError, "Database error")
		}
		return
	}

	// 2. Check Stock
	if currentQty < quantity {
		c.String(http.StatusBadRequest, "Insufficient stock")
		return
	}

	total := price * float64(quantity)
	date := time.Now().Format("2006-01-02 15:04:05")

	// 3. Insert Sale
	_, err = tx.Exec("INSERT INTO sales (product_id, quantity, total, date) VALUES (?, ?, ?, ?)",
		productID, quantity, total, date)
	if err != nil {
		c.String(http.StatusInternalServerError, "Failed to record sale")
		return
	}

	// 4. Update Inventory
	_, err = tx.Exec("UPDATE products SET quantity = quantity - ? WHERE id = ?", quantity, productID)
	if err != nil {
		c.String(http.StatusInternalServerError, "Failed to update stock")
		return
	}

	// 5. Commit
	if err := tx.Commit(); err != nil {
		c.String(http.StatusInternalServerError, "Failed to commit transaction")
		return
	}

	c.Redirect(http.StatusFound, "/sales")
}

func Dashboard(c *gin.Context) {
	var productCount, salesCount int
	var totalRevenue float64

	database.DB.QueryRow("SELECT COUNT(*) FROM products").Scan(&productCount)
	database.DB.QueryRow("SELECT COUNT(*) FROM sales").Scan(&salesCount)
	database.DB.QueryRow("SELECT COALESCE(SUM(total), 0) FROM sales").Scan(&totalRevenue)

	c.HTML(http.StatusOK, "dashboard.html", gin.H{
		"title":        "Dashboard",
		"productCount": productCount,
		"salesCount":   salesCount,
		"totalRevenue": totalRevenue,
	})
}
