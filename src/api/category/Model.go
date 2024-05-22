package category

import (
	httperrors "github.com/myrachanto/erroring"
	"github.com/myrachanto/estate/src/support"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Category struct {
	ID          primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Name        string             `json:"name,omitempty"`
	Title       string             `json:"title,omitempty"`
	Description string             `json:"description,omitempty"`
	Meta        string             `json:"meta"`
	Content     string             `json:"content"`
	Code        string             `json:"code,omitempty"`
	Base        support.Base       `json:"base,omitempty"`
}

type Results struct {
	Data        []Category `json:"results"`
	Total       int        `json:"total"`
	Pages       int        `json:"pages"`
	CurrentPage int        `json:"currentpage"`
}

var MessageResp struct {
	Message string `json:"message,omitempty"`
}

func (l Category) Validate() httperrors.HttpErr {
	if l.Name == "" {
		return httperrors.NewBadRequestError("Name should not be empty")
	}
	if l.Title == "" {
		return httperrors.NewBadRequestError("Title should not be empty")
	}
	// if l.Description == "" {
	// 	return httperrors.NewBadRequestError("Description should not be empty")
	// }
	return nil
}
