package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"badminton-reservation-api/utils"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

// WARNING: This program will DROP ALL TABLES in the public schema of the configured database.
// Use with caution.
func main() {
	_ = godotenv.Load()

	ds := utils.GetDataSource()
	if ds == "" {
		log.Fatal("no data source available; set DB_URL or DB_HOST/DB_* env vars")
	}

	db, err := sql.Open("postgres", ds)
	if err != nil {
		log.Fatalf("failed to open db: %v", err)
	}
	defer db.Close()

	// Ping to ensure connection
	if err := db.Ping(); err != nil {
		log.Fatalf("ping failed: %v", err)
	}

	rows, err := db.Query("SELECT table_name FROM information_schema.tables WHERE table_schema='public' AND table_type='BASE TABLE'")
	if err != nil {
		log.Fatalf("failed to list tables: %v", err)
	}
	defer rows.Close()

	var tables []string
	for rows.Next() {
		var t string
		if err := rows.Scan(&t); err != nil {
			log.Fatalf("scan error: %v", err)
		}
		tables = append(tables, t)
	}

	if len(tables) == 0 {
		fmt.Println("No tables found in public schema.")
		os.Exit(0)
	}

	fmt.Println("Tables to drop:")
	for _, t := range tables {
		fmt.Println(" -", t)
	}

	// Double-check environment to avoid accidental runs in dev: require confirmation via env var DROP_DB_CONFIRM=yes
	if os.Getenv("DROP_DB_CONFIRM") != "yes" {
		fmt.Println("To actually drop the above tables, re-run with environment variable DROP_DB_CONFIRM=yes")
		os.Exit(1)
	}

	for _, t := range tables {
		q := fmt.Sprintf("DROP TABLE IF EXISTS \"%s\" CASCADE", t)
		if _, err := db.Exec(q); err != nil {
			log.Fatalf("failed to drop table %s: %v", t, err)
		}
		fmt.Printf("Dropped %s\n", t)
	}

	fmt.Println("All tables dropped.")
}
