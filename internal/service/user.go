package service

import (
	"htmx-reddit/internal/adapter"
	"htmx-reddit/internal/db"
	"time"

	"context"

	"github.com/alexedwards/scs/v2"
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
	Login(context.Context, string, string) error
	Logout(ctx context.Context) error
}

type user struct {
	db.UserStore // use UserStore methods unless specifically overridden
	sess         *scs.SessionManager
}

func NewUser(m db.UserStore, s *scs.SessionManager) User {
	return user{
		UserStore: m,
		sess:      s,
	}
}

func (u user) Get(id int) (UserData, error) {
	return adapter.Get("user", u.UserStore.Get, asUserData)(id)
}

func (u user) Add(name, pass string) error {
	return u.UserStore.Add(db.User{
		Name:     name,
		Password: pass,
	})
}

func (u user) Login(ctx context.Context, name, password string) error {
	id, err := u.UserStore.Authenticate(name, password)
	if err != nil {
		return err
	}

	u.sess.Put(ctx, "authenticatedUserID", id)
	u.sess.Put(ctx, "username", name)

	return nil
}

func (u user) Logout(ctx context.Context) error {
	err := u.sess.RenewToken(ctx)
	if err != nil {
		return err
	}

	u.sess.Remove(ctx, "authenticatedUserID")
	u.sess.Remove(ctx, "username")
	return nil
}
