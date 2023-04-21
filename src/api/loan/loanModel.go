package loan

import (
	"github.com/myrachanto/demyst/src/support"
	httperrors "github.com/myrachanto/erroring"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	PENDING  = "PENDING"
	DECLINED = "DECLINED"
	APPROVED = "APPROVED"
)

type Loan struct {
	ID                 primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Name               string             `json:"name,omitempty"`
	BusinessPin        string             `json:"business_pin,omitempty"`
	YearEstablished    int32              `json:"year_established,omitempty"`
	ProfitGrowth       float64            `json:"profit_growth,omitempty"`
	HighestProfit      float64            `json:"highest_profit,omitempty"`
	LastQuaterProfit   float64            `json:"last_quater_profit,omitempty"`
	LastYearProfit     float64            `json:"last_year_profit,omitempty"`
	AssetValue         float64            `json:"asset_value,omitempty"`
	Amount             float64            `json:"amount,omitempty"`
	PreAssesment       int                `json:"pre_assesment,omitempty"`
	AccountingSoftware string             `json:"accounting_software,omitempty"`
	Status             string             `json:"status,omitempty"`
	Code               string             `json:"code,omitempty"`
	ActiveSoftware     bool               `json:"active_software,omitempty"`
	ActiveDecision     bool               `json:"active_decision,omitempty"`
	Base               support.Base       `json:"base,omitempty"`
}
type BalanceSheet struct {
	Year         int     `json:"year,omitempty"`
	Month        int     `json:"month,omitempty"`
	ProfitOrLoss float64 `json:"profitOrLoss,omitempty"`
	AssetsValue  float64 `json:"assetsValue,omitempty"`
}
type Results struct {
	Data []BalanceSheet `json:"data,omitempty"`
	Loan Loan           `json:"loan,omitempty"`
}

var MessageResp struct {
	Message string `json:"message,omitempty"`
}

func (l Loan) Validate() httperrors.HttpErr {
	if l.Name == "" {
		return httperrors.NewBadRequestError("Name should not be empty")
	}
	if l.BusinessPin == "" {
		return httperrors.NewBadRequestError("Business Pin should not be empty")
	}
	if l.YearEstablished == 0 {
		return httperrors.NewBadRequestError("you need an already established business")
	}
	return nil
}
