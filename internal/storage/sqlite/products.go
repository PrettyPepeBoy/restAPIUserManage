package sqlite

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/mattn/go-sqlite3"
	"tstUser/internal/storage"
)

type Products struct {
	db         *sql.DB
	StmtCreate *sql.Stmt
	StmtUpdate *sql.Stmt
}

func NewProductsTable(storagePath string) (*Products, error) {
	const op = "internal/storage/sqlite/NewProductsTable"

	db, err := sql.Open("sqlite3", storagePath)

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	stmt, err := db.Prepare(`
	CREATE TABLE IF NOT EXISTS products(
	    id INTEGER PRIMARY KEY,
	    name TEXT NOT NULL UNIQUE,
	    price INT NOT NULL,
	    amount INT);
	CREATE INDEX IF NOT EXISTS idx_name  ON products(name)
	`)

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	_, err = stmt.Exec()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	stmtCreate, err := db.Prepare(`INSERT INTO products(name, price, amount) VALUES (?,?,?)`)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	stmtUpdate, err := db.Prepare(`UPDATE products SET amount = ? WHERE name = ?`)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return &Products{db: db,
		StmtCreate: stmtCreate,
		StmtUpdate: stmtUpdate,
	}, nil
}

func (p *Products) CreateProducts(name string, price, amount int64) (int64, error) {
	const op = "internal/storage/sqlite/CreateProducts"
	res, err := p.StmtCreate.Exec(name, price, amount)
	if err != nil {
		var sqliteErr sqlite3.Error
		if errors.As(err, &sqliteErr) && errors.Is(sqliteErr.ExtendedCode, sqlite3.ErrConstraintUnique) {
			return 0, fmt.Errorf("%s, %w", op, storage.ErrProductsExist)
		}
		return 0, fmt.Errorf("%s: execute statemnt: %w", op, err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("%s: failed to get last insert id: %w", op, err)
	}
	return id, nil

}

func (p *Products) UpdateProducts(name string, amount int64) (int64, error) {
	const op = "internal/storage/sqlite/UpdateProducts"
	_, err := p.StmtUpdate.Exec(name, amount)
	if err != nil {
		return 0, fmt.Errorf("%s: execute statement: %w", op, err)
	}
	return 0, nil
}
