package subLocation

import (
	httperrors "github.com/myrachanto/erroring"
	"github.com/myrachanto/estate/src/support"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type SubLocation struct {
	ID           primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Name         string             `json:"name"`
	Title        string             `json:"title"`
	Description  string             `json:"description"`
	Code         string             `json:"code"`
	Url          string             `json:"url"`
	Location     string             `json:"location"`
	Meta         string             `json:"meta"`
	Content      string             `json:"content"`
	PropertyType string             `json:"propertyType"`
	Base         support.Base       `json:"base"`
}

type Results struct {
	Data        []SubLocation `json:"results"`
	Total       int           `json:"total"`
	Pages       int           `json:"pages"`
	CurrentPage int           `json:"currentpage"`
}

var MessageResp struct {
	Message string `json:"message,omitempty"`
}

func (l SubLocation) Validate() httperrors.HttpErr {
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
