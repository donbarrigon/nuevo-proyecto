package model

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type User struct {
	ID        bson.ObjectID `bson:"_id,omitempty" json:"-"`
	Name      string        `bson:"name" json:"name" fillable:"true"`
	Email     string        `bson:"email,omitempty" json:"email" fillable:"true"`
	Phone     string        `bson:"phone,omitempty" json:"phone" fillable:"true"`
	Password  string        `bson:"password" json:"password"`
	Tokens    *[]Token      `bson:"tokens" json:"tokens"`
	CreatedAt time.Time     `bson:"created_at" json:"-"`
	UpdatedAt time.Time     `bson:"updated_at" json:"-"`
	DeletedAt *time.Time    `bson:"deleted_at,omitempty" json:"-"`
}

func NewUser() *User {
	return &User{
		ID:        bson.NewObjectID(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

func (u *User) TableName() string {
	return "users"
}

func (u *User) Default() {
	if u.CreatedAt.IsZero() {
		u.CreatedAt = time.Now()
	}
	u.UpdatedAt = time.Now()
}

func (u *User) Anonymous() *User {
	var id bson.ObjectID // zero value: "000000000000000000000000"
	var timeZero time.Time

	return &User{
		ID:        id,
		Name:      "Anonymous",
		Email:     "anonymous@gmail.com",
		Phone:     "+57 320 000 0000",
		Password:  "anonymous",
		CreatedAt: timeZero,
		UpdatedAt: timeZero,
		DeletedAt: &timeZero,
	}
}
