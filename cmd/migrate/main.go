package main

import (
	"fmt"
	"os"

	"badminton-reservation-api/database"

	"github.com/joho/godotenv"
)

func main() {
	// load .env if present
	_ = godotenv.Load()

	fmt.Println("Running GORM migrations...")
	if err := database.RunGormMigrations(); err != nil {
		fmt.Println("Migration failed:", err)
		os.Exit(1)
	}
	fmt.Println("Migrations applied successfully")
}
