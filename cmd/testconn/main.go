package main

import (
	"fmt"
	"net/url"
	"os"
	"strings"
	"time"

	"badminton-reservation-api/utils"

	"github.com/joho/godotenv"
)

func maskDSN(ds string) string {
	ds = strings.Trim(ds, "\"' ")
	// Expect pq-style: host=... port=... user=... password=... dbname=... sslmode=...
	parts := strings.Fields(ds)
	for _, p := range parts {
		if strings.HasPrefix(p, "host=") {
			return strings.TrimPrefix(p, "host=")
		}
	}
	// Fallback: try to parse as URL
	if strings.HasPrefix(ds, "postgres://") || strings.HasPrefix(ds, "postgresql://") {
		u, err := url.Parse(ds)
		if err != nil {
			return "(unable to parse DSN)"
		}
		return u.Hostname()
	}
	return "(unknown)"
}

func main() {
	_ = godotenv.Load()

	ds := utils.GetDataSource()
	source := "env vars (DB_HOST etc)"

	fmtHost := maskDSN(ds)
	fmt.Printf("Testing DB connection (source: %s, host: %s)...\n", source, fmtHost)

	// Use a longer timeout for remote Neon connections
	err := utils.TestDBConnection(ds, 15*time.Second)
	if err != nil {
		fmt.Printf("Connection FAILED: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("Connection OK")
}
