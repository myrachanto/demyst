package business

import (
	"github.com/myrachanto/demyst/src/support"
	httperrors "github.com/myrachanto/erroring"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Business struct {
	ID              primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Name            string             `json:"name,omitempty"`
	BusinessPin     string             `json:"business_pin,omitempty"`
	YearEstablished int                `json:"year_established,omitempty"`
	Industry        string             `json:"industry,omitempty"`
	Code            string             `json:"code,omitempty"`
	Base            support.Base       `json:"base,omitempty"`
}

func (b Business) Validate() httperrors.HttpErr {
	if b.Name == "" {
		return httperrors.NewBadRequestError("Name should not be empty")
	}
	if b.BusinessPin == "" {
		return httperrors.NewBadRequestError("Business Pin should not be empty")
	}
	if b.YearEstablished == 0 {
		return httperrors.NewBadRequestError("Year of Establishment should not be empty")
	}
	return nil
}
