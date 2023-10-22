package post

import (
	"database/sql"
	"fmt"
	"time"
)

type Post struct {
	ID          int
	Title       string
	Body        string
	TimeCreated time.Time
}

type Model interface {
	Get(int) (Post, error)
	GetAll() ([]Post, error)
	Add(Post) (int, error)
	Delete(int) error
}

type model struct {
	DB *sql.DB
}

var _ Model = &model{}

func createTable(DB *sql.DB) {
	_, err := DB.Exec(`
	CREATE TABLE IF NOT EXISTS post (
		id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
		title TEXT NOT NULL,
		body TEXT NOT NULL,
		created_time DATETIME NOT NULL
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

func (m *model) Get(id int) (Post, error) {
	var post Post

	row := m.DB.QueryRow(
		"SELECT id, title, body, created_time FROM post WHERE id=?", id,
	)
	err := row.Scan(
		&post.ID,
		&post.Title,
		&post.Body,
		&post.TimeCreated,
	)
	if err != nil {
		return post, err
	}

	return post, nil
}

func (t *model) GetAll() ([]Post, error) {
	var posts []Post
	var post Post

	rows, err := t.DB.Query("SELECT id, title, body, created_time FROM post")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(
			&post.ID,
			&post.Title,
			&post.Body,
			&post.TimeCreated,
		)
		if err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}

	return posts, nil
}

func (t *model) Add(post Post) (int, error) {
	result, err := t.DB.Exec(
		"INSERT INTO post (title, body, created_time) VALUES (?,?,?);",
		post.Title,
		post.Body,
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
	_, err := t.DB.Exec("DELETE FROM post WHERE id = ?", id)
	if err != nil {
		return err
	}
	return nil
}
