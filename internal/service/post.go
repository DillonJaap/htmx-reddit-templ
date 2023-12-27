package service

import (
	"context"
	"errors"
	"fmt"
	"htmx-reddit/internal/adapter"
	"htmx-reddit/internal/db"
	"htmx-reddit/internal/helpers"
	err_ext "htmx-reddit/lib/errors_ext"
	"time"

	"github.com/alexedwards/scs/v2"
)

var (
	ErrNotLoggedIn          = errors.New("user is not logged in")
	ErrUnauthorizedToDelete = errors.New("unauthorized to delete post")
	ErrSqlError             = errors.New("")
)

type PostData struct {
	ID          int
	Title       string
	Body        string
	Owner       string
	OwnerID     int
	TimeCreated time.Time
}

func asPostData(p db.Post) PostData {
	return PostData{
		ID:          p.ID,
		Title:       p.Title,
		Body:        p.Body,
		Owner:       p.Owner,
		OwnerID:     p.OwnerID,
		TimeCreated: p.TimeCreated,
	}
}

type Post interface {
	Get(int) (PostData, error)
	GetAll() ([]PostData, error)
	Delete(context.Context, int) error
	Add(context.Context, string, string) error
}

type post struct {
	db.PostStore // use PostStore methods unless specifically overridden
	sess         *scs.SessionManager
}

func NewPost(m db.PostStore, s *scs.SessionManager) Post {
	return post{
		PostStore: m,
		sess:      s,
	}
}

func (p post) Get(id int) (PostData, error) {
	return adapter.Get("post", p.PostStore.Get, asPostData)(id)
}

func (p post) GetAll() ([]PostData, error) {
	return adapter.GetAll("post", p.PostStore.GetAll, asPostData)()
}

func (p post) Add(ctx context.Context, title, body string) error {
	if !helpers.IsLoggedIn(ctx, p.sess) {
		return ErrNotLoggedIn
	}
	_, err := p.PostStore.Add(db.Post{
		Title:   title,
		Body:    body,
		Owner:   p.sess.GetString(ctx, "username"),
		OwnerID: p.sess.GetInt(ctx, "authenticatedUserID"),
	})
	if err != nil {
		return err_ext.Join(ErrSqlError, err)
	}
	return nil
}

func (p post) Delete(ctx context.Context, id int) error {
	if !helpers.IsLoggedInUser(ctx, p.sess, id) {
		return fmt.Errorf("%w: LoggedInUserid is: %d", ErrUnauthorizedToDelete, id)
	}

	return p.PostStore.Delete(id)
}
