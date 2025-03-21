package model

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type User struct {
	ID        bson.ObjectID `bson:"_id"`
	Name      string        `bson:"name"`
	Email     *string       `bson:"email,omitempty"`
	Phone     *string       `bson:"phone,omitempty"`
	Password  string        `bson:"password"`
	CreatedAt time.Time     `bson:"createdAt"`
	UpdatedAt time.Time     `bson:"updatedAt"`
	DeletedAt *time.Time    `bson:"deletedAt,omitempty"`
}

func NewUser() *User {
	return &User{
		ID:        bson.NewObjectID(),
		CreatedAt: time.Now(),
	}
}

func (u *User) CollectionName() string {
	return "users"
}

func (u *User) GetID() bson.ObjectID {
	return u.ID
}

func (u *User) Index() map[string]string {
	return map[string]string{
		"token":     "unique",
		"email":     "index",
		"phone":     "index",
		"deletedAt": "index",
	}
}
