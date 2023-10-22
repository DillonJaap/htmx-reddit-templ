package post

import (
	"htmx-reddit/internal/adapter"
	"htmx-reddit/internal/db/post"
	"htmx-reddit/internal/templ"
)

func AsViewData(dbPost post.Post) templ.PostInput {
	return templ.PostInput{
		ID:          dbPost.ID,
		Title:       dbPost.Title,
		Body:        dbPost.Body,
		TimeCreated: dbPost.TimeCreated,
		DeleteButton: templ.DeleteButtonInput{
			ID:         dbPost.ID,
			Target:     "post",
			DeletePath: "/post/delete",
		},
		EditButton: templ.EditButtonInput{
			ID:       dbPost.ID,
			Target:   "post",
			EditPath: "/post/edit",
		},
	}
}

func GetAll(model post.Model) ([]templ.PostInput, error) {
	return adapter.GetAll("post", model.GetAll, AsViewData)()
}

func Get(model post.Model, id int) (templ.PostInput, error) {
	return adapter.Get("post", model.Get, AsViewData)(id)
}
