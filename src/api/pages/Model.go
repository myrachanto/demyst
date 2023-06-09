package pages

import (
	httperrors "github.com/myrachanto/erroring"
	"github.com/myrachanto/sports/src/support"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Page struct {
	ID      primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Name    string             `json:"name,omitempty"`
	Title   string             `json:"title,omitempty"`
	Meta    string             `json:"meta,omitempty"`
	Content string             `json:"content,omitempty"`
	Code    string             `json:"code,omitempty"`
	Base    support.Base       `json:"base,omitempty"`
}

type Results struct {
	Data        []*Page `json:"results"`
	Total       int     `json:"total"`
	Pages       int     `json:"pages"`
	CurrentPage int     `json:"currentpage"`
}

var MessageResp struct {
	Message string `json:"message,omitempty"`
}

func (p Page) Validate() httperrors.HttpErr {
	if p.Name == "" {
		return httperrors.NewBadRequestError("Name should not be empty")
	}
	if p.Title == "" {
		return httperrors.NewBadRequestError("Title should not be empty")
	}
	if p.Content == "" {
		return httperrors.NewBadRequestError("Content should not be empty")
	}
	return nil
}
