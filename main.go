package main

import (
	"inventory-system/database"
	"inventory-system/internal/ui"
)

func main() {
	// Initialize SQLite database (unchanged schema)
	database.InitDB()
	defer database.DB.Close()

	// Launch Fyne native desktop application
	ui.Run()
}
