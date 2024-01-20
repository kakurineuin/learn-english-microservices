package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	Id        primitive.ObjectID `bson:"_id,omitempty"`
	Username  string             `bson:"username"`
	Password  string             `bson:"password"`
	Role      string             `bson:"role"`
	CreatedAt time.Time          `bson:"createdAt"`
	UpdatedAt time.Time          `bson:"updatedAt"`
}

/*
Autofill created_at and updated_at in golang struct while pushing into mongodb
https://stackoverflow.com/questions/71902455/autofill-created-at-and-updated-at-in-golang-struct-while-pushing-into-mongodb
*/
func (u *User) MarshalBSON() ([]byte, error) {
	now := time.Now()

	if u.CreatedAt.IsZero() {
		u.CreatedAt = now
	}

	u.UpdatedAt = now

	type my User
	return bson.Marshal((*my)(u))
}
