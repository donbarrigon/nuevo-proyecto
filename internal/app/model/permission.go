package model

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type Permission struct {
	ID        bson.ObjectID `bson:"_id,omitempty" json:"-"`
	Name      string        `bson:"name" json:"name"`
	DeletedAt *time.Time    `bson:"deleted_at,omitempty" json:"-"`
}

func NewPermission() *Permission {
	return &Permission{
		ID: bson.NewObjectID(),
	}
}
func (p *Permission) CollectionName() string {
	return "permissions"
}

func (p *Permission) Default() {
	//...
}
