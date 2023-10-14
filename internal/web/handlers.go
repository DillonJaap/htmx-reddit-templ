package handlers

import (
	"htmx-reddit/internal/models/comment"
	"htmx-reddit/internal/models/post"
	"htmx-reddit/internal/models/user"
	"htmx-reddit/internal/render"
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
	renderer render.Renderer,
) *Handler {
	return &Handler{
		Comment: cmntHandler.New(comments, renderer),
		Post:    postHandler.New(posts, renderer),
		User:    userHandler.New(users, renderer),
	}
}
