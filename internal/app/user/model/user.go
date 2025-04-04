package model

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type User struct {
	ID        bson.ObjectID `bson:"_id" json:"-"`
	Name      string        `bson:"name" json:"name"`
	Email     string        `bson:"email,omitempty" json:"email"`
	Phone     string        `bson:"phone,omitempty" json:"phone"`
	Password  string        `bson:"password" json:"password"`
	CreatedAt time.Time     `bson:"createdAt" json:"-"`
	UpdatedAt time.Time     `bson:"updatedAt" json:"-"`
	DeletedAt *time.Time    `bson:"deletedAt,omitempty" json:"-"`
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
