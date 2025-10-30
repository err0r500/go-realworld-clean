package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/jackc/pgx/v5/stdlib"
)

// NewPostgresConnection creates a new PostgreSQL database connection pool using pgx
func NewPostgresConnection() (*sql.DB, error) {
	// Get connection string from environment with sensible defaults
	host := getEnvOrDefault("PG_HOST", "localhost")
	port := getEnvOrDefault("PG_PORT", "5432")
	user := getEnvOrDefault("PG_USER", "postgres")
	password := getEnvOrDefault("PG_PASSWORD", "postgres")
	dbname := getEnvOrDefault("PG_DB", "postgres")
	sslmode := getEnvOrDefault("PG_SSLMODE", "disable")

	// Construct connection string for pgx
	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		host, port, user, password, dbname, sslmode,
	)

	// Open connection using pgx driver
	db, err := sql.Open("pgx", connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to open database connection: %w", err)
	}

	// Verify connection
	if err := db.Ping(); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	// Configure connection pool
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)

	log.Printf("Successfully connected to PostgreSQL at %s:%s", host, port)
	return db, nil
}

// GetDatabaseURL returns the database URL for use with migrations
func GetDatabaseURL() string {
	host := getEnvOrDefault("PG_HOST", "localhost")
	port := getEnvOrDefault("PG_PORT", "5432")
	user := getEnvOrDefault("PG_USER", "postgres")
	password := getEnvOrDefault("PG_PASSWORD", "postgres")
	dbname := getEnvOrDefault("PG_DB", "postgres")
	sslmode := getEnvOrDefault("PG_SSLMODE", "disable")

	return fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		user, password, host, port, dbname, sslmode,
	)
}

func getEnvOrDefault(key, defaultValue string) string {
	value, ok := os.LookupEnv(key)
	if !ok {
		return defaultValue
	}
	return value
}
