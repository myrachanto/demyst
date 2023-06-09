package pasetos

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	httperrors "github.com/myrachanto/erroring"
)

var (
	ErrExpiredToken = "token has expired"
)

type Payload struct {
	IDs       uuid.UUID `json:"id"`
	Code      string    `json:"code"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Admin     bool      `json:"admin"`
	IssuedAt  time.Time `json:"issued_at"`
	ExpiredAt time.Time `json:"expired_at"`
}
type Data struct {
	Code     string `json:"code"`
	Usercode string `json:"usercode"`
	Username string `json:"username"`
	Admin    bool   `json:"admin"`
	Email    string `json:"email"`
}

func NewPayload(data *Data, duration time.Duration) (*Payload, httperrors.HttpErr) {
	tokenid, err := uuid.NewRandom()
	if err != nil {
		return nil, httperrors.NewBadRequestError(fmt.Sprintf("error with uuid generation, %d", err))
	}
	// fmt.Println("paseto setting ==============", data)
	return &Payload{
		IDs:       tokenid,
		Username:  data.Username,
		Email:     data.Email,
		Code:      data.Code,
		Admin:     data.Admin,
		IssuedAt:  time.Now(),
		ExpiredAt: time.Now().Add(duration),
	}, nil
}
func (payload *Payload) Valid() httperrors.HttpErr {
	if time.Now().After(payload.ExpiredAt) {
		return httperrors.NewBadRequestError(ErrExpiredToken)
	}
	return nil
}
