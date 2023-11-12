package db

import (
	"database/sql"
	"time"

	"github.com/charmbracelet/log"
)

type Post struct {
	ID          int
	Title       string
	Body        string
	Owner       string
	OwnerID     int
	TimeCreated time.Time
}

type PostStore interface {
	Get(int) (Post, error)
	GetAll() ([]Post, error)
	Add(Post) (int, error)
	Delete(int) error
}

type postModel struct {
	DB *sql.DB
}

var _ PostStore = &postModel{}

func NewPostStore(db *sql.DB) PostStore {
	return &postModel{db}
}

func (m *postModel) Get(id int) (Post, error) {
	var post Post

	row := m.DB.QueryRow(
		"SELECT id, title, body, owner, created_time FROM post WHERE id=?", id,
	)
	err := row.Scan(
		&post.ID,
		&post.Title,
		&post.Body,
		&post.Owner,
		&post.TimeCreated,
	)
	if err != nil {
		return post, err
	}

	return post, nil
}

func (t *postModel) GetAll() ([]Post, error) {
	var posts []Post
	var post Post

	rows, err := t.DB.Query("SELECT id, title, body, owner, created_time FROM post")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(
			&post.ID,
			&post.Title,
			&post.Body,
			&post.Owner,
			&post.TimeCreated,
		)
		if err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}

	return posts, nil
}

func (t *postModel) Add(post Post) (int, error) {
	log.Printf("%+v", post)
	result, err := t.DB.Exec(
		"INSERT INTO post (title, body, created_time, owner, owner_id) VALUES (?,?,?,?,?);",
		post.Title,
		post.Body,
		time.Now(),
		post.Owner,
		post.OwnerID,
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

func (t *postModel) Delete(id int) error {
	_, err := t.DB.Exec("DELETE FROM post WHERE id = ?", id)
	if err != nil {
		return err
	}
	return nil
}
