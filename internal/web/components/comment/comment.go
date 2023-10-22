package comment

import (
	"htmx-reddit/internal/db/comment"

	"htmx-reddit/internal/templ"

	"github.com/charmbracelet/log"
)

func asViewData(dbCmnt comment.Comment) templ.CommentInput {
	return templ.CommentInput{
		ID:              dbCmnt.ID,
		ParentID:        dbCmnt.ParentID,
		Description:     dbCmnt.Description,
		TimeCreated:     dbCmnt.TimeCreated,
		ReplyBoxVisible: false,
	}
}

func GetByParentID(model comment.Model, id int) []templ.CommentInput {
	dbComments, err := model.GetByParentID(id)
	if err != nil {
		log.Error("getting comments", "error", err)
		return nil
	}

	var comments []templ.CommentInput

	for _, dbCmnt := range dbComments {
		comment := asViewData(dbCmnt)
		comment.Replies = GetByParentID(model, comment.ID)
		comments = append(comments, comment)
	}

	return comments
}

func GetByPostID(model comment.Model, id int) []templ.CommentInput {
	dbComments, err := model.GetByPostID(id)
	if err != nil {
		log.Error("getting comments", "error", err)
		return nil
	}

	var comments []templ.CommentInput

	for _, dbCmnt := range dbComments {
		comment := asViewData(dbCmnt)
		comment.Replies = GetByParentID(model, comment.ID)
		comments = append(comments, comment)
	}

	return comments
}
