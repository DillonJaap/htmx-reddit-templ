package service

import (
	"htmx-reddit/internal/adapter"
	"htmx-reddit/internal/db"
	"time"

	"github.com/charmbracelet/log"
)

type CommentData struct {
	ID          int
	ParentID    int
	Description string
	TimeCreated time.Time
	Replies     []CommentData
}

func asCommentData(c db.Comment) CommentData {
	return CommentData{
		ID:          c.ID,
		ParentID:    c.ParentID,
		Description: c.Description,
		TimeCreated: c.TimeCreated,
	}
}

type Comment interface {
	Get(int) (CommentData, error)
	GetAll() ([]CommentData, error)
	Delete(int) error
	Add(int, string) (int, error)
	GetByParentID(id int) []CommentData
	GetByPostID(id int) []CommentData
}

type comment struct {
	db.CommentStore // use CommentStore methods unless specifically overridden
}

func NewComment(m db.CommentStore) Comment {
	return comment{CommentStore: m}
}

func (cs comment) Get(id int) (CommentData, error) {
	return adapter.Get("comment", cs.CommentStore.Get, asCommentData)(id)
}

func (cs comment) GetAll() ([]CommentData, error) {
	return adapter.GetAll("comment", cs.CommentStore.GetAll, asCommentData)()
}

func (cs comment) Add(parentID int, desc string) (int, error) {
	return cs.CommentStore.Add(db.Comment{
		ParentID:    parentID,
		Description: desc,
	})
}

func (cs comment) GetByParentID(id int) []CommentData {
	dbComments, err := cs.CommentStore.GetByParentID(id)
	if err != nil {
		log.Error("getting comments", "error", err)
		return nil
	}

	var comments []CommentData

	for _, dbCmnt := range dbComments {
		comment := asCommentData(dbCmnt)
		comment.Replies = cs.GetByParentID(comment.ID)
		comments = append(comments, comment)
	}

	return comments
}

func (cs comment) GetByPostID(id int) []CommentData {
	dbComments, err := cs.CommentStore.GetByPostID(id)
	if err != nil {
		log.Error("getting comments", "error", err)
		return nil
	}

	var comments []CommentData

	for _, dbCmnt := range dbComments {
		comment := asCommentData(dbCmnt)
		comment.Replies = cs.GetByParentID(comment.ID)
		comments = append(comments, comment)
	}

	return comments
}
