package handlers

import (
	"net/http"
	"strconv"

	"inventory-system/database"
	"inventory-system/models"

	"github.com/gin-gonic/gin"
)

func ListProducts(c *gin.Context) {
	rows, err := database.DB.Query("SELECT id, name, category, price, quantity FROM products")
	if err != nil {
		c.String(http.StatusInternalServerError, "Database error")
		return
	}
	defer rows.Close()

	var products []models.Product
	for rows.Next() {
		var p models.Product
		if err := rows.Scan(&p.ID, &p.Name, &p.Category, &p.Price, &p.Quantity); err != nil {
			continue
		}
		p.CalculateStatus()
		products = append(products, p)
	}

	c.HTML(http.StatusOK, "products.html", gin.H{
		"title":    "Products",
		"products": products,
	})
}

func AddProduct(c *gin.Context) {
	var p models.Product
	p.Name = c.PostForm("name")
	p.Category = c.PostForm("category")
	p.Price, _ = strconv.ParseFloat(c.PostForm("price"), 64)
	p.Quantity, _ = strconv.Atoi(c.PostForm("quantity"))

	if err := p.Validate(); err != nil {
		// In a real app, we'd pass this back to the template
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	_, err := database.DB.Exec("INSERT INTO products (name, category, price, quantity) VALUES (?, ?, ?, ?)",
		p.Name, p.Category, p.Price, p.Quantity)
	if err != nil {
		c.String(http.StatusInternalServerError, "Failed to add product")
		return
	}
	c.Redirect(http.StatusFound, "/products")
}

func UpdateProduct(c *gin.Context) {
	id := c.Param("id")
	var p models.Product
	p.Name = c.PostForm("name")
	p.Category = c.PostForm("category")
	p.Price, _ = strconv.ParseFloat(c.PostForm("price"), 64)
	p.Quantity, _ = strconv.Atoi(c.PostForm("quantity"))

	if err := p.Validate(); err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	_, err := database.DB.Exec("UPDATE products SET name=?, category=?, price=?, quantity=? WHERE id=?",
		p.Name, p.Category, p.Price, p.Quantity, id)
	if err != nil {
		c.String(http.StatusInternalServerError, "Failed to update product")
		return
	}
	c.Redirect(http.StatusFound, "/products")
}

func DeleteProduct(c *gin.Context) {
	id := c.Param("id")

	// Check if product is in sales (Foreign Key constraint handles this if ON DELETE RESTRICT)
	_, err := database.DB.Exec("DELETE FROM products WHERE id = ?", id)
	if err != nil {
		c.String(http.StatusBadRequest, "Cannot delete product (it might have sales records)")
		return
	}
	c.Redirect(http.StatusFound, "/products")
}
