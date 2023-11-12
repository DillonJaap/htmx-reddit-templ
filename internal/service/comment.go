package service

import (
	"context"
	"errors"
	"htmx-reddit/internal/adapter"
	"htmx-reddit/internal/db"
	"htmx-reddit/internal/helpers"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/charmbracelet/log"
)

var (
	ErrUnauthorizedToDeleteComment = errors.New("unauthorized to delete comment")
)

type CommentData struct {
	ID          int
	ParentID    int
	Description string
	TimeCreated time.Time
	Owner       string
	OwnerID     int
	Replies     []CommentData
}

func asCommentData(c db.Comment) CommentData {
	return CommentData{
		ID:          c.ID,
		ParentID:    c.ParentID,
		Description: c.Description,
		Owner:       c.Owner,
		OwnerID:     c.OwnerID,
		TimeCreated: c.TimeCreated,
	}
}

type Comment interface {
	Get(int) (CommentData, error)
	GetAll() ([]CommentData, error)
	Add(context.Context, int, string) (int, error)
	Reply(context.Context, int, string) (int, error)
	GetByParentID(id int) []CommentData
	GetByPostID(id int) []CommentData
	Delete(context.Context, int) error
}

type comment struct {
	db.CommentStore // use CommentStore methods unless specifically overridden
	sess            *scs.SessionManager
}

func NewComment(m db.CommentStore, s *scs.SessionManager) Comment {
	return comment{
		CommentStore: m,
		sess:         s,
	}
}

func (c comment) Get(id int) (CommentData, error) {
	return adapter.Get("comment", c.CommentStore.Get, asCommentData)(id)
}

func (c comment) GetAll() ([]CommentData, error) {
	return adapter.GetAll("comment", c.CommentStore.GetAll, asCommentData)()
}

func (c comment) Add(ctx context.Context, postID int, desc string) (int, error) {
	return c.CommentStore.Add(db.Comment{
		PostID:      postID,
		Description: desc,
		Owner:       c.sess.GetString(ctx, "username"),
		OwnerID:     c.sess.GetInt(ctx, "authenticatedUserID"),
	})
}

func (c comment) Reply(ctx context.Context, parentID int, desc string) (int, error) {
	return c.CommentStore.Add(db.Comment{
		ParentID:    parentID,
		Description: desc,
		Owner:       c.sess.GetString(ctx, "username"),
		OwnerID:     c.sess.GetInt(ctx, "authenticatedUserID"),
	})
}

func (c comment) GetByParentID(id int) []CommentData {
	dbComments, err := c.CommentStore.GetByParentID(id)
	if err != nil {
		log.Error("getting comments", "error", err)
		return nil
	}

	var comments []CommentData

	for _, dbCmnt := range dbComments {
		comment := asCommentData(dbCmnt)
		comment.Replies = c.GetByParentID(comment.ID)
		comments = append(comments, comment)
	}

	return comments
}

func (c comment) GetByPostID(id int) []CommentData {
	dbComments, err := c.CommentStore.GetByPostID(id)
	if err != nil {
		log.Error("getting comments", "error", err)
		return nil
	}

	var comments []CommentData

	for _, dbCmnt := range dbComments {
		comment := asCommentData(dbCmnt)
		comment.Replies = c.GetByParentID(comment.ID)
		comments = append(comments, comment)
	}

	return comments
}

func (c comment) Delete(ctx context.Context, id int) error {
	if !helpers.IsLoggedInUser(ctx, c.sess, id) {
		return ErrUnauthorizedToDeleteComment
	}

	return c.CommentStore.Delete(id)
}
