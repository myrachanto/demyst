package location

import (
	httperrors "github.com/myrachanto/erroring"
	"github.com/myrachanto/estate/src/support"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Location struct {
	ID           primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Name         string             `json:"name,omitempty"`
	Title        string             `json:"title,omitempty"`
	Description  string             `json:"description,omitempty"`
	Code         string             `json:"code,omitempty"`
	Picture      string             `json:"picture"`
	Url          string             `json:"url"`
	Meta         string             `json:"meta"`
	Content      string             `json:"content"`
	PropertyType string             `json:"propertyType"`
	Base         support.Base       `json:"base,omitempty"`
}

type Results struct {
	Data        []Location `json:"results"`
	Total       int        `json:"total"`
	Pages       int        `json:"pages"`
	CurrentPage int        `json:"currentpage"`
}

var MessageResp struct {
	Message string `json:"message,omitempty"`
}

func (l Location) Validate() httperrors.HttpErr {
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
