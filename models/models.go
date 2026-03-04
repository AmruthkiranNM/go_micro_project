package models

import (
	"errors"
	"fmt"
	"strings"
)

// Product represents an item in the inventory
type Product struct {
	ID             int     `json:"id"`
	Name           string  `json:"name"`
	Category       string  `json:"category"`
	Price          float64 `json:"price"`
	FormattedPrice string  `json:"formatted_price"`
	Quantity       int     `json:"quantity"`
	Status         string  `json:"status"` // Calculated field
}

// Sale represents a sales transaction
type Sale struct {
	ID             int     `json:"id"`
	ProductID      int     `json:"product_id"`
	Quantity       int     `json:"quantity"`
	Total          float64 `json:"total"`
	FormattedTotal string  `json:"formatted_total"`
	Date           string  `json:"date"`

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

// FormatIndianRupees formats a float64 into the Indian numbering system format with a ₹ symbol.
func FormatIndianRupees(amount float64) string {
	// Handle integer parts and fractional parts separately
	// We'll format it manually to handle the 3,2,2 grouping of the Indian system
	amountStr := fmt.Sprintf("%.2f", amount)
	parts := strings.Split(amountStr, ".")

	integerPart := parts[0]
	decimalPart := parts[1]

	if len(integerPart) <= 3 {
		return fmt.Sprintf("₹%s.%s", integerPart, decimalPart)
	}

	// Get the last 3 digits
	lastThree := integerPart[len(integerPart)-3:]
	// Get the remaining digits
	remainingParams := integerPart[:len(integerPart)-3]

	var resultStrings []string

	// Group the remaining digits by 2
	for i := len(remainingParams); i > 0; i -= 2 {
		start := i - 2
		if start < 0 {
			start = 0
		}
		resultStrings = append([]string{remainingParams[start:i]}, resultStrings...)
	}

	resultStrings = append(resultStrings, lastThree)
	formattedInteger := strings.Join(resultStrings, ",")

	return fmt.Sprintf("₹%s.%s", formattedInteger, decimalPart)
}
