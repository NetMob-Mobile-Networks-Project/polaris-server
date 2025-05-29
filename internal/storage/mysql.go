package database

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"

	"github.com/NetMob-Mobile-Networks-Project/polaris-server/internal/storage/db"
)

// DB represents the database connection and queries
type DB struct {
	*db.Queries
	conn *sql.DB
}

// NewMySQLConnection creates a new MySQL database connection with sqlc queries
func NewMySQLConnection(dbUser, dbPass, dbHost, dbPort, dbName string) (*DB, error) {
	// Construct DSN
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true",
		dbUser, dbPass, dbHost, dbPort, dbName)

	// Open database connection
	conn, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("error opening database: %v", err)
	}

	// Test the connection
	if err := conn.Ping(); err != nil {
		return nil, fmt.Errorf("error connecting to database: %v", err)
	}

	// Create queries instance
	queries := db.New(conn)

	return &DB{
		Queries: queries,
		conn:    conn,
	}, nil
}

// Close closes the database connection
func (db *DB) Close() error {
	return db.conn.Close()
}
