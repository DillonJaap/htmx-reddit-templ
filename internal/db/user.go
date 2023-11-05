package db

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidCredentials = errors.New("Invalid Credentails")
)

type User struct {
	ID          int
	Name        string
	Password    string
	TimeCreated time.Time
}

type UserStore interface {
	Get(int) (User, error)
	Add(User) error
	Delete(int) error
	Exists(id int) (bool, error)
}

type userModel struct {
	DB *sql.DB
}

var _ UserStore = &userModel{}

func createUserTable(DB *sql.DB) {
	_, err := DB.Exec(`
	CREATE TABLE IF NOT EXISTS user (
		id              INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
		name            TEXT NOT NULL,
		hashed_password TEXT NOT NULL,
		created_time    DATETIME NOT NULL
	);
	`)
	if err != nil {
		fmt.Printf("error creating table: %s", err)
	}
}

func NewUserStore(db *sql.DB) UserStore {
	createUserTable(db)
	return &userModel{db}
}

func (m *userModel) Get(id int) (User, error) {
	var user User

	row := m.DB.QueryRow(
		"SELECT id, name, created_time FROM post WHERE id=?", id,
	)
	err := row.Scan(
		&user.ID,
		&user.Name,
		&user.TimeCreated,
	)
	if err != nil {
		return user, err
	}

	return user, nil
}

func (t *userModel) Delete(id int) error {
	_, err := t.DB.Exec("DELETE FROM user WHERE id = ?", id)
	if err != nil {
		return err
	}
	return nil
}

func (m *userModel) Add(user User) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), 12)
	if err != nil {
		return err
	}

	stmt := `INSERT INTO user (name, hashed_password, created_time)
	VALUES(?, ?, ?)`

	_, err = m.DB.Exec(
		stmt,
		user.Name,
		string(hashedPassword),
		time.Now(),
	)
	if err != nil {
		return err
	}
	return nil
}

func (m *userModel) Authenticate(name, password string) (int, error) {
	var id int
	var hashedPassword []byte

	stmt := "SELECT id, hashed_password FROM user WHERE name = ?"

	err := m.DB.QueryRow(stmt, name).Scan(&id, &hashedPassword)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, ErrInvalidCredentials
		}
		return 0, err
	}

	err = bcrypt.CompareHashAndPassword(hashedPassword, []byte(password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return 0, ErrInvalidCredentials
		}
		return 0, err
	}

	return id, nil
}

func (m *userModel) Exists(id int) (bool, error) {
	var exists bool

	stmt := "SELECT EXISTS(SELECT true FROM user WHERE id = ?)"

	err := m.DB.QueryRow(stmt, id).Scan(&exists)
	return exists, err
}
