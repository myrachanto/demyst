package accounting

import (
	"github.com/myrachanto/demyst/src/support"
	httperrors "github.com/myrachanto/erroring"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Accounting struct {
	ID          primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Name        string             `json:"name,omitempty"`
	BusinessPin string             `json:"business_pin,omitempty"`
	UrlEndpoint string             `json:"urlEndpoint,omitempty"`
	Code        string             `json:"code,omitempty"`
	Base        support.Base       `json:"base,omitempty"`
}
type BalanceSheet struct {
	Year         int     `json:"year,omitempty"`
	Month        int     `json:"month,omitempty"`
	ProfitOrLoss float64 `json:"profit_or_loss,omitempty"`
	AssetsValue  float64 `json:"assets_value,omitempty"`
}

var MessageResp struct {
	Message string `json:"message,omitempty"`
}

func (l Accounting) Validate() httperrors.HttpErr {
	if l.Name == "" {
		return httperrors.NewBadRequestError("Name should not be empty")
	}
	if l.BusinessPin == "" {
		return httperrors.NewBadRequestError("Business Pin should not be empty")
	}
	if l.UrlEndpoint == "" {
		return httperrors.NewBadRequestError("Url endpoint should not be empty")
	}
	return nil
}
