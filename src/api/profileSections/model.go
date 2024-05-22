package profilesections

import (
	"github.com/myrachanto/estate/src/support"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ProfileSection struct {
	ID        primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Name      string             `json:"name"`
	Content   string             `json:"content"`
	Image     string             `json:"image"`
	Highlight string             `json:"highlight"`
	Code      string             `json:"code"`
	Instagram string             `json:"instagram"`
	Twitter   string             `json:"twitter"`
	Youtube   string             `json:"youtube"`
	Facebook  string             `json:"facebook"`
	Base      support.Base       `json:"base"`
}
