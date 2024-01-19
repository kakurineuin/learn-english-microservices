package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserHistory struct {
	Id        primitive.ObjectID `json:"_id"       bson:"_id,omitempty"`
	UserId    string             `json:"userId"    bson:"userId"`
	Method    string             `json:"method"    bson:"method"`
	Path      string             `json:"path"      bson:"path"`
	CreatedAt time.Time          `json:"createdAt" bson:"createdAt"`
	UpdatedAt time.Time          `json:"updatedAt" bson:"updatedAt"`
}
