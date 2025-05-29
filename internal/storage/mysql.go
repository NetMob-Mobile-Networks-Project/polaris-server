package database

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"

	"github.com/NetMob-Mobile-Networks-Project/polaris-server/internal/storage/db"
)

type DB struct {
	*db.Queries
	conn *sql.DB
}

func NewMySQLConnection(dbUser, dbPass, dbHost, dbPort, dbName string) (*DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true",
		dbUser, dbPass, dbHost, dbPort, dbName)

	conn, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("error opening database: %v", err)
	}

	if err := conn.Ping(); err != nil {
		return nil, fmt.Errorf("error connecting to database: %v", err)
	}

	queries := db.New(conn)

	return &DB{
		Queries: queries,
		conn:    conn,
	}, nil
}

func (db *DB) Close() error {
	return db.conn.Close()
}
