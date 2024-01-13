package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type FavoriteWordMeaning struct {
	Id            primitive.ObjectID `json:"_id"           bson:"_id,omitempty"`
	UserId        string             `json:"userId"        bson:"userId"`
	WordMeaningId primitive.ObjectID `json:"wordMeaningId" bson:"wordMeaningId"`
	CreatedAt     time.Time          `json:"createdAt"     bson:"createdAt"`
	UpdatedAt     time.Time          `json:"updatedAt"     bson:"updatedAt"`
}
