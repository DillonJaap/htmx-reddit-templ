package service

import (
	"htmx-reddit/internal/adapter"
	"htmx-reddit/internal/db/post"
	"os/user"
	"time"
)

type User struct {
	ID          int
	Name        string
	Password    string
	TimeCreated time.Time
}

func asViewData(dbUser user.User) UserData {
	return UserData{
		ID:          dbUser.ID,
		Name:        dbUser.Name,
		TimeCreated: dbUser.TimeCreated,
	}
}

type UserService interface {
	Get(int) (Post, error)
	GetAll() ([]Post, error)
	Delete(int) error
	Add(title, body string) error
}

type postService struct {
	m post.Model
}

func NewPostService(m post.Model) PostService {
	return postService{m: m}
}

func (ps postService) Get(id int) (Post, error) {
	return adapter.Get("post", ps.m.Get, asServicePost)(id)
}

func (ps postService) GetAll() ([]Post, error) {
	return adapter.GetAll("post", ps.m.GetAll, asServicePost)()
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
