package model

import (
	"go.mongodb.org/mongo-driver/v2/bson"
)

type Role struct {
	ID          bson.ObjectID `bson:"_id,omitempty" json:"-"`
	Name        string        `bson:"name" json:"name"`
	Permissions []Permission  `bson:"permissions" json:"-"`
}

func NewRole() *Role {
	return &Role{
		ID: bson.NewObjectID(),
	}
}
func (r *Role) TableName() string {
	return "roles"
}

func (r *Role) Default() {
	//...
}
