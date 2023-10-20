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

func asViewData(dbUser user.User) UserData {
	return UserData{
		ID:          dbUser.ID,
		Name:        dbUser.Name,
		TimeCreated: dbUser.TimeCreated,
	}
}
