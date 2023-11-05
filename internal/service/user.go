package service

import (
	"htmx-reddit/internal/adapter"
	"htmx-reddit/internal/db"
	"time"
)

type UserData struct {
	ID          int
	Name        string
	Password    string
	TimeCreated time.Time
}

func asUserData(u db.User) UserData {
	return UserData{
		ID:          u.ID,
		Name:        u.Name,
		TimeCreated: u.TimeCreated,
	}
}

type User interface {
	Get(int) (UserData, error)
	Delete(int) error
	Add(string, string) error
	Exists(id int) (bool, error)
}

type user struct {
	db.UserStore // use UserStore methods unless specifically overridden
}

func NewUser(m db.UserStore) User {
	return user{UserStore: m}
}

func (us user) Get(id int) (UserData, error) {
	return adapter.Get("user", us.UserStore.Get, asUserData)(id)
}

func (ps user) Add(name, pass string) error {
	return ps.UserStore.Add(db.User{
		Name:     name,
		Password: pass,
	})
}
