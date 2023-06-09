package tags

import (
	httperrors "github.com/myrachanto/erroring"
	"github.com/myrachanto/sports/src/support"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Tag struct {
	ID          primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Name        string             `json:"name,omitempty"`
	Title       string             `json:"title,omitempty"`
	Description string             `json:"description,omitempty"`
	Code        string             `json:"code,omitempty"`
	Base        support.Base       `json:"base,omitempty"`
}

type Results struct {
	Data        []*Tag `json:"results"`
	Total       int    `json:"total"`
	Pages       int    `json:"pages"`
	CurrentPage int    `json:"currentpage"`
}

var MessageResp struct {
	Message string `json:"message,omitempty"`
}

func (t Tag) Validate() httperrors.HttpErr {
	if t.Name == "" {
		return httperrors.NewBadRequestError("Name should not be empty")
	}
	if t.Title == "" {
		return httperrors.NewBadRequestError("Title should not be empty")
	}
	if t.Description == "" {
		return httperrors.NewBadRequestError("Description should not be empty")
	}
	return nil
}
