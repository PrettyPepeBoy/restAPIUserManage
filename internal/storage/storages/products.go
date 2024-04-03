package storages

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/mattn/go-sqlite3"
	"tstUser/internal/storage/storages/errs"
)

type ProductStorage struct {
	db         *sql.DB
	StmtCreate *sql.Stmt
	StmtUpdate *sql.Stmt
	StmtGet    *sql.Stmt
	StmtDelete *sql.Stmt
}

type Product struct {
	Id     int64
	Name   string
	Price  int
	Amount int
}

func NewProductsTable(storagePath string) (*ProductStorage, error) {
	const op = "internal/storage/storages/NewProductsTable"

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

	stmtUpdate, err := db.Prepare(`UPDATE products SET name = ?, price = ?, amount =? WHERE id = ?`)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	stmtGet, err := db.Prepare(`SELECT name, price, amount FROM products WHERE id = ?`)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	stmtDelete, err := db.Prepare(`DELETE FROM products WHERE id = ?`)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &ProductStorage{db: db,
		StmtCreate: stmtCreate,
		StmtUpdate: stmtUpdate,
		StmtGet:    stmtGet,
		StmtDelete: stmtDelete,
	}, nil
}

func (p *ProductStorage) CreateProducts(product Product) (int64, error) {
	const op = "internal/storage/storages/CreateProducts"
	res, err := p.StmtCreate.Exec(product.Name, product.Price, product.Amount)
	if err != nil {
		var sqliteErr sqlite3.Error
		if errors.As(err, &sqliteErr) && errors.Is(sqliteErr.ExtendedCode, sqlite3.ErrConstraintUnique) {
			return 0, fmt.Errorf("%s, %w", op, errs.ErrProductsExist)
		}
		return 0, fmt.Errorf("%s: execute statemnt: %w", op, err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("%s: failed to get last insert id: %w", op, err)
	}
	return id, nil

}

func (p *ProductStorage) UpdateProducts(product Product) error {
	const op = "internal/storage/storages/UpdateProducts"
	res, err := p.StmtUpdate.Exec(product.Name, product.Price, product.Amount, product.Id)
	if err != nil {
		return fmt.Errorf("%s: execute statement: %w", op, err)
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("%s: execute statement: %w", op, err)
	}
	if rows == 0 {
		return fmt.Errorf("%s, %w", op, errs.ErrProductNotFound)
	}
	return nil
}

func (p *ProductStorage) GetProducts(id int64) (Product, error) {
	const op = "internal/storage/storages/GetProducts"
	var product Product
	err := p.StmtGet.QueryRow(id).Scan(&product.Name, &product.Price, &product.Amount)
	if errors.Is(err, sql.ErrNoRows) {
		return Product{}, fmt.Errorf("%s, %w", op, errs.ErrProductNotFound)
	}
	if err != nil {
		return Product{}, fmt.Errorf("%s: execute statement %w", op, err)
	}
	product.Id = id
	return product, nil
}

func (p *ProductStorage) DeleteProduct(id int64) error {
	const op = "internal/storage/storages/DeleteProduct"
	res, err := p.StmtDelete.Exec(id)
	if err != nil {
		return fmt.Errorf("%s: execute statement %w", op, err)
	}
	row, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("%s: execute statement %w", op, err)
	}
	if row == 0 {
		return fmt.Errorf("%s: execute statement %w", op, errs.ErrProductNotFound)
	}
	return nil
}
