package users

import (
	"regexp"
	"time"

	httperrors "github.com/myrachanto/erroring"
	"github.com/myrachanto/estate/src/support"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID        primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Fullname  string             `json:"fullname,omitempty"`
	Username  string             `json:"username,omitempty"`
	Birthday  string             `json:"birthday,omitempty"`
	Address   string             `json:"address,omitempty"`
	Phone     string             `json:"phone,omitempty"`
	Email     string             `json:"email,omitempty"`
	Password  string             `json:"password,omitempty"`
	Usercode  string             `json:"usercode,omitempty"`
	Picture   string             `json:"picture,omitempty"`
	UserAgent string             `json:"user_agent,omitempty"`
	Admin     bool               `json:"admin,omitempty"`
	Auditor   bool               `json:"auditor,omitempty"`
	Base      support.Base       `json:"base,omitempty"`
}
type LoginUser struct {
	Email     string `json:"email,omitempty"`
	Password  string `json:"password,omitempty"`
	UserAgent string `json:"user_agent,omitempty"`
}

type Auth struct {
	Usercode            string    `json:"usercode,omitempty"`
	UserName            string    `json:"username,omitempty"`
	Picture             string    `json:"picture,omitempty"`
	Admin               bool      `json:"admin,omitempty"`
	Auditor             bool      `json:"auditor,omitempty"`
	Token               string    `bson:"token" json:"token,omitempty"`
	TokenExpires        time.Time `json:"token_expires,omitempty"`
	RefleshToken        string    `json:"reflesh_token,omitempty"`
	RefleshTokenExpires time.Time `json:"reflesh_token_expires,omitempty"`
	SessionCode         string    `json:"session_code,omitempty"`
}
type Session struct {
	ID           primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Code         string             `json:"code,omitempty"`
	Username     string             `json:"username,omitempty"`
	Usercode     string             `json:"usercode,omitempty"`
	RefleshToken string             `json:"reflesh_token,omitempty"`
	TokenId      string             `json:"token_id,omitempty"`
	UserAgent    string             `json:"user_agent,omitempty"`
	ClientIp     string             `json:"client_ip,omitempty"`
	IsBlocked    bool               `json:"is_blocked,omitempty"`
	ExpiresAt    time.Time          `json:"expires_at,omitempty"`
	Base         support.Base       `json:"base,omitempty"`
}

type Results struct {
	Data        []*User `json:"results"`
	Total       int     `json:"total"`
	Pages       int     `json:"pages"`
	CurrentPage int     `json:"currentpage"`
}
type Base struct {
	Created_At time.Time  `bson:"created_at"`
	Updated_At time.Time  `bson:"updated_at"`
	Delete_At  *time.Time `bson:"deleted_at"`
}

func (user User) ValidateEmail(email string) (matchedString bool) {
	stringresults := httperrors.ValidStringNotEmpty(email)
	if stringresults.Noerror() {
		return false
	}
	re := regexp.MustCompile("^[a-zA-Z0-9.!#$%&amp;'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
	matchedString = re.MatchString(email)
	return matchedString
}
func (user User) ValidatePassword(password string) (bool, httperrors.HttpErr) {
	stringresults := httperrors.ValidStringNotEmpty(password)
	if stringresults.Noerror() {
		return false, stringresults
	}
	if len(password) < 5 {
		return false, httperrors.NewBadRequestError("your password need more characters!")
	} else if len(password) > 32 {
		return false, httperrors.NewBadRequestError("your password is way too long!")
	}
	return true, nil
}
func (user User) HashPassword(password string) (string, httperrors.HttpErr) {
	stringresults := httperrors.ValidStringNotEmpty(password)
	if stringresults.Noerror() {
		return "", httperrors.NewBadRequestError("your password Must not be empty!")
	}
	pass, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	if err != nil {
		return "", httperrors.NewNotFoundError("soemthing went wrong!")
	}
	return string(pass), nil

}

func (user User) Compare(p1, p2 string) bool {
	stringresults := httperrors.ValidStringNotEmpty(p1)
	if stringresults.Noerror() {
		return false
	}
	stringresults2 := httperrors.ValidStringNotEmpty(p2)
	if stringresults2.Noerror() {
		return false
	}
	err := bcrypt.CompareHashAndPassword([]byte(p2), []byte(p1))
	return err == nil
}
func (user LoginUser) Compare(p1, p2 string) bool {
	stringresults := httperrors.ValidStringNotEmpty(p1)
	if stringresults.Noerror() {
		return false
	}
	stringresults2 := httperrors.ValidStringNotEmpty(p2)
	if stringresults2.Noerror() {
		return false
	}
	err := bcrypt.CompareHashAndPassword([]byte(p2), []byte(p1))
	return err == nil
}
func (u User) Validate() httperrors.HttpErr {
	if u.Username == "" {
		return httperrors.NewBadRequestError("Username should not be empty")
	}
	if u.Email == "" {
		return httperrors.NewBadRequestError("Email should not be empty")
	}
	if u.Password == "" {
		return httperrors.NewBadRequestError("Password should not be empty")
	}
	return nil
}

func (u LoginUser) Validate() httperrors.HttpErr {
	if u.Email == "" {
		return httperrors.NewNotFoundError("Invalid Email")
	}
	if u.Password == "" {
		return httperrors.NewNotFoundError("Invalid password")
	}
	return nil
}
