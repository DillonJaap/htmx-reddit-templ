package components

import (
	"htmx-reddit/internal/service"
)

type Handler struct {
	Comment *comment
	Post    *post
	User    *user
}

func NewHandler(
	comments service.Comment,
	posts service.Post,
	users service.User,
) *Handler {
	return &Handler{
		Comment: newComment(comments),
		Post:    newPost(posts),
		User:    newUser(users),
	}
}
