package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AnswerWrong struct {
	Id         primitive.ObjectID `json:"_id"        bson:"_id,omitempty"`
	ExamId     string             `json:"examId"     bson:"examId"`
	QuestionId string             `json:"questionId" bson:"questionId"`
	Times      int32              `json:"times"      bson:"times"`
	UserId     string             `json:"userId"     bson:"userId"`
	CreatedAt  time.Time          `json:"createdAt"  bson:"createdAt"`
	UpdatedAt  time.Time          `json:"updatedAt"  bson:"updatedAt"`
}
