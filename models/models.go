package models

import (
	"errors"
	"strings"
)

// Product represents an item in the inventory
type Product struct {
	ID       int     `json:"id"`
	Name     string  `json:"name"`
	Category string  `json:"category"`
	Price    float64 `json:"price"`
	Quantity int     `json:"quantity"`
	Status   string  `json:"status"` // Calculated field
}

// Sale represents a sales transaction
type Sale struct {
	ID        int     `json:"id"`
	ProductID int     `json:"product_id"`
	Quantity  int     `json:"quantity"`
	Total     float64 `json:"total"`
	Date      string  `json:"date"`

	// Join fields
	ProductName string `json:"product_name"`
}

// User represents an admin user
type User struct {
	ID           int    `json:"id"`
	Username     string `json:"username"`
	PasswordHash string `json:"-"`
	Email        string `json:"email"`
}

// CalculateStatus determines the stock status
func (p *Product) CalculateStatus() {
	if p.Quantity > 10 {
		p.Status = "In Stock"
	} else if p.Quantity > 0 {
		p.Status = "Low Stock"
	} else {
		p.Status = "Out of Stock"
	}
}

// Validate checks if product data is valid
func (p *Product) Validate() error {
	if strings.TrimSpace(p.Name) == "" {
		return errors.New("product name cannot be empty")
	}
	if p.Price <= 0 {
		return errors.New("price must be greater than 0")
	}
	if p.Quantity < 0 {
		return errors.New("quantity cannot be negative")
	}
	return nil
}
