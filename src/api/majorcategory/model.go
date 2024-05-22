package majorcategory

import (
	httperrors "github.com/myrachanto/erroring"
	"github.com/myrachanto/estate/src/support"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Majorcategory struct {
	ID            primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Name          string             `json:"name"`
	Title         string             `json:"title"`
	Description   string             `json:"description"`
	Meta          string             `json:"meta"`
	Content       string             `json:"content"`
	Category      []Categs           `json:"category"`
	Supercategory string             `json:"supercategory"`
	Code          string             `json:"code"`
	Shopalias     string             `json:"shopalias"`
	Picture       string             `json:"picture"`
	Url           string             `json:"url"`
	Used          bool               `json:"used"`
	Base          support.Base
}
type Results struct {
	Data        []Majorcategory `json:"results"`
	Total       int             `json:"total"`
	Pages       float32         `json:"pages"`
	CurrentPage int             `json:"currentpage"`
}

type Majorcat struct {
	Code string `json:"code"`
}
type Supercat struct {
	Code string `json:"code"`
}
type Categs struct {
	Code string `json:"code"`
}

func (majorcategory Majorcategory) Validate() httperrors.HttpErr {
	if majorcategory.Name == "" {
		return httperrors.NewNotFoundError("Invalid Name")
	}
	if majorcategory.Title == "" {
		return httperrors.NewNotFoundError("Invalid title")
	}
	// if majorcategory.Description == "" {
	// 	return httperrors.NewNotFoundError("Invalid Description")
	// }
	// if majorcategory.Shopalias == "" {
	// 	return httperrors.NewNotFoundError("Invalid Shopalias")
	// }
	return nil
}
