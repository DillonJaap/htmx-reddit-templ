package handlers

import (
	"htmx-reddit/internal/models/comment"
	"htmx-reddit/internal/models/post"
	"htmx-reddit/internal/models/user"
	cmntHandler "htmx-reddit/internal/web/components/comment"
	postHandler "htmx-reddit/internal/web/components/post"
	userHandler "htmx-reddit/internal/web/components/user"
)

type Handler struct {
	Comment *cmntHandler.Handler
	Post    *postHandler.Handler
	User    *userHandler.Handler
}

func New(
	comments comment.Model,
	posts post.Model,
	users user.Model,
) *Handler {
	return &Handler{
		Comment: cmntHandler.New(comments),
		Post:    postHandler.New(posts),
		User:    userHandler.New(users),
	}
}
