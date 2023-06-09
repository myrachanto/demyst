package newssections

import (
	"github.com/myrachanto/sports/src/support"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type NewsSection struct {
	ID        primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Name      string             `json:"name,omitempty"`
	Content   string             `json:"content,omitempty"`
	Image     string             `json:"image,omitempty"`
	Highlight string             `json:"highlight,omitempty"`
	Code      string             `json:"code,omitempty"`
	Base      support.Base       `json:"base,omitempty"`
}
