package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

type Product struct {
	ID       string
	Name     string
	Category string
	Price    float64
	Quantity int
	Status   string
}

type Sale struct {
	InvoiceID string
	Product   string
	Quantity  int
	Total     float64
	Date      string
}

func main() {
	router := gin.Default()

	// Serve static files
	router.Static("/static", "./static")

	// Load templates
	router.LoadHTMLGlob("templates/*")

	// Routes
	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "dashboard.html", gin.H{
			"title": "Dashboard",
		})
	})

	router.GET("/products", func(c *gin.Context) {
		products := []Product{
			{ID: "1", Name: "Laptop", Category: "Electronics", Price: 1200.00, Quantity: 10, Status: "In Stock"},
			{ID: "2", Name: "Mouse", Category: "Accessories", Price: 25.00, Quantity: 50, Status: "In Stock"},
			{ID: "3", Name: "Keyboard", Category: "Accessories", Price: 45.00, Quantity: 30, Status: "In Stock"},
			{ID: "4", Name: "Monitor", Category: "Electronics", Price: 300.00, Quantity: 5, Status: "Low Stock"},
		}
		c.HTML(http.StatusOK, "products.html", gin.H{
			"title":    "Products",
			"products": products,
		})
	})

	router.GET("/sales", func(c *gin.Context) {
		sales := []Sale{
			{InvoiceID: "INV-001", Product: "Laptop", Quantity: 1, Total: 1200.00, Date: "2023-10-25"},
			{InvoiceID: "INV-002", Product: "Mouse", Quantity: 2, Total: 50.00, Date: "2023-10-26"},
			{InvoiceID: "INV-003", Product: "Monitor", Quantity: 1, Total: 300.00, Date: "2023-10-26"},
		}
		c.HTML(http.StatusOK, "sales.html", gin.H{
			"title": "Sales",
			"sales": sales,
		})
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
