package post

import (
	"htmx-reddit/internal/models/post"
	"time"
)

type Data struct {
	ID           int
	Title        string
	Body         string
	TimeCreated  time.Time
	DeleteButton deleteButton
	EditButton   editButton
}

func AsViewData(dbPost post.Post) Data {
	return Data{
		ID:          dbPost.ID,
		Title:       dbPost.Title,
		Body:        dbPost.Body,
		TimeCreated: dbPost.TimeCreated,
		DeleteButton: deleteButton{
			ID:         dbPost.ID,
			Target:     "post",
			DeletePath: "/post/delete",
		},
		EditButton: editButton{
			ID:       dbPost.ID,
			Target:   "post",
			EditPath: "/post/edit",
		},
	}
}

func (data Data) asDBData() post.Post {
	return post.Post{
		ID:          data.ID,
		Title:       data.Title,
		Body:        data.Body,
		TimeCreated: data.TimeCreated,
	}
}
