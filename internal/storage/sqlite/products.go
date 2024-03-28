package sqlite

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/mattn/go-sqlite3"
	"tstUser/internal/http-server/transport/productDTO"
	"tstUser/internal/storage"
)

type Products struct {
	db         *sql.DB
	StmtCreate *sql.Stmt
	StmtUpdate *sql.Stmt
	StmtGet    *sql.Stmt
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

	stmtUpdate, err := db.Prepare(`UPDATE products SET name = ?, amount = ?, price =? WHERE id = ?`)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	stmtGet, err := db.Prepare(`SELECT name, price, amount FROM products WHERE id = ?`)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return &Products{db: db,
		StmtCreate: stmtCreate,
		StmtUpdate: stmtUpdate,
		StmtGet:    stmtGet,
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

func (p *Products) UpdateProducts(up productDTO.ProductDTO) error {
	const op = "internal/storage/sqlite/UpdateProducts"
	res, err := p.StmtUpdate.Exec(up.Name, up.Amount, up.Price, up.ID)
	if err != nil {
		return fmt.Errorf("%s: execute statement: %w", op, err)
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("%s: execute statement: %w", op, err)
	}
	if rows == 0 {
		return fmt.Errorf("%s, %w", op, storage.ErrProductNotFound)
	}
	return nil
}

func (p *Products) GetProducts(ID int64) (productDTO.ProductDTO, error) {
	const op = "internal/storage/sqlite/GetProducts"
	var product productDTO.ProductDTO
	err := p.StmtGet.QueryRow(ID).Scan(&product.Name, &product.Price, &product.Amount)
	if errors.Is(err, sql.ErrNoRows) {
		return productDTO.ProductDTO{}, fmt.Errorf("%s, %w", op, storage.ErrProductNotFound)
	}
	if err != nil {
		return productDTO.ProductDTO{}, fmt.Errorf("%s: execute statement %w", op, err)
	}
	product.ID = ID
	return product, nil
}
