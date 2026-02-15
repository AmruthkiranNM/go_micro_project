package main

import (
	"database/sql"
	"log"

	_ "modernc.org/sqlite"
)

var DB *sql.DB

func InitDB() {
	var err error
	DB, err = sql.Open("sqlite", "./inventory.db")
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	createTables()
}

func createTables() {
	// Create products table
	productTable := `
	CREATE TABLE IF NOT EXISTS products (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		category TEXT NOT NULL,
		price REAL NOT NULL,
		quantity INTEGER NOT NULL
	);`

	_, err := DB.Exec(productTable)
	if err != nil {
		log.Fatal("Failed to create products table:", err)
	}

	// Create sales table
	salesTable := `
	CREATE TABLE IF NOT EXISTS sales (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		product_id INTEGER,
		quantity INTEGER,
		total REAL,
		date TEXT,
		FOREIGN KEY(product_id) REFERENCES products(id)
	);`

	_, err = DB.Exec(salesTable)
	if err != nil {
		log.Fatal("Failed to create sales table:", err)
	}

	log.Println("Database initialized successfully")
}
