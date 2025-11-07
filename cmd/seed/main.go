package main

import (
	"bufio"
	"database/sql"
	"fmt"
	"io"
	"os"
	"strings"

	"badminton-reservation-api/utils"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	_ = godotenv.Load()

	ds := utils.GetDataSource()
	if ds == "" {
		fmt.Println("No data source available. Set DB_URL or DB_* env vars.")
		os.Exit(1)
	}

	db, err := sql.Open("postgres", ds)
	if err != nil {
		fmt.Println("Failed to open DB:", err)
		os.Exit(1)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		fmt.Println("DB ping failed:", err)
		os.Exit(1)
	}

	f, err := os.Open("database/seeds/seed_data.sql")
	if err != nil {
		fmt.Println("Failed to open seed file:", err)
		os.Exit(1)
	}
	defer f.Close()

	reader := bufio.NewReader(f)
	var sb strings.Builder
	for {
		line, err := reader.ReadString('\n')
		if err != nil && err != io.EOF {
			fmt.Println("Error reading seed file:", err)
			os.Exit(1)
		}
		sb.WriteString(line)
		if err == io.EOF {
			break
		}
	}

	// Split statements by semicolon; naive but OK for simple seeds
	sqlText := sb.String()
	stmts := strings.Split(sqlText, ";")
	for _, s := range stmts {
		stmt := strings.TrimSpace(s)
		if stmt == "" {
			continue
		}
		if _, err := db.Exec(stmt); err != nil {
			fmt.Printf("Failed to execute statement: %v\nError: %v\n", stmt, err)
			os.Exit(1)
		}
	}

	fmt.Println("Seed applied successfully")
}
