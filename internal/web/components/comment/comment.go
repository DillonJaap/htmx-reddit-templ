package comment

import (
	"htmx-reddit/internal/models/comment"
	"time"

	"github.com/charmbracelet/log"
)

type Data struct {
	ID              int
	ParentID        int
	Description     string
	TimeCreated     time.Time
	Replies         []Data
	ReplyBoxVisible bool
}

type Controller struct {
	GetComments func() ([]Data, error)
	GetReplies  func(int) []Data
	Get         func(int) (Data, error)
	Add         func(string, int) (int, error)
	Delete      func(int) error
}

func asViewData(dbCmnt comment.Comment) Data {
	return Data{
		ID:              dbCmnt.ID,
		ParentID:        dbCmnt.ParentID,
		Description:     dbCmnt.Description,
		TimeCreated:     dbCmnt.TimeCreated,
		ReplyBoxVisible: false,
	}
}

func (data Data) asDBData() comment.Comment {
	return comment.Comment{
		ID:       data.ID,
		ParentID: data.ParentID,
		//PostID:      data.PostID, // TODO fill this in
		Description: data.Description,
		TimeCreated: data.TimeCreated,
	}
}

func GetByParentID(model comment.Model, id int) []Data {
	dbComments, err := model.GetByParentID(id)
	if err != nil {
		log.Error("getting comments", "error", err)
		return nil
	}

	var comments []Data

	for _, dbCmnt := range dbComments {
		comment := asViewData(dbCmnt)
		comment.Replies = GetByParentID(model, comment.ID)
		comments = append(comments, comment)
	}

	return comments
}

func GetByPostID(model comment.Model, id int) []Data {
	dbComments, err := model.GetByPostID(id)
	if err != nil {
		log.Error("getting comments", "error", err)
		return nil
	}

	var comments []Data

	for _, dbCmnt := range dbComments {
		comment := asViewData(dbCmnt)
		comment.Replies = GetByParentID(model, comment.ID)
		comments = append(comments, comment)
	}

	return comments
}
