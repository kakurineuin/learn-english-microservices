package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ExamRecord struct {
	Id        primitive.ObjectID `json:"_id"       bson:"_id,omitempty"`
	ExamId    string             `json:"examId"    bson:"examId"`
	Score     int32              `json:"score"     bson:"score"`
	UserId    string             `json:"userId"    bson:"userId"`
	CreatedAt time.Time          `json:"createdAt" bson:"createdAt"`
	UpdatedAt time.Time          `json:"updatedAt" bson:"updatedAt"`
}
