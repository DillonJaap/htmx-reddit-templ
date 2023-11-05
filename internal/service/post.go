package service

import (
	"htmx-reddit/internal/adapter"
	"htmx-reddit/internal/db"
	"time"
)

type PostData struct {
	ID          int
	Title       string
	Body        string
	TimeCreated time.Time
}

func asPostData(p db.Post) PostData {
	return PostData{
		ID:          p.ID,
		Title:       p.Title,
		Body:        p.Body,
		TimeCreated: p.TimeCreated,
	}
}

type Post interface {
	Get(int) (PostData, error)
	GetAll() ([]PostData, error)
	Delete(int) error
	Add(title, body string) error
}

type post struct {
	db.PostStore // use PostStore methods unless specifically overridden
}

func NewPost(m db.PostStore) Post {
	return post{PostStore: m}
}

func (ps post) Get(id int) (PostData, error) {
	return adapter.Get("post", ps.PostStore.Get, asPostData)(id)
}

func (ps post) GetAll() ([]PostData, error) {
	return adapter.GetAll("post", ps.PostStore.GetAll, asPostData)()
}

func (ps post) Add(title, body string) error {
	_, err := ps.PostStore.Add(db.Post{
		Title: title,
		Body:  body,
	})
	return err
}
