package main

import (
	"log"
	"os"

	"inventory-system/database"
	"inventory-system/handlers"
	"inventory-system/middleware"

	"github.com/gin-contrib/multitemplate"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

func main() {
	// Initialize Database
	database.InitDB()
	defer database.DB.Close()

	router := gin.Default()

	// Sessions
	store := cookie.NewStore([]byte("secret"))
	router.Use(sessions.Sessions("mysession", store))

	// Use multitemplate renderer
	router.HTMLRender = createMyRender()

	// Serve static files
	router.Static("/static", "./static")

	// Public Routes
	router.GET("/login", handlers.ShowLogin)
	router.POST("/login", handlers.Login)
	router.GET("/logout", handlers.Logout)

	// Protected Routes
	protected := router.Group("/")
	protected.Use(middleware.AuthRequired())
	{
		protected.GET("/", handlers.Dashboard)

		// Products
		protected.GET("/products", handlers.ListProducts)
		protected.POST("/products/add", handlers.AddProduct)
		protected.POST("/products/update/:id", handlers.UpdateProduct)
		protected.POST("/products/delete/:id", handlers.DeleteProduct)

		// Sales
		protected.GET("/sales", handlers.ListSales)
		protected.POST("/sales/add", handlers.RecordSale)

		// Reports
		protected.GET("/reports", handlers.ShowReports)

		// Admin Profile
		protected.GET("/profile", handlers.Profile)
		protected.POST("/profile", handlers.UpdateProfile)
	}

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
	r.AddFromFiles("profile.html", "templates/layout.html", "templates/profile.html")
	// Login does not use layout
	r.AddFromFiles("login.html", "templates/login.html")
	return r
}
