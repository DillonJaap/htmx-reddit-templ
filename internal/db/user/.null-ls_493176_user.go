package main

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
	ID             int
	name           string
	hashedPassword string
	TimeCreated    time.Time
}

type Model interface {
	Get(int) (User, error)
	Add(User) (int, error)
	Delete(int) error
}

type model struct {
	DB *sql.DB
}

var _ Model = &model{}

func createTable(DB *sql.DB) {
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

func New(db *sql.DB) Model {
	createTable(db)
	return &model{db}
}

func (m *model) Get(id int) (User, error) {
	var user User

	row := m.DB.QueryRow(
		"SELECT id, name, created_time FROM post WHERE id=?", id,
	)
	err := row.Scan(
		&user.ID,
		&user.name,
		&user.TimeCreated,
	)
	if err != nil {
		return user, err
	}

	return user, nil
}

func (t *model) Add(user User) (int, error) {
	result, err := t.DB.Exec(
		"INSERT INTO user (name, password, created_time) VALUES (?,?,?);",
		user.name,
		user.hashedPassword,
		time.Now(),
	)
	if err != nil {
		return 0, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), err
}

func (t *model) Delete(id int) error {
	_, err := t.DB.Exec("DELETE FROM user WHERE id = ?", id)
	if err != nil {
		return err
	}
	return nil
}

func (m *model) Insert(name, email, password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return err
	}

	stmt := `INSERT INTO users (name, email, hashed_password, created)
	VALUES(?, ?, ?, UTC_TIMESTAMP())`

	_, err = m.DB.Exec(stmt, name, email, string(hashedPassword))
	if err != nil {
		return err
	}
	return nil
}

func (m *model) Authenticate(name, password string) (int, error) {
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

func (m *model) Exists(id int) (bool, error) {
	var exists bool

	stmt := "SELECT EXISTS(SELECT true FROM users WHERE id = ?)"

	err := m.DB.QueryRow(stmt, id).Scan(&exists)
	return exists, err
}
