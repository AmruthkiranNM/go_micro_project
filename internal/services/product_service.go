package services

import (
	"fmt"
	"inventory-system/database"
	"inventory-system/models"
	"strings"
)

// ListProducts returns all products from the database.
func ListProducts() ([]models.Product, error) {
	rows, err := database.DB.Query("SELECT id, name, category, price, quantity FROM products ORDER BY id")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []models.Product
	for rows.Next() {
		var p models.Product
		if err := rows.Scan(&p.ID, &p.Name, &p.Category, &p.Price, &p.Quantity); err != nil {
			continue
		}
		p.CalculateStatus()
		p.FormattedPrice = models.FormatIndianRupees(p.Price)
		products = append(products, p)
	}
	return products, rows.Err()
}

// ListProductsInStock returns only products that have stock > 0 (for sales dropdown).
func ListProductsInStock() ([]models.Product, error) {
	rows, err := database.DB.Query("SELECT id, name, quantity FROM products WHERE quantity > 0 ORDER BY name")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []models.Product
	for rows.Next() {
		var p models.Product
		if err := rows.Scan(&p.ID, &p.Name, &p.Quantity); err != nil {
			continue
		}
		products = append(products, p)
	}
	return products, rows.Err()
}

// AddProduct inserts a new product after validation.
func AddProduct(p models.Product) error {
	if err := p.Validate(); err != nil {
		return err
	}
	_, err := database.DB.Exec(
		"INSERT INTO products (name, category, price, quantity) VALUES (?, ?, ?, ?)",
		p.Name, p.Category, p.Price, p.Quantity,
	)
	return err
}

// UpdateProduct updates an existing product identified by id.
func UpdateProduct(id int, p models.Product) error {
	if err := p.Validate(); err != nil {
		return err
	}
	_, err := database.DB.Exec(
		"UPDATE products SET name=?, category=?, price=?, quantity=? WHERE id=?",
		p.Name, p.Category, p.Price, p.Quantity, id,
	)
	return err
}

// DeleteProduct removes a product; returns error if it has sales records.
func DeleteProduct(id int) error {
	_, err := database.DB.Exec("DELETE FROM products WHERE id = ?", id)
	if err != nil {
		if strings.Contains(err.Error(), "FOREIGN KEY constraint failed") {
			return fmt.Errorf("product cannot be deleted because it has recorded sales history. Try setting its quantity to 0 instead")
		}
		return err
	}
	return nil
}
