package repository

import (
	"database/sql"
	"refactor/interfaces"

	_ "github.com/mattn/go-sqlite3"
)

type SQLiteOrderRepository struct {
	db *sql.DB
}

func NewSQLiteOrderRepository(db *sql.DB) *SQLiteOrderRepository {
	return &SQLiteOrderRepository{db: db}
}

func (r *SQLiteOrderRepository) Init() error {
	_, err := r.db.Exec(`
        CREATE TABLE IF NOT EXISTS orders (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            customer TEXT NOT NULL,
            products TEXT NOT NULL,
            total REAL NOT NULL,
            status TEXT NOT NULL
        )`)
	return err
}

func (r *SQLiteOrderRepository) Save(order *interfaces.Order) error {
	_, err := r.db.Exec(
		"INSERT INTO orders (customer, products, total, status) VALUES (?, ?, ?, ?)",
		order.Customer, order.Products, order.Total, order.Status,
	)
	return err
}
