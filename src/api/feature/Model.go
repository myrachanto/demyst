package feature

import (
	httperrors "github.com/myrachanto/erroring"
	"github.com/myrachanto/estate/src/support"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Feature struct {
	ID          primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Name        string             `json:"name,omitempty"`
	Title       string             `json:"title,omitempty"`
	Description string             `json:"description,omitempty"`
	Code        string             `json:"code,omitempty"`
	Base        support.Base       `json:"base,omitempty"`
}

type Results struct {
	Data        []Feature `json:"results"`
	Total       int        `json:"total"`
	Pages       int        `json:"pages"`
	CurrentPage int        `json:"currentpage"`
}

var MessageResp struct {
	Message string `json:"message,omitempty"`
}

func (l Feature) Validate() httperrors.HttpErr {
	if l.Name == "" {
		return httperrors.NewBadRequestError("Name should not be empty")
	}
	if l.Title == "" {
		return httperrors.NewBadRequestError("Title should not be empty")
	}
	if l.Description == "" {
		return httperrors.NewBadRequestError("Description should not be empty")
	}
	return nil
}
