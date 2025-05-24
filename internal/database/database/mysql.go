package database

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/themhh/polaris-server/internal/database/db"
)

// DB represents the database connection and queries
type DB struct {
	*db.Queries
	db *sql.DB
}

// NewMySQLConnection creates a new MySQL database connection with sqlc queries
func NewMySQLConnection() (*DB, error) {
	// Get database configuration from environment variables
	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASS")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")

	// Construct DSN
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true",
		dbUser, dbPass, dbHost, dbPort, dbName)

	// Open database connection
	sqlDB, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("error opening database: %v", err)
	}

	// Test the connection
	if err := sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("error connecting to database: %v", err)
	}

	// Create queries instance
	queries := db.New(sqlDB)

	return &DB{
		Queries: queries,
		db:      sqlDB,
	}, nil
}

// Close closes the database connection
func (db *DB) Close() error {
	return db.db.Close()
}
