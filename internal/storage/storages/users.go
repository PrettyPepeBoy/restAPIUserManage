package storages

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/mattn/go-sqlite3"
	"tstUser/internal/storage/storages/errs"
)

type UserStorage struct {
	db              *sql.DB
	StmtDelete      *sql.Stmt
	StmtCreate      *sql.Stmt
	StmtFindUser    *sql.Stmt
	StmtCheckUserID *sql.Stmt
	StmtUpdateUser  *sql.Stmt
}

type User struct {
	Id      int64
	Name    string
	Surname string
	Mail    string
	Date    string
	Cash    int
}

func NewUserTable(storagePath string) (*UserStorage, error) {
	const op = "data-logic/pack/storage/storages/New"

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

	stmtFindUser, err := db.Prepare("SELECT id, name, surname, mail, cash, date  FROM user WHERE id = ?")
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	stmtUpdateUser, err := db.Prepare("UPDATE user SET name = ?, surname = ?, mail = ?, cash = ?, date = ? WHERE id = ?")
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &UserStorage{db: db,
		StmtDelete:      stmtDelete,
		StmtCreate:      stmtCreate,
		StmtCheckUserID: stmtCheckUserId,
		StmtFindUser:    stmtFindUser,
		StmtUpdateUser:  stmtUpdateUser,
	}, nil
}

func (s *UserStorage) CreateUser(user User) (int64, error) {
	const op = "storage/storages/CreateUser"
	res, err := s.StmtCreate.Exec(user.Name, user.Surname, user.Mail, user.Cash, user.Date)
	if err != nil {
		var sqliteErr sqlite3.Error
		if errors.As(err, &sqliteErr) && errors.Is(sqliteErr.ExtendedCode, sqlite3.ErrConstraintUnique) {
			return 0, fmt.Errorf("%s, %w", op, errs.ErrUserExist)
		}
		return 0, fmt.Errorf("%s: execute statemnt: %w", op, err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("%s: failed to get last insert id: %w", op, err)
	}

	return id, nil
}

func (s *UserStorage) DeleteUser(id int64) error {
	const op = "storage/storages/DeleteUser"
	res, err := s.StmtDelete.Exec(id)
	if err != nil {
		return fmt.Errorf("%s, %w", op, err)
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("%s, %w", op, err)
	}
	if rows == 0 {
		return fmt.Errorf("%s, %w", op, errs.ErrUserNotFound)
	}
	return nil
}

func (s *UserStorage) CheckUserID(Id int64) (string, error) {
	const op = "storage/storages/CheckUserID"
	var userMail string
	err := s.StmtCheckUserID.QueryRow(Id).Scan(&userMail)
	if errors.Is(err, sql.ErrNoRows) {
		return "", errs.ErrUserNotFound
	}
	if err != nil {
		return "", fmt.Errorf("%s: execute statement: %w", op, err)
	}
	return userMail, nil
}

func (s *UserStorage) GetUserInfo(Id int64) (User, error) {
	const op = "storage/storages/GetUser"
	var user User
	err := s.StmtFindUser.QueryRow(Id).Scan(&user.Id, &user.Name, &user.Surname, &user.Mail, &user.Cash, &user.Date)
	if errors.Is(err, sql.ErrNoRows) {
		return User{}, errs.ErrUserNotFound
	}
	if err != nil {
		return User{}, fmt.Errorf("%s: execute statement: %w", op, err)
	}
	return user, nil
}

func (s *UserStorage) UpdateUser(user User) error {
	const op = "storage/storages/UpdateUser"
	res, err := s.StmtUpdateUser.Exec(user.Name, user.Surname, user.Mail, user.Cash, user.Date, user.Id)
	if err != nil {
		return fmt.Errorf("%s: execute statment: %w", op, err)
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("%s: execute statment: %w", op, err)
	}
	if rows == 0 {
		return fmt.Errorf("%s, %w", op, errs.ErrUserNotFound)
	}
	return nil
}
