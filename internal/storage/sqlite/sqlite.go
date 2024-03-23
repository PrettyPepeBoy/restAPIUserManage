package sqlite

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/mattn/go-sqlite3"
	"tstUser/internal/storage"
)

type Storage struct {
	db              *sql.DB
	StmtDelete      *sql.Stmt
	StmtCreate      *sql.Stmt
	StmtFindUser    *sql.Stmt
	StmtCheckUserID *sql.Stmt
}

func NewUserTable(storagePath string) (*Storage, error) {
	const op = "data-logic/pack/storage/sqlite/New"

	db, err := sql.Open("sqlite3", storagePath)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	stmt, err := db.Prepare(`
	CREATE TABLE IF NOT EXISTS user(
	    id INTEGER PRIMARY KEY,
	    name TEXT NOT NULL,
	    surname TEXT NOT NULL,
	    mail TEXT NOT NULL UNIQUE ,
	    cash INTEGER,
	    date TEXT NOT NULL );
	CREATE INDEX IF NOT EXISTS idx_mail ON user(mail)
`)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	_, err = stmt.Exec()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	stmtDelete, err := db.Prepare("DELETE FROM  user WHERE id = ?  ")
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	stmtCreate, err := db.Prepare("INSERT INTO user(name, surname, mail, cash, date) VALUES(?,?,?,?,?)")
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	stmtCheckUserId, err := db.Prepare("SELECT mail FROM user WHERE id = ?")
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	stmtFindUser, err := db.Prepare("SELECT mail FROM user WHERE mail = ?")
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Storage{db: db,
		StmtDelete:      stmtDelete,
		StmtCreate:      stmtCreate,
		StmtCheckUserID: stmtCheckUserId,
		StmtFindUser:    stmtFindUser,
	}, nil
}

func (s *Storage) CreateUser(name, surname, mail, date string, cash int) (int64, error) {
	const op = "storage/sqlite/CreateUser"
	res, err := s.StmtCreate.Exec(name, surname, mail, cash, date)
	if err != nil {
		var sqliteErr sqlite3.Error
		if errors.As(err, &sqliteErr) && errors.Is(sqliteErr.ExtendedCode, sqlite3.ErrConstraintUnique) {
			return 0, fmt.Errorf("%s, %w", op, storage.ErrUserExist)
		}
		return 0, fmt.Errorf("%s: execute statemnt: %w", op, err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("%s: failed to get last insert id: %w", op, err)
	}

	return id, nil
}

func (s *Storage) DeleteUser(id int64) error {
	const op = "storage/sqlite/DeleteUser"
	res, err := s.StmtDelete.Exec(id)
	if err != nil {
		return fmt.Errorf("%s, %w", op, err)
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("%s, %w", op, err)
	}
	if rows == 0 {
		return fmt.Errorf("%s, %w", op, storage.ErrUserNotFound)
	}
	return nil
}

func (s *Storage) CheckUserID(Id int64) (string, error) {
	const op = "storage/sqlite/CheckUserID"
	var userMail string
	err := s.StmtCheckUserID.QueryRow(Id).Scan(&userMail)
	if errors.Is(err, sql.ErrNoRows) {
		return "", storage.ErrUserNotFound
	}
	if err != nil {
		return "", fmt.Errorf("%s: execute statement: %w", op, err)
	}
	return userMail, nil
}

func (s *Storage) GetUser(userMail string) error {
	const op = "storage/sqlite/GetUser"
	var mail string
	err := s.StmtFindUser.QueryRow(userMail).Scan(mail)
	if errors.Is(err, sql.ErrNoRows) {
		return storage.ErrUserNotFound
	}
	if err != nil {
		return fmt.Errorf("%s: execute statement: %w", op, err)
	}
	return nil
}
