package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Question struct {
	Id        primitive.ObjectID `json:"_id"       bson:"_id,omitempty"`
	ExamId    string             `json:"examId"    bson:"examId"`
	Ask       string             `json:"ask"       bson:"ask"`
	Answers   []string           `json:"answers"   bson:"answers"`
	UserId    string             `json:"userId"    bson:"userId"`
	CreatedAt time.Time          `json:"createdAt" bson:"createdAt"`
	UpdatedAt time.Time          `json:"updatedAt" bson:"updatedAt"`
}
