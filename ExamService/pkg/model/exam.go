package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson"
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

/*
Autofill created_at and updated_at in golang struct while pushing into mongodb
https://stackoverflow.com/questions/71902455/autofill-created-at-and-updated-at-in-golang-struct-while-pushing-into-mongodb
*/
func (u *Exam) MarshalBSON() ([]byte, error) {
	now := time.Now()

	if u.CreatedAt.IsZero() {
		u.CreatedAt = now
	}

	u.UpdatedAt = now

	type my Exam
	return bson.Marshal((*my)(u))
}
