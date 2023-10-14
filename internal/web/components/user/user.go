package user

import (
	"htmx-reddit/internal/models/user"
	"time"
)

type UserData struct {
	ID          int
	Name        string
	Password    string
	TimeCreated time.Time
}

type Controller struct {
	Get    func(int) (UserData, error)
	Add    func(string, int) (int, error)
	Delete func(int) error
}

func asViewData(dbUser user.User) UserData {
	return UserData{
		ID:          dbUser.ID,
		Name:        dbUser.Name,
		TimeCreated: dbUser.TimeCreated,
	}
}

func (data UserData) asDBData() user.User {
	return user.User{
		ID:          data.ID,
		Password:    data.Password,
		TimeCreated: data.TimeCreated,
	}
}
