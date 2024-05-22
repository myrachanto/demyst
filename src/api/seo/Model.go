package seo

import (
	httperrors "github.com/myrachanto/erroring"
	"github.com/myrachanto/estate/src/support"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Seo struct {
	ID          primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Page        string             `json:"page"`
	Name        string             `json:"name"`
	Title       string             `json:"title"`
	Description string             `json:"description"`
	Code        string             `json:"code"`
	Meta        string             `json:"meta"`
	Kind        string             `json:"kind"`
	Location    string             `json:"location"`
	Sublocation string             `json:"sublocation"`
	Base        support.Base       `json:"base"`
}

type Results struct {
	Data        []Seo `json:"results"`
	Total       int   `json:"total"`
	Pages       int   `json:"pages"`
	CurrentPage int   `json:"currentpage"`
}

var MessageResp struct {
	Message string `json:"message,omitempty"`
}

func (l Seo) Validate() httperrors.HttpErr {
	if l.Title == "" {
		return httperrors.NewBadRequestError("Title should not be empty")
	}
	if l.Description == "" {
		return httperrors.NewBadRequestError("Description should not be empty")
	}
	return nil
}
