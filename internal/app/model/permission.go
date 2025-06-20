package model

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type Permission struct {
	ID        bson.ObjectID `bson:"_id,omitempty" json:"id"`
	Name      string        `bson:"name" json:"name" fillable`
	DeletedAt *time.Time    `bson:"deleted_at,omitempty" json:"deletedAt,omitempty"`
}

func NewPermission() *Permission {
	return &Permission{
		ID: bson.NewObjectID(),
	}
}
func (p *Permission) TableName() string {
	return "permissions"
}

func (p *Permission) Default() {
	//...
}
