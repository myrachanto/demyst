package rating

import (
	httperrors "github.com/myrachanto/erroring"
	"github.com/myrachanto/estate/src/support"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Rating struct {
	ID          primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Author      string             `json:"author"`
	Productname string             `json:"productname"`
	Productcode string             `json:"productcode"`
	Description string             `json:"description"`
	Code        string             `json:"code"`
	Rate        int64              `json:"rate"`
	Shopalias   string             `json:"shopalias"`
	Used        bool               `json:"used"`
	Featured    bool               `json:"featured"`
	Base        support.Base
}

func (t Rating) Validate() httperrors.HttpErr {
	if t.Productcode == "" {
		return httperrors.NewBadRequestError("Invalid Code")
	}
	if t.Description == "" {
		return httperrors.NewBadRequestError("Invalid Description")
	}
	if t.Rate < 0 {
		return httperrors.NewBadRequestError("Invalid Rate")
	}
	return nil
}
