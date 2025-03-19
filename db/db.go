package db

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "github.com/go-sql-driver/mysql"
)

type Database struct {
	db *sqlx.DB
}

func NewDatabase(url string) (*Database, error) {
	db, err := sqlx.Open("mysql", url)
	if err != nil {
		return nil, fmt.Errorf("error opening database: %w", err)
	}

	return &Database{db: db}, nil
}

func (d *Database) Close() error {
	return d.db.Close()
}

func (d *Database) GetDB() *sqlx.DB {
	return d.db
}