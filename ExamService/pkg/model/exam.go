package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Exam struct {
	Id          primitive.ObjectID `json:"_id"         bson:"_id,omitempty"`
	Topic       string             `json:"topic"       bson:"topic"`
	Description string             `json:"description" bson:"description"`
	Tags        []string           `json:"tags"        bson:"tags"`
	IsPublic    bool               `json:"isPublic"    bson:"isPublic"`
	UserId      string             `json:"userId"      bson:"userId"`
	CreatedAt   time.Time          `json:"createdAt"   bson:"createdAt"`
	UpdatedAt   time.Time          `json:"updatedAt"   bson:"updatedAt"`
}
