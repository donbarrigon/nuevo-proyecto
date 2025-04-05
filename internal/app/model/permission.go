package model

import (
	"github.com/donbarrigon/nuevo-proyecto/pkg/errors"
	"github.com/donbarrigon/nuevo-proyecto/pkg/lang"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type Permission struct {
	ID   bson.ObjectID `bson:"_id"`
	Name string        `bson:"name"`
}

func NewPermission() *Permission {
	return &Permission{
		ID: bson.NewObjectID(),
	}
}
func (p *Permission) CollectionName() string {
	return "permissions"
}
func (p *Permission) GetID() bson.ObjectID {
	return p.ID
}
func (p *Permission) SetID(id bson.ObjectID) {
	p.ID = id
}
func (p *Permission) Default() {
	//...
}

func (p *Permission) Validate(l string) errors.Error {
	err := &errors.Err{}

	if len(p.Name) > 255 {
		err.Append("name", lang.TT(l, "Maximo %v caracteres", 255))
	}
	return err.Errors()
}
