package utils

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"strings"
	"time"

	_ "github.com/lib/pq"
)

// GetDataSource returns a Postgres data source string. It prefers DB_URL if set,
// GetDataSource returns a Postgres data source string composed from explicit
// environment variables: DB_HOST, DB_PORT, DB_USER, DB_PASSWORD, DB_NAME, DB_SSLMODE.
// We intentionally do NOT read DB_URL here (per project policy) so callers must
// configure the individual vars for Neon connections.
func GetDataSource() string {
	host := strings.Trim(os.Getenv("DB_HOST"), "\"' ")
	port := strings.Trim(os.Getenv("DB_PORT"), "\"' ")
	user := strings.Trim(os.Getenv("DB_USER"), "\"' ")
	pass := strings.Trim(os.Getenv("DB_PASSWORD"), "\"' ")
	name := strings.Trim(os.Getenv("DB_NAME"), "\"' ")
	ssl := strings.Trim(os.Getenv("DB_SSLMODE"), "\"' ")
	if ssl == "" {
		ssl = "require"
	}

	// Build pq-style connection string
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s", host, port, user, pass, name, ssl)
}

// TestDBConnection tries to open and ping the database using the provided data source name.
// It returns an error if the connection or ping fails.
func TestDBConnection(dataSource string, timeout time.Duration) error {
	if dataSource == "" {
		return fmt.Errorf("empty data source")
	}

	db, err := sql.Open("postgres", dataSource)
	if err != nil {
		return err
	}
	defer db.Close()

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	return db.PingContext(ctx)
}
