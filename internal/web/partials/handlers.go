package components

import (
	"htmx-reddit/internal/service"

	"github.com/alexedwards/scs/v2"
)

type Handler struct {
	Comment *comment
	Post    *post
	User    *user
}

func NewHandler(
	sess *scs.SessionManager,
	comments service.Comment,
	posts service.Post,
	users service.User,
) *Handler {
	return &Handler{
		Comment: newComment(comments, sess),
		Post:    newPost(posts),
		User:    newUser(users),
	}
}
