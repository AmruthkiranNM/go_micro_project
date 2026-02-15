package main

// Product represents an item in the inventory
type Product struct {
	ID       int     `json:"id"`
	Name     string  `json:"name"`
	Category string  `json:"category"`
	Price    float64 `json:"price"`
	Quantity int     `json:"quantity"`
	Status   string  `json:"status"` // Calculated field, not in DB
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
