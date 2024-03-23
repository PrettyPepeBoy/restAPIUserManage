package sqlite

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/mattn/go-sqlite3"
	"tstUser/internal/storage"
)

type Goods struct {
	db         *sql.DB
	StmtCreate *sql.Stmt
}

func NewStuffTable(storagePath string) (*Goods, error) {
	const op = "internal/storage/sqlite/NewStuffTable"

	db, err := sql.Open("sqlite3", storagePath)

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	stmt, err := db.Prepare(`
	CREATE TABLE IF NOT EXISTS goods(
	    id INTEGER PRIMARY KEY,
	    name TEXT NOT NULL UNIQUE,
	    price INT NOT NULL,
	    amount INT);
	CREATE INDEX IF NOT EXISTS idx_name  ON goods(name)
	`)

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	_, err = stmt.Exec()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	stmtCreate, err := db.Prepare(`INSERT INTO goods(name, price, amount) VALUES (?,?,?)`)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Goods{db: db,
		StmtCreate: stmtCreate,
	}, nil
}

func (g *Goods) CreateGoods(name string, price, amount int64) (int64, error) {
	const op = "internal/storage/sqlite/CreateGoods"
	res, err := g.StmtCreate.Exec(name, price, amount)
	if err != nil {
		var sqliteErr sqlite3.Error
		if errors.As(err, &sqliteErr) && errors.Is(sqliteErr.ExtendedCode, sqlite3.ErrConstraintUnique) {
			return 0, fmt.Errorf("%s, %w", op, storage.ErrGoodsExist)
		}
		return 0, fmt.Errorf("%s: execute statemnt: %w", op, err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("%s: failed to get last insert id: %w", op, err)
	}
	return id, nil

}
