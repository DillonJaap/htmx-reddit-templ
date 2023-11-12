package db

import (
	"database/sql"
	"time"
)

type Comment struct {
	ID          int
	ParentID    int
	PostID      int
	Description string
	Owner       string
	OwnerID     int
	TimeCreated time.Time
}

type CommentStore interface {
	Get(int) (Comment, error)
	GetAll() ([]Comment, error)
	GetByParentID(int) ([]Comment, error)
	GetByPostID(int) ([]Comment, error)
	Add(Comment) (int, error)
	Delete(int) error
}

type commentModel struct {
	DB *sql.DB
}

var _ CommentStore = &commentModel{}

func NewCommentStore(db *sql.DB) CommentStore {
	return &commentModel{db}
}

func (m *commentModel) Get(id int) (Comment, error) {
	var comment Comment

	row := m.DB.QueryRow(
		"SELECT id, description, created_time, parent_id, owner FROM comment WHERE id=?", id,
	)
	err := row.Scan(
		&comment.ID,
		&comment.Description,
		&comment.TimeCreated,
		&comment.ParentID,
		&comment.Owner,
	)
	if err != nil {
		return comment, err
	}

	return comment, nil
}

// TODO update this to actually get all comments
func (m *commentModel) GetAll() ([]Comment, error) {
	var comments []Comment
	var comment Comment

	// Get all top level comments
	rows, err := m.DB.Query("SELECT id, description, created_time, parent_id, owner FROM comment WHERE parent_id = 0")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(
			&comment.ID,
			&comment.Description,
			&comment.TimeCreated,
			&comment.ParentID,
			&comment.Owner,
		)
		if err != nil {
			return nil, err
		}
		comments = append(comments, comment)
	}

	return comments, nil
}

func (m *commentModel) GetByParentID(parentID int) ([]Comment, error) {
	var comments []Comment
	var comment Comment

	// Get all top level comments
	rows, err := m.DB.Query("SELECT id, description, created_time, parent_id, owner FROM comment WHERE parent_id=?", parentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(
			&comment.ID,
			&comment.Description,
			&comment.TimeCreated,
			&comment.ParentID,
			&comment.Owner,
		)
		if err != nil {
			return nil, err
		}
		comments = append(comments, comment)
	}

	return comments, nil
}

func (m *commentModel) GetByPostID(postID int) ([]Comment, error) {
	var comments []Comment
	var comment Comment

	// Get all top level comments
	rows, err := m.DB.Query("SELECT id, description, created_time, owner FROM comment WHERE post_id=?", postID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(
			&comment.ID,
			&comment.Description,
			&comment.TimeCreated,
			&comment.Owner,
		)
		if err != nil {
			return nil, err
		}
		comments = append(comments, comment)
	}

	return comments, nil
}
func (m *commentModel) Add(comment Comment) (int, error) {
	// TODO create SQL query to get reply depth of parent comment
	result, err := m.DB.Exec(
		"INSERT INTO comment (description, created_time, parent_id, post_id, owner, owner_id) VALUES (?,?,?,?,?,?);",
		comment.Description,
		time.Now(),
		comment.ParentID,
		comment.PostID,
		comment.Owner,
		comment.OwnerID,
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

func (m *commentModel) Delete(id int) error {
	_, err := m.DB.Exec("DELETE FROM comment WHERE id = ?", id)
	if err != nil {
		return err
	}
	return nil
}
