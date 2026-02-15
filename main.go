package main

import (
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-contrib/multitemplate"
	"github.com/gin-gonic/gin"
)

func main() {
	// Initialize Database
	InitDB()
	defer DB.Close()

	router := gin.Default()

	// Use multitemplate renderer
	router.HTMLRender = createMyRender()

	// Serve static files
	router.Static("/static", "./static")

	// Routes
	router.GET("/", func(c *gin.Context) {
		// Fetch counts for dashboard
		var productCount, salesCount int
		var totalRevenue float64

		DB.QueryRow("SELECT COUNT(*) FROM products").Scan(&productCount)
		DB.QueryRow("SELECT COUNT(*) FROM sales").Scan(&salesCount)
		DB.QueryRow("SELECT COALESCE(SUM(total), 0) FROM sales").Scan(&totalRevenue)

		c.HTML(http.StatusOK, "dashboard.html", gin.H{
			"title":        "Dashboard",
			"productCount": productCount,
			"salesCount":   salesCount,
			"totalRevenue": totalRevenue,
		})
	})

	router.GET("/products", func(c *gin.Context) {
		rows, err := DB.Query("SELECT id, name, category, price, quantity FROM products")
		if err != nil {
			c.String(http.StatusInternalServerError, "Database error")
			return
		}
		defer rows.Close()

		var products []Product
		for rows.Next() {
			var p Product
			if err := rows.Scan(&p.ID, &p.Name, &p.Category, &p.Price, &p.Quantity); err != nil {
				continue
			}
			// Determine status
			if p.Quantity > 10 {
				p.Status = "In Stock"
			} else if p.Quantity > 0 {
				p.Status = "Low Stock"
			} else {
				p.Status = "Out of Stock"
			}
			products = append(products, p)
		}

		c.HTML(http.StatusOK, "products.html", gin.H{
			"title":    "Products",
			"products": products,
		})
	})

	router.POST("/products/add", func(c *gin.Context) {
		name := c.PostForm("name")
		category := c.PostForm("category")
		price, _ := strconv.ParseFloat(c.PostForm("price"), 64)
		quantity, _ := strconv.Atoi(c.PostForm("quantity"))

		_, err := DB.Exec("INSERT INTO products (name, category, price, quantity) VALUES (?, ?, ?, ?)",
			name, category, price, quantity)
		if err != nil {
			log.Println("Error inserting product:", err)
		}
		c.Redirect(http.StatusFound, "/products")
	})

	router.POST("/products/delete/:id", func(c *gin.Context) {
		id := c.Param("id")
		DB.Exec("DELETE FROM products WHERE id = ?", id)
		c.Redirect(http.StatusFound, "/products")
	})

	// Basic Update implementation (Updating name/price/qty)
	router.POST("/products/update/:id", func(c *gin.Context) {
		id := c.Param("id")
		name := c.PostForm("name")
		category := c.PostForm("category")
		price, _ := strconv.ParseFloat(c.PostForm("price"), 64)
		quantity, _ := strconv.Atoi(c.PostForm("quantity"))

		_, err := DB.Exec("UPDATE products SET name=?, category=?, price=?, quantity=? WHERE id=?",
			name, category, price, quantity, id)
		if err != nil {
			log.Println("Error updating product:", err)
		}
		c.Redirect(http.StatusFound, "/products")
	})

	router.GET("/sales", func(c *gin.Context) {
		rows, err := DB.Query(`
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

		var sales []Sale
		for rows.Next() {
			var s Sale
			if err := rows.Scan(&s.ID, &s.ProductName, &s.Quantity, &s.Total, &s.Date); err != nil {
				continue
			}
			sales = append(sales, s)
		}

		// Fetch products for the dropdown in the 'Add Sale' modal
		prodRows, _ := DB.Query("SELECT id, name FROM products WHERE quantity > 0")
		defer prodRows.Close()
		var products []Product
		for prodRows.Next() {
			var p Product
			prodRows.Scan(&p.ID, &p.Name)
			products = append(products, p)
		}

		c.HTML(http.StatusOK, "sales.html", gin.H{
			"title":    "Sales",
			"sales":    sales,
			"products": products,
		})
	})

	router.POST("/sales/add", func(c *gin.Context) {
		productID, _ := strconv.Atoi(c.PostForm("product_id"))
		quantity, _ := strconv.Atoi(c.PostForm("quantity"))

		// Get product price and current quantity
		var price float64
		var currentQty int
		err := DB.QueryRow("SELECT price, quantity FROM products WHERE id = ?", productID).Scan(&price, &currentQty)
		if err != nil {
			log.Println("Error finding product:", err)
			c.Redirect(http.StatusFound, "/sales")
			return
		}

		if currentQty < quantity {
			// Not enough stock
			// In a real app we would show an error message
			log.Println("Not enough stock")
			c.Redirect(http.StatusFound, "/sales")
			return
		}

		total := price * float64(quantity)
		date := time.Now().Format("2006-01-02 15:04:05")

		// Transaction
		tx, err := DB.Begin()
		if err != nil {
			return
		}

		// Insert Sale
		_, err = tx.Exec("INSERT INTO sales (product_id, quantity, total, date) VALUES (?, ?, ?, ?)",
			productID, quantity, total, date)
		if err != nil {
			tx.Rollback()
			c.Redirect(http.StatusFound, "/sales")
			return
		}

		// Update Inventory
		_, err = tx.Exec("UPDATE products SET quantity = quantity - ? WHERE id = ?", quantity, productID)
		if err != nil {
			tx.Rollback()
			c.Redirect(http.StatusFound, "/sales")
			return
		}

		tx.Commit()
		c.Redirect(http.StatusFound, "/sales")
	})

	router.GET("/reports", func(c *gin.Context) {
		c.HTML(http.StatusOK, "reports.html", gin.H{
			"title": "Reports",
		})
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s", port)
	router.Run(":" + port)
}

func createMyRender() multitemplate.Renderer {
	r := multitemplate.NewRenderer()
	r.AddFromFiles("dashboard.html", "templates/layout.html", "templates/dashboard.html")
	r.AddFromFiles("products.html", "templates/layout.html", "templates/products.html")
	r.AddFromFiles("sales.html", "templates/layout.html", "templates/sales.html")
	r.AddFromFiles("reports.html", "templates/layout.html", "templates/reports.html")
	return r
}
