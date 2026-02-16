package database

import (
	"database/sql"
	"log"

	"golang.org/x/crypto/bcrypt"
	_ "modernc.org/sqlite"
)

var DB *sql.DB

func InitDB() {
	var err error
	DB, err = sql.Open("sqlite", "./inventory.db")
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Enable Foreign Keys
	if _, err := DB.Exec("PRAGMA foreign_keys = ON;"); err != nil {
		log.Fatal("Failed to enable foreign keys:", err)
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
		price REAL NOT NULL CHECK(price > 0),
		quantity INTEGER NOT NULL CHECK(quantity >= 0)
	);`

	if _, err := DB.Exec(productTable); err != nil {
		log.Fatal("Failed to create products table:", err)
	}

	// Create sales table
	salesTable := `
	CREATE TABLE IF NOT EXISTS sales (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		product_id INTEGER,
		quantity INTEGER CHECK(quantity > 0),
		total REAL,
		date TEXT,
		FOREIGN KEY(product_id) REFERENCES products(id) ON DELETE RESTRICT
	);`

	if _, err := DB.Exec(salesTable); err != nil {
		log.Fatal("Failed to create sales table:", err)
	}

	// Create users table
	usersTable := `
	CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		username TEXT UNIQUE NOT NULL,
		password_hash TEXT NOT NULL,
		email TEXT
	);`

	if _, err := DB.Exec(usersTable); err != nil {
		log.Fatal("Failed to create users table:", err)
	}

	log.Println("Database initialized successfully")
	SeedAdmin()
}

func SeedAdmin() {
	var count int
	err := DB.QueryRow("SELECT COUNT(*) FROM users").Scan(&count)
	if err != nil {
		log.Println("Failed to count users:", err)
		return
	}

	if count == 0 {
		log.Println("Creating default admin user...")
		password := "admin123"
		hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			log.Fatal("Failed to hash password:", err)
		}

		_, err = DB.Exec("INSERT INTO users (username, password_hash, email) VALUES (?, ?, ?)",
			"admin", string(hash), "admin@example.com")
		if err != nil {
			log.Fatal("Failed to seed admin user:", err)
		}
		log.Println("Default admin user created: admin / admin123")
	}
}
