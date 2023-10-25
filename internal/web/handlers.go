package handlers

import (
	"htmx-reddit/internal/service"
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
	comments service.Comment,
	posts service.Post,
	users service.User,
) *Handler {
	return &Handler{
		Comment: cmntHandler.New(comments),
		Post:    postHandler.New(posts),
		User:    userHandler.New(users),
	}
}
