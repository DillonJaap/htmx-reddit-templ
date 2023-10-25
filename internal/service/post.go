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
	m db.PostStore
}

func NewPost(m db.PostStore) Post {
	return post{m: m}
}

func (ps post) Get(id int) (PostData, error) {
	return adapter.Get("post", ps.m.Get, asPostData)(id)
}

func (ps post) GetAll() ([]PostData, error) {
	return adapter.GetAll("post", ps.m.GetAll, asPostData)()
}

func (ps post) Delete(id int) error {
	return ps.m.Delete(id)
}

func (ps post) Add(title, body string) error {
	_, err := ps.m.Add(db.Post{
		Title: title,
		Body:  body,
	})
	return err
}
