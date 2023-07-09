package news

import (
	httperrors "github.com/myrachanto/erroring"
	newssections "github.com/myrachanto/sports/src/api/newsSections"
	"github.com/myrachanto/sports/src/support"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type News struct {
	ID            primitive.ObjectID          `json:"_id,omitempty" bson:"_id,omitempty"`
	Name          string                      `json:"name"`
	Title         string                      `json:"title"`
	Caption       string                      `json:"caption"`
	Meta          string                      `json:"meta"`
	Url           string                      `json:"url"`
	Sport         string                      `json:"sport"`
	Featured      bool                        `json:"featured"`
	Exclusive     bool                        `json:"exclusive"`
	Trending      bool                        `json:"trending"`
	Content       string                      `json:"content"`
	SectionsCodes []string                    `json:"sectionscodes"`
	Sections      []*newssections.NewsSection `bson:"-" json:"sections"`
	Picture       string                      `json:"picture"`
	Code          string                      `json:"code"`
	TimeAgo       string                      `json:"time_ago"`
	Author        string                      `json:"author"`
	Credit        string                      `json:"credit"`
	PhotoCredit   string                      `json:"photocredit"`
	Comments      []Comment                   `json:"comments"`
	Base          support.Base                `json:"base"`
}

//	type NewsB struct {
//		ID        primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
//		Name      string             `json:"name"`
//		Title     string             `json:"title"`
//		Caption   string             `json:"caption"`
//		Meta      string             `json:"meta"`
//		Url       string             `json:"url"`
//		Sport     string             `json:"sport"`
//		Featured  bool               `json:"featured"`
//		Exclusive bool               `json:"exclusive"`
//		Trending  bool               `json:"trending"`
//		Content   string             `json:"content"`
//		Sections  []Coder            `json:"sections"`
//		Picture   string             `json:"picture"`
//		Code      string             `json:"code"`
//		TimeAgo   string             `json:"time_ago"`
//		Comments  []Comment          `json:"comments"`
//		Base      support.Base       `json:"base"`
//	}
type ByNews struct {
	News     *News   `json:"news"`
	Trending []*News `json:"trending"`
}
type Coder struct {
	Name string `json:"name,omitempty"`
}
type Section struct {
	Name      string `json:"name"`
	Content   string `json:"content"`
	Image     string `json:"image"`
	Highlight bool   `json:"highlight"`
}
type Comment struct {
	Name    string `json:"name,omitempty"`
	Message string `json:"message,omitempty"`
	TimeAgo string `json:"time_ago,omitempty"`
	Image   string `json:"image,omitempty"`
}
type Results struct {
	Data        []*News `json:"results"`
	Total       int     `json:"total"`
	Pages       int     `json:"pages"`
	CurrentPage int     `json:"currentpage"`
}

type SportCount struct {
	Sport string `bson:"_id" json:"sport" `
	Count int    `json:"count"`
}
type Weekly struct {
	Id    int32 `bson:"_id" json:"id"`
	Count int32 `bson:"count" json:"count"`
}

func (l News) Validate() httperrors.HttpErr {
	if l.Name == "" {
		return httperrors.NewBadRequestError("Name should not be empty")
	}
	if l.Title == "" {
		return httperrors.NewBadRequestError("Title should not be empty")
	}
	return nil
}
