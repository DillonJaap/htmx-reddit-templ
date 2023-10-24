package service

import (
	"htmx-reddit/internal/adapter"
	"htmx-reddit/internal/db/post"
	"time"
)

type Post struct {
	ID          int
	Title       string
	Body        string
	TimeCreated time.Time
}

func asDBPost(dbPost post.Post) Post {
	return Post{
		ID:          dbPost.ID,
		Title:       dbPost.Title,
		Body:        dbPost.Body,
		TimeCreated: dbPost.TimeCreated,
	}
}

type PostService interface {
	Get(int) (Post, error)
	Delete(int) error
	Add(title, body string) error
}

type postService struct {
	m post.Model
}

func (ps postService) Get(id int) (Post, error) {
	return adapter.Get("post", ps.m.Get, asDBPost)(id)
}

func (ps postService) Delete(id int) error {
	return ps.m.Delete(id)
}

func (ps postService) Add(title, body string) error {
	_, err := ps.m.Add(post.Post{
		Title: title,
		Body:  body,
	})
	return err
}
